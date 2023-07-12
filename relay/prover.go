package relay

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/datachainlab/ethereum-ibc-relay-chain/pkg/client"
	"github.com/datachainlab/ethereum-ibc-relay-chain/pkg/relay/ethereum"
	"github.com/datachainlab/ethereum-ibc-relay-prover/beacon"
	lctypes "github.com/datachainlab/ethereum-ibc-relay-prover/light-clients/ethereum/types"
	"github.com/hyperledger-labs/yui-relayer/core"
)

type Prover struct {
	chain           *ethereum.Chain
	config          ProverConfig
	executionClient *client.ETHClient
	beaconClient    beacon.Client
	codec           codec.ProtoCodecMarshaler
}

func NewProver(chain *ethereum.Chain, config ProverConfig) *Prover {
	beaconClient := beacon.NewClient(config.BeaconEndpoint)
	return &Prover{chain: chain, config: config, executionClient: chain.Client(), beaconClient: beaconClient}
}

//--------- Prover implementation ---------//

var _ core.Prover = (*Prover)(nil)

// Init initializes the chain
func (pr *Prover) Init(homePath string, timeout time.Duration, codec codec.ProtoCodecMarshaler, debug bool) error {
	pr.codec = codec
	return nil
}

// SetRelayInfo sets source's path and counterparty's info to the chain
func (pr *Prover) SetRelayInfo(path *core.PathEnd, counterparty *core.ProvableChain, counterpartyPath *core.PathEnd) error {
	return nil
}

// SetupForRelay performs chain-specific setup before starting the relay
func (pr *Prover) SetupForRelay(ctx context.Context) error {
	return nil
}

//--------- LightClient implementation ---------//

var _ core.LightClient = (*Prover)(nil)

// CreateMsgCreateClient creates a CreateClientMsg to this chain
func (pr *Prover) CreateMsgCreateClient(clientID string, dstHeader core.Header, signer sdk.AccAddress) (*clienttypes.MsgCreateClient, error) {
	// NOTE: `dstHeader` generated by GetLatestFinalizedHeader doesn't have the next sync committe
	// So, we don't use it here, use the result of `light_client_update` to create an initial state instead.
	_ = dstHeader.(*lctypes.Header)

	genesis, err := pr.beaconClient.GetGenesis()
	if err != nil {
		return nil, err
	}
	bootstrap, err := pr.getLightClientBootstrap()
	if err != nil {
		return nil, err
	}
	accountUpdate, err := pr.buildAccountUpdate(bootstrap.Header.Execution.BlockNumber)
	if err != nil {
		return nil, err
	}

	clientState := pr.newClientState()
	clientState.GenesisValidatorsRoot = genesis.GenesisValidatorsRoot[:]
	clientState.GenesisTime = genesis.GenesisTimeSeconds
	clientState.LatestSlot = uint64(bootstrap.Header.Beacon.Slot)
	clientState.LatestExecutionBlockNumber = bootstrap.Header.Execution.BlockNumber

	consensusState := &lctypes.ConsensusState{
		Slot:                 clientState.LatestSlot,
		StorageRoot:          accountUpdate.AccountStorageRoot,
		Timestamp:            bootstrap.Header.Execution.Timestamp,
		CurrentSyncCommittee: bootstrap.CurrentSyncCommittee.AggregatePubKey,
	}

	return clienttypes.NewMsgCreateClient(clientState, consensusState, signer.String())
}

// SetupHeadersForUpdate returns the finalized header and any intermediate headers needed to apply it to the client on the counterpaty chain
// The order of the returned header slice should be as: [<intermediate headers>..., <update header>]
// if the header slice's length == 0 and err == nil, the relayer should skips the update-client
func (pr *Prover) SetupHeadersForUpdate(dstChain core.ChainInfoICS02Querier, latestFinalizedHeader core.Header) ([]core.Header, error) {
	finalizedHeader := latestFinalizedHeader.(*lctypes.Header)

	latestHeight, err := dstChain.LatestHeight()
	if err != nil {
		return nil, err
	}

	// retrieve counterparty client from dst chain
	counterpartyClientRes, err := dstChain.QueryClientState(core.NewQueryContext(context.TODO(), latestHeight))
	if err != nil {
		return nil, err
	}
	var cs ibcexported.ClientState
	if err := pr.codec.UnpackAny(counterpartyClientRes.ClientState, &cs); err != nil {
		return nil, err
	}

	if cs.GetLatestHeight().GetRevisionHeight() >= finalizedHeader.ExecutionUpdate.BlockNumber {
		return nil, fmt.Errorf("the latest finalized header is equal to or older than the latest height of client state: finalized_block_number=%v client_latest_height=%v", finalizedHeader.ExecutionUpdate.BlockNumber, cs.GetLatestHeight().GetRevisionHeight())
	}

	latestPeriod := pr.computeSyncCommitteePeriod(pr.computeEpoch(finalizedHeader.ConsensusUpdate.FinalizedHeader.Slot))
	statePeriod, err := pr.findLCPeriodByHeight(cs.GetLatestHeight(), latestPeriod)
	if err != nil {
		return nil, err
	}

	log.Printf("try to setup headers for updating the light-client: lc_latest_height=%v lc_latest_height_period=%v latest_period=%v", cs.GetLatestHeight(), statePeriod, latestPeriod)

	if statePeriod > latestPeriod {
		return nil, fmt.Errorf("the light-client server's response is old: client_state_period=%v latest_finalized_period=%v", statePeriod, latestPeriod)
	} else if statePeriod == latestPeriod {
		latestHeight := cs.GetLatestHeight().(clienttypes.Height)
		res, err := pr.beaconClient.GetLightClientUpdate(statePeriod)
		if err != nil {
			return nil, err
		}
		root, err := res.Data.FinalizedHeader.Beacon.HashTreeRoot()
		if err != nil {
			return nil, err
		}
		bootstrapRes, err := pr.beaconClient.GetBootstrap(root[:])
		if err != nil {
			return nil, err
		}
		finalizedHeader.TrustedSyncCommittee = &lctypes.TrustedSyncCommittee{
			TrustedHeight: &latestHeight,
			SyncCommittee: bootstrapRes.Data.CurrentSyncCommittee.ToProto(),
			IsNext:        false,
		}
		return []core.Header{finalizedHeader}, nil
	}

	//--------- In case statePeriod < latestPeriod ---------//

	var (
		headers              []core.Header
		trustedSyncCommittee *lctypes.SyncCommittee
		trustedHeight        = cs.GetLatestHeight().(clienttypes.Height)
	)

	for p := statePeriod; p < latestPeriod; p++ {
		var header *lctypes.Header
		if p == statePeriod {
			header, err = pr.buildNextSyncCommitteeUpdateForCurrent(statePeriod, trustedHeight)
			if err != nil {
				return nil, err
			}
		} else {
			header, err = pr.buildNextSyncCommitteeUpdateForNext(p, trustedHeight)
			if err != nil {
				return nil, err
			}
		}
		headers = append(headers, header)
		trustedHeight = clienttypes.NewHeight(0, header.ExecutionUpdate.BlockNumber)
		trustedSyncCommittee = header.ConsensusUpdate.NextSyncCommittee
	}

	finalizedHeader.TrustedSyncCommittee = &lctypes.TrustedSyncCommittee{
		TrustedHeight: &trustedHeight,
		SyncCommittee: trustedSyncCommittee,
		IsNext:        true,
	}
	headers = append(headers, finalizedHeader)
	return headers, nil
}

// GetLatestFinalizedHeader returns the latest finalized header on this chain
// The returned header is expected to be the latest one of headers that can be verified by the light client
func (pr *Prover) GetLatestFinalizedHeader() (headers core.Header, err error) {
	res, err := pr.beaconClient.GetLightClientFinalityUpdate()
	if err != nil {
		return nil, err
	}
	lcUpdate := res.Data.ToProto()
	executionHeader := res.Data.FinalizedHeader.Execution
	executionUpdate, err := pr.buildExecutionUpdate(executionHeader)
	if err != nil {
		return nil, err
	}
	executionRoot, err := executionHeader.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(executionRoot[:], lcUpdate.FinalizedExecutionRoot) {
		return nil, fmt.Errorf("execution root mismatch: %X != %X", executionRoot, lcUpdate.FinalizedExecutionRoot)
	}

	accountUpdate, err := pr.buildAccountUpdate(executionHeader.BlockNumber)
	if err != nil {
		return nil, err
	}
	return &lctypes.Header{
		ConsensusUpdate: lcUpdate,
		ExecutionUpdate: executionUpdate,
		AccountUpdate:   accountUpdate,
		Timestamp:       executionHeader.Timestamp,
	}, nil
}

func (pr *Prover) newClientState() *lctypes.ClientState {
	var commitmentsSlot [32]byte
	ibcAddress := pr.chain.Config().IBCAddress()

	return &lctypes.ClientState{
		ForkParameters:               pr.config.getForkParameters(),
		SecondsPerSlot:               pr.config.getSecondsPerSlot(),
		SlotsPerEpoch:                pr.config.getSlotsPerEpoch(),
		EpochsPerSyncCommitteePeriod: pr.config.getEpochsPerSyncCommitteePeriod(),

		MinSyncCommitteeParticipants: 1,

		IbcAddress:         ibcAddress.Bytes(),
		IbcCommitmentsSlot: commitmentsSlot[:],
		TrustLevel: &lctypes.Fraction{
			Numerator:   2,
			Denominator: 3,
		},
		TrustingPeriod: 0,
	}
}

func (pr *Prover) buildNextSyncCommitteeUpdateForCurrent(period uint64, trustedHeight clienttypes.Height) (*lctypes.Header, error) {
	res, err := pr.beaconClient.GetLightClientUpdate(period)
	if err != nil {
		return nil, err
	}
	lcUpdate := res.Data.ToProto()
	executionHeader := res.Data.FinalizedHeader.Execution
	executionUpdate, err := pr.buildExecutionUpdate(executionHeader)
	if err != nil {
		return nil, err
	}
	executionRoot, err := executionHeader.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(executionRoot[:], lcUpdate.FinalizedExecutionRoot) {
		return nil, fmt.Errorf("execution root mismatch: %X != %X", executionRoot, lcUpdate.FinalizedExecutionRoot)
	}

	accountUpdate, err := pr.buildAccountUpdate(executionHeader.BlockNumber)
	if err != nil {
		return nil, err
	}

	root, err := res.Data.FinalizedHeader.Beacon.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	bootstrapRes, err := pr.beaconClient.GetBootstrap(root[:])
	if err != nil {
		return nil, err
	}

	return &lctypes.Header{
		TrustedSyncCommittee: &lctypes.TrustedSyncCommittee{
			TrustedHeight: &trustedHeight,
			SyncCommittee: bootstrapRes.Data.CurrentSyncCommittee.ToProto(),
			IsNext:        false,
		},
		ConsensusUpdate: lcUpdate,
		ExecutionUpdate: executionUpdate,
		AccountUpdate:   accountUpdate,
		Timestamp:       executionHeader.Timestamp,
	}, nil
}

func (pr *Prover) buildNextSyncCommitteeUpdateForNext(period uint64, trustedHeight clienttypes.Height) (*lctypes.Header, error) {
	res, err := pr.beaconClient.GetLightClientUpdate(period)
	if err != nil {
		return nil, err
	}
	lcUpdate := res.Data.ToProto()
	executionHeader := res.Data.FinalizedHeader.Execution
	executionUpdate, err := pr.buildExecutionUpdate(executionHeader)
	if err != nil {
		return nil, err
	}
	executionRoot, err := executionHeader.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(executionRoot[:], lcUpdate.FinalizedExecutionRoot) {
		return nil, fmt.Errorf("execution root mismatch: %X != %X", executionRoot, lcUpdate.FinalizedExecutionRoot)
	}

	accountUpdate, err := pr.buildAccountUpdate(executionHeader.BlockNumber)
	if err != nil {
		return nil, err
	}
	return &lctypes.Header{
		TrustedSyncCommittee: &lctypes.TrustedSyncCommittee{
			TrustedHeight: &trustedHeight,
			SyncCommittee: lcUpdate.NextSyncCommittee,
			IsNext:        true,
		},
		ConsensusUpdate: lcUpdate,
		ExecutionUpdate: executionUpdate,
		AccountUpdate:   accountUpdate,
		Timestamp:       executionHeader.Timestamp,
	}, nil
}

//--------- IBCProvableQuerier implementation ---------//

var _ core.Prover = (*Prover)(nil)

// ProveState returns the proof of an IBC state specified by `path` and `value`
func (pr *Prover) ProveState(ctx core.QueryContext, path string, value []byte) ([]byte, clienttypes.Height, error) {
	// clientCtx := pr.chain.CLIContext(int64(ctx.Height().GetRevisionHeight()))
	// if v, proof, proofHeight, err := ibcclient.QueryTendermintProof(clientCtx, []byte(path)); err != nil {
	// 	return nil, clienttypes.Height{}, err
	// } else if !bytes.Equal(v, value) {
	// 	return nil, clienttypes.Height{}, fmt.Errorf("value unmatch: %x != %x", v, value)
	// } else {
	// 	return proof, proofHeight, nil
	// }
	h := sha256.Sum256(value)
	proof := h[:]
	return proof, ctx.Height().(clienttypes.Height), nil
}

// // QueryClientConsensusState returns the ClientConsensusState and its proof
// func (pr *Prover) QueryClientConsensusStateWithProof(ctx core.QueryContext, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
// 	res, err := pr.chain.QueryClientConsensusState(ctx, dstClientConsHeight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	proofHeight := int64(ctx.Height().GetRevisionHeight())
// 	res.Proof, err = pr.buildStateProof(host.FullConsensusStateKey(
// 		pr.chain.Path().ClientID,
// 		dstClientConsHeight,
// 	), proofHeight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res.ProofHeight = pr.newHeight(proofHeight)
// 	return res, nil
// }

// // QueryClientStateWithProof returns the ClientState and its proof
// func (pr *Prover) QueryClientStateWithProof(ctx core.QueryContext) (*clienttypes.QueryClientStateResponse, error) {
// 	res, err := pr.chain.QueryClientState(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	proofHeight := int64(ctx.Height().GetRevisionHeight())
// 	res.Proof, err = pr.buildStateProof(host.FullClientStateKey(
// 		pr.chain.Path().ClientID,
// 	), proofHeight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res.ProofHeight = pr.newHeight(proofHeight)
// 	return res, nil
// }

// // QueryConnectionWithProof returns the Connection and its proof
// func (pr *Prover) QueryConnectionWithProof(ctx core.QueryContext) (*conntypes.QueryConnectionResponse, error) {
// 	res, err := pr.chain.QueryConnection(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if res.Connection.State == conntypes.UNINITIALIZED {
// 		// connection not found
// 		return res, nil
// 	}
// 	proofHeight := int64(ctx.Height().GetRevisionHeight())
// 	res.Proof, err = pr.buildStateProof(host.ConnectionKey(
// 		pr.chain.Path().ConnectionID,
// 	), proofHeight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res.ProofHeight = pr.newHeight(proofHeight)
// 	return res, nil
// }

// // QueryChannelWithProof returns the Channel and its proof
// func (pr *Prover) QueryChannelWithProof(ctx core.QueryContext) (chanRes *chantypes.QueryChannelResponse, err error) {
// 	res, err := pr.chain.QueryChannel(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if res.Channel.State == chantypes.UNINITIALIZED {
// 		// channel not found
// 		return res, nil
// 	}
// 	proofHeight := int64(ctx.Height().GetRevisionHeight())
// 	res.Proof, err = pr.buildStateProof(host.ChannelKey(
// 		pr.chain.Path().PortID,
// 		pr.chain.Path().ChannelID,
// 	), proofHeight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res.ProofHeight = pr.newHeight(proofHeight)
// 	return res, nil
// }

// // QueryPacketCommitmentWithProof returns the packet commitment and its proof
// func (pr *Prover) QueryPacketCommitmentWithProof(ctx core.QueryContext, seq uint64) (comRes *chantypes.QueryPacketCommitmentResponse, err error) {
// 	res, err := pr.chain.QueryPacketCommitment(ctx, seq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	proofHeight := int64(ctx.Height().GetRevisionHeight())
// 	res.Proof, err = pr.buildStateProof(host.PacketCommitmentKey(
// 		pr.chain.Path().PortID,
// 		pr.chain.Path().ChannelID,
// 		seq,
// 	), proofHeight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res.ProofHeight = pr.newHeight(proofHeight)
// 	return res, nil
// }

// // QueryPacketAcknowledgementCommitmentWithProof returns the packet acknowledgement commitment and its proof
// func (pr *Prover) QueryPacketAcknowledgementCommitmentWithProof(ctx core.QueryContext, seq uint64) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error) {
// 	res, err := pr.chain.QueryPacketAcknowledgementCommitment(ctx, seq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	proofHeight := int64(ctx.Height().GetRevisionHeight())
// 	res.Proof, err = pr.buildStateProof(host.PacketAcknowledgementKey(
// 		pr.chain.Path().PortID,
// 		pr.chain.Path().ChannelID,
// 		seq,
// 	), proofHeight)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res.ProofHeight = pr.newHeight(proofHeight)
// 	return res, nil
// }

// func (pr *Prover) newHeight(blockNumber int64) clienttypes.Height {
// 	return clienttypes.NewHeight(0, uint64(blockNumber))
// }
