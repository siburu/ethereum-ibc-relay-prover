// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: relay/ethereum_light_client/config/config.proto

package relay

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type ProverConfig struct {
	BeaconEndpoint string `protobuf:"bytes,1,opt,name=beacon_endpoint,json=beaconEndpoint,proto3" json:"beacon_endpoint,omitempty"`
	Network        string `protobuf:"bytes,2,opt,name=network,proto3" json:"network,omitempty"`
}

func (m *ProverConfig) Reset()         { *m = ProverConfig{} }
func (m *ProverConfig) String() string { return proto.CompactTextString(m) }
func (*ProverConfig) ProtoMessage()    {}
func (*ProverConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_85af615077598949, []int{0}
}
func (m *ProverConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProverConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProverConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProverConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProverConfig.Merge(m, src)
}
func (m *ProverConfig) XXX_Size() int {
	return m.Size()
}
func (m *ProverConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ProverConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ProverConfig proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ProverConfig)(nil), "relay.ethereum_light_client.config.ProverConfig")
}

func init() {
	proto.RegisterFile("relay/ethereum_light_client/config/config.proto", fileDescriptor_85af615077598949)
}

var fileDescriptor_85af615077598949 = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x50, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x8d, 0x19, 0x40, 0x44, 0x08, 0xa4, 0x88, 0x21, 0x74, 0xb0, 0x50, 0x17, 0x58, 0x12, 0x0f,
	0x0c, 0xec, 0x20, 0x76, 0x40, 0x62, 0x61, 0x89, 0x6c, 0xf7, 0x6a, 0x5b, 0xb8, 0xbe, 0xc8, 0x72,
	0x40, 0xfd, 0x0b, 0x3e, 0xab, 0x63, 0x47, 0x46, 0x48, 0x7e, 0x04, 0x71, 0x2e, 0x4c, 0x4c, 0xf6,
	0x7b, 0xf7, 0xde, 0xe9, 0xdd, 0x2b, 0x45, 0x04, 0x2f, 0xd7, 0x02, 0x92, 0x85, 0x08, 0xc3, 0xaa,
	0xf3, 0xce, 0xd8, 0xd4, 0x69, 0xef, 0x20, 0x24, 0xa1, 0x31, 0x2c, 0x9d, 0xd9, 0x3d, 0x6d, 0x1f,
	0x31, 0x61, 0x35, 0x27, 0x43, 0xfb, 0xaf, 0xa1, 0xcd, 0xca, 0xd9, 0xa9, 0x41, 0x83, 0x24, 0x17,
	0x3f, 0xbf, 0xec, 0x9c, 0x9d, 0x19, 0x44, 0xe3, 0x41, 0x10, 0x52, 0xc3, 0x52, 0xc8, 0xb0, 0xce,
	0xa3, 0xf9, 0x43, 0x79, 0x74, 0x1f, 0xf1, 0x15, 0xe2, 0x2d, 0x2d, 0xa8, 0x2e, 0xca, 0x13, 0x05,
	0x52, 0x63, 0xe8, 0x20, 0x2c, 0x7a, 0x74, 0x21, 0xd5, 0xec, 0x9c, 0x5d, 0x1e, 0x3e, 0x1e, 0x67,
	0xfa, 0x6e, 0xc7, 0x56, 0x75, 0x79, 0x10, 0x20, 0xbd, 0x61, 0x7c, 0xa9, 0xf7, 0x48, 0xf0, 0x0b,
	0x6f, 0x9e, 0x36, 0x5f, 0xbc, 0xd8, 0x8c, 0x9c, 0x6d, 0x47, 0xce, 0x3e, 0x47, 0xce, 0xde, 0x27,
	0x5e, 0x6c, 0x27, 0x5e, 0x7c, 0x4c, 0xbc, 0x78, 0xbe, 0x36, 0x2e, 0xd9, 0x41, 0xb5, 0x1a, 0x57,
	0x62, 0x21, 0x93, 0xd4, 0x56, 0xba, 0xe0, 0xa5, 0xfa, 0x2b, 0xa2, 0x71, 0x4a, 0x37, 0x74, 0x6a,
	0xd3, 0x53, 0xb2, 0x5c, 0x94, 0xda, 0xa7, 0xc0, 0x57, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xb4,
	0x50, 0x61, 0x02, 0x38, 0x01, 0x00, 0x00,
}

func (m *ProverConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProverConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProverConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Network) > 0 {
		i -= len(m.Network)
		copy(dAtA[i:], m.Network)
		i = encodeVarintConfig(dAtA, i, uint64(len(m.Network)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.BeaconEndpoint) > 0 {
		i -= len(m.BeaconEndpoint)
		copy(dAtA[i:], m.BeaconEndpoint)
		i = encodeVarintConfig(dAtA, i, uint64(len(m.BeaconEndpoint)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProverConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.BeaconEndpoint)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	l = len(m.Network)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	return n
}

func sovConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozConfig(x uint64) (n int) {
	return sovConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ProverConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowConfig
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProverConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProverConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BeaconEndpoint", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BeaconEndpoint = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Network", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Network = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthConfig
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowConfig
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthConfig
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupConfig
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthConfig
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthConfig        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowConfig          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupConfig = fmt.Errorf("proto: unexpected end of group")
)
