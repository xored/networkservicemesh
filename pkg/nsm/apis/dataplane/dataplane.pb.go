// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dataplane.proto

package dataplane

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common "github.com/ligato/networkservicemesh/pkg/nsm/apis/common"
	math "math"
)

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Message sent by dataplane module informing NSM of any changes in its
// operations parameters or constraints
type MechanismUpdate struct {
	RemoteMechanisms     []*common.RemoteMechanism `protobuf:"bytes,1,rep,name=remote_mechanisms,json=remoteMechanisms,proto3" json:"remote_mechanisms,omitempty"`
	LocalMechanisms      []*common.LocalMechanism  `protobuf:"bytes,2,rep,name=local_mechanisms,json=localMechanisms,proto3" json:"local_mechanisms,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *MechanismUpdate) Reset()         { *m = MechanismUpdate{} }
func (m *MechanismUpdate) String() string { return proto.CompactTextString(m) }
func (*MechanismUpdate) ProtoMessage()    {}
func (*MechanismUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_617387e490a04ffa, []int{0}
}

func (m *MechanismUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MechanismUpdate.Unmarshal(m, b)
}
func (m *MechanismUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MechanismUpdate.Marshal(b, m, deterministic)
}
func (m *MechanismUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MechanismUpdate.Merge(m, src)
}
func (m *MechanismUpdate) XXX_Size() int {
	return xxx_messageInfo_MechanismUpdate.Size(m)
}
func (m *MechanismUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_MechanismUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_MechanismUpdate proto.InternalMessageInfo

func (m *MechanismUpdate) GetRemoteMechanisms() []*common.RemoteMechanism {
	if m != nil {
		return m.RemoteMechanisms
	}
	return nil
}

func (m *MechanismUpdate) GetLocalMechanisms() []*common.LocalMechanism {
	if m != nil {
		return m.LocalMechanisms
	}
	return nil
}

type Connection struct {
	LocalSource *common.LocalMechanism `protobuf:"bytes,1,opt,name=local_source,json=localSource,proto3" json:"local_source,omitempty"`
	// Types that are valid to be assigned to Destination:
	//	*Connection_Local
	//	*Connection_Remote
	Destination isConnection_Destination `protobuf_oneof:"destination"`
	// For ConnectRequest this will be nil, but for DisconnectRequest it will
	// carry a connection id returned by the dataplane controller
	ConnectionId         string   `protobuf:"bytes,4,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Connection) Reset()         { *m = Connection{} }
func (m *Connection) String() string { return proto.CompactTextString(m) }
func (*Connection) ProtoMessage()    {}
func (*Connection) Descriptor() ([]byte, []int) {
	return fileDescriptor_617387e490a04ffa, []int{1}
}

func (m *Connection) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Connection.Unmarshal(m, b)
}
func (m *Connection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Connection.Marshal(b, m, deterministic)
}
func (m *Connection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Connection.Merge(m, src)
}
func (m *Connection) XXX_Size() int {
	return xxx_messageInfo_Connection.Size(m)
}
func (m *Connection) XXX_DiscardUnknown() {
	xxx_messageInfo_Connection.DiscardUnknown(m)
}

var xxx_messageInfo_Connection proto.InternalMessageInfo

func (m *Connection) GetLocalSource() *common.LocalMechanism {
	if m != nil {
		return m.LocalSource
	}
	return nil
}

type isConnection_Destination interface {
	isConnection_Destination()
}

type Connection_Local struct {
	Local *common.LocalMechanism `protobuf:"bytes,2,opt,name=local,proto3,oneof"`
}

type Connection_Remote struct {
	Remote *common.RemoteMechanism `protobuf:"bytes,3,opt,name=remote,proto3,oneof"`
}

func (*Connection_Local) isConnection_Destination() {}

func (*Connection_Remote) isConnection_Destination() {}

func (m *Connection) GetDestination() isConnection_Destination {
	if m != nil {
		return m.Destination
	}
	return nil
}

func (m *Connection) GetLocal() *common.LocalMechanism {
	if x, ok := m.GetDestination().(*Connection_Local); ok {
		return x.Local
	}
	return nil
}

func (m *Connection) GetRemote() *common.RemoteMechanism {
	if x, ok := m.GetDestination().(*Connection_Remote); ok {
		return x.Remote
	}
	return nil
}

func (m *Connection) GetConnectionId() string {
	if m != nil {
		return m.ConnectionId
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Connection) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Connection_OneofMarshaler, _Connection_OneofUnmarshaler, _Connection_OneofSizer, []interface{}{
		(*Connection_Local)(nil),
		(*Connection_Remote)(nil),
	}
}

func _Connection_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Connection)
	// destination
	switch x := m.Destination.(type) {
	case *Connection_Local:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Local); err != nil {
			return err
		}
	case *Connection_Remote:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Remote); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Connection.Destination has unexpected type %T", x)
	}
	return nil
}

func _Connection_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Connection)
	switch tag {
	case 2: // destination.local
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(common.LocalMechanism)
		err := b.DecodeMessage(msg)
		m.Destination = &Connection_Local{msg}
		return true, err
	case 3: // destination.remote
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(common.RemoteMechanism)
		err := b.DecodeMessage(msg)
		m.Destination = &Connection_Remote{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Connection_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Connection)
	// destination
	switch x := m.Destination.(type) {
	case *Connection_Local:
		s := proto.Size(x.Local)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Connection_Remote:
		s := proto.Size(x.Remote)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Reply struct {
	Success       bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ExtendedError string `protobuf:"bytes,2,opt,name=extended_error,json=extendedError,proto3" json:"extended_error,omitempty"`
	// Connection in case of success contains information about mechanisms that have been chosen
	Connection           *Connection `protobuf:"bytes,3,opt,name=connection,proto3" json:"connection,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Reply) Reset()         { *m = Reply{} }
func (m *Reply) String() string { return proto.CompactTextString(m) }
func (*Reply) ProtoMessage()    {}
func (*Reply) Descriptor() ([]byte, []int) {
	return fileDescriptor_617387e490a04ffa, []int{2}
}

func (m *Reply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Reply.Unmarshal(m, b)
}
func (m *Reply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Reply.Marshal(b, m, deterministic)
}
func (m *Reply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Reply.Merge(m, src)
}
func (m *Reply) XXX_Size() int {
	return xxx_messageInfo_Reply.Size(m)
}
func (m *Reply) XXX_DiscardUnknown() {
	xxx_messageInfo_Reply.DiscardUnknown(m)
}

var xxx_messageInfo_Reply proto.InternalMessageInfo

func (m *Reply) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *Reply) GetExtendedError() string {
	if m != nil {
		return m.ExtendedError
	}
	return ""
}

func (m *Reply) GetConnection() *Connection {
	if m != nil {
		return m.Connection
	}
	return nil
}

func init() {
	proto.RegisterType((*MechanismUpdate)(nil), "dataplane.MechanismUpdate")
	proto.RegisterType((*Connection)(nil), "dataplane.Connection")
	proto.RegisterType((*Reply)(nil), "dataplane.Reply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DataplaneOperationsClient is the client API for DataplaneOperations service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DataplaneOperationsClient interface {
	MonitorMechanisms(ctx context.Context, in *common.Empty, opts ...grpc.CallOption) (DataplaneOperations_MonitorMechanismsClient, error)
	ConnectRequest(ctx context.Context, in *Connection, opts ...grpc.CallOption) (*Reply, error)
	DisconnectRequest(ctx context.Context, in *Connection, opts ...grpc.CallOption) (*Reply, error)
}

type dataplaneOperationsClient struct {
	cc *grpc.ClientConn
}

func NewDataplaneOperationsClient(cc *grpc.ClientConn) DataplaneOperationsClient {
	return &dataplaneOperationsClient{cc}
}

func (c *dataplaneOperationsClient) MonitorMechanisms(ctx context.Context, in *common.Empty, opts ...grpc.CallOption) (DataplaneOperations_MonitorMechanismsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_DataplaneOperations_serviceDesc.Streams[0], "/dataplane.DataplaneOperations/MonitorMechanisms", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataplaneOperationsMonitorMechanismsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DataplaneOperations_MonitorMechanismsClient interface {
	Recv() (*MechanismUpdate, error)
	grpc.ClientStream
}

type dataplaneOperationsMonitorMechanismsClient struct {
	grpc.ClientStream
}

func (x *dataplaneOperationsMonitorMechanismsClient) Recv() (*MechanismUpdate, error) {
	m := new(MechanismUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dataplaneOperationsClient) ConnectRequest(ctx context.Context, in *Connection, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/dataplane.DataplaneOperations/ConnectRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataplaneOperationsClient) DisconnectRequest(ctx context.Context, in *Connection, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/dataplane.DataplaneOperations/DisconnectRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataplaneOperationsServer is the server API for DataplaneOperations service.
type DataplaneOperationsServer interface {
	MonitorMechanisms(*common.Empty, DataplaneOperations_MonitorMechanismsServer) error
	ConnectRequest(context.Context, *Connection) (*Reply, error)
	DisconnectRequest(context.Context, *Connection) (*Reply, error)
}

func RegisterDataplaneOperationsServer(s *grpc.Server, srv DataplaneOperationsServer) {
	s.RegisterService(&_DataplaneOperations_serviceDesc, srv)
}

func _DataplaneOperations_MonitorMechanisms_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(common.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DataplaneOperationsServer).MonitorMechanisms(m, &dataplaneOperationsMonitorMechanismsServer{stream})
}

type DataplaneOperations_MonitorMechanismsServer interface {
	Send(*MechanismUpdate) error
	grpc.ServerStream
}

type dataplaneOperationsMonitorMechanismsServer struct {
	grpc.ServerStream
}

func (x *dataplaneOperationsMonitorMechanismsServer) Send(m *MechanismUpdate) error {
	return x.ServerStream.SendMsg(m)
}

func _DataplaneOperations_ConnectRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Connection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataplaneOperationsServer).ConnectRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dataplane.DataplaneOperations/ConnectRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataplaneOperationsServer).ConnectRequest(ctx, req.(*Connection))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataplaneOperations_DisconnectRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Connection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataplaneOperationsServer).DisconnectRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dataplane.DataplaneOperations/DisconnectRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataplaneOperationsServer).DisconnectRequest(ctx, req.(*Connection))
	}
	return interceptor(ctx, in, info, handler)
}

var _DataplaneOperations_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dataplane.DataplaneOperations",
	HandlerType: (*DataplaneOperationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnectRequest",
			Handler:    _DataplaneOperations_ConnectRequest_Handler,
		},
		{
			MethodName: "DisconnectRequest",
			Handler:    _DataplaneOperations_DisconnectRequest_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MonitorMechanisms",
			Handler:       _DataplaneOperations_MonitorMechanisms_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "dataplane.proto",
}

func init() { proto.RegisterFile("dataplane.proto", fileDescriptor_617387e490a04ffa) }

var fileDescriptor_617387e490a04ffa = []byte{
	// 426 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xad, 0x5b, 0x5a, 0xc8, 0xa4, 0x69, 0x92, 0x45, 0x80, 0x95, 0x53, 0x14, 0x84, 0xd4, 0x93,
	0x0d, 0x45, 0x1c, 0x2a, 0x71, 0xe0, 0x23, 0x91, 0x8a, 0x44, 0x85, 0xb4, 0x88, 0x73, 0xb4, 0x5d,
	0x8f, 0x92, 0x55, 0xbd, 0xbb, 0x66, 0x67, 0x0d, 0xf4, 0xc4, 0xff, 0xe0, 0x57, 0xf1, 0x0b, 0xf8,
	0x2d, 0x28, 0x6b, 0x3b, 0x76, 0x11, 0xe4, 0xc0, 0xc9, 0x9a, 0xb7, 0xef, 0xbd, 0x99, 0x37, 0x1e,
	0x18, 0x66, 0xc2, 0x8b, 0x22, 0x17, 0x06, 0x93, 0xc2, 0x59, 0x6f, 0x59, 0x6f, 0x0b, 0x4c, 0x16,
	0x2b, 0xe5, 0xd7, 0xe5, 0x55, 0x22, 0xad, 0x4e, 0x73, 0xb5, 0x12, 0xde, 0xa6, 0x06, 0xfd, 0x57,
	0xeb, 0xae, 0x09, 0xdd, 0x17, 0x25, 0x51, 0x23, 0xad, 0xd3, 0xe2, 0x7a, 0x95, 0x1a, 0xd2, 0xa9,
	0x28, 0x14, 0xa5, 0xd2, 0x6a, 0x6d, 0x4d, 0xfd, 0xa9, 0x1c, 0x67, 0x3f, 0x22, 0x18, 0x5e, 0xa2,
	0x5c, 0x0b, 0xa3, 0x48, 0x7f, 0x2a, 0x32, 0xe1, 0x91, 0xcd, 0x61, 0xec, 0x50, 0x5b, 0x8f, 0x4b,
	0xdd, 0xbc, 0x50, 0x1c, 0x4d, 0x0f, 0x4e, 0xfb, 0x67, 0x8f, 0x92, 0x5a, 0xcd, 0x03, 0x61, 0xab,
	0xe4, 0x23, 0x77, 0x1b, 0x20, 0xf6, 0x1a, 0x46, 0xb9, 0x95, 0x22, 0xef, 0x9a, 0xec, 0x07, 0x93,
	0x87, 0x8d, 0xc9, 0xfb, 0xcd, 0x7b, 0xeb, 0x31, 0xcc, 0x6f, 0xd5, 0x34, 0xfb, 0x15, 0x01, 0xbc,
	0xb5, 0xc6, 0xa0, 0xf4, 0xca, 0x1a, 0x76, 0x0e, 0xc7, 0x95, 0x23, 0xd9, 0xd2, 0x49, 0x8c, 0xa3,
	0x69, 0xb4, 0xc3, 0xad, 0x1f, 0xb8, 0x1f, 0x03, 0x95, 0x25, 0x70, 0x18, 0xca, 0x78, 0x7f, 0x97,
	0xe6, 0x62, 0x8f, 0x57, 0x34, 0xf6, 0x0c, 0x8e, 0xaa, 0x40, 0xf1, 0x41, 0x10, 0xfc, 0x2b, 0xf7,
	0xc5, 0x1e, 0xaf, 0x89, 0xec, 0x31, 0x0c, 0xe4, 0x76, 0xd6, 0xa5, 0xca, 0xe2, 0x3b, 0xd3, 0xe8,
	0xb4, 0xc7, 0x8f, 0x5b, 0xf0, 0x5d, 0xf6, 0x66, 0x00, 0xfd, 0x0c, 0xc9, 0x2b, 0x23, 0x36, 0xc0,
	0xec, 0x3b, 0x1c, 0x72, 0x2c, 0xf2, 0x1b, 0x16, 0xc3, 0x5d, 0x2a, 0xa5, 0x44, 0xa2, 0x90, 0xea,
	0x1e, 0x6f, 0x4a, 0xf6, 0x04, 0x4e, 0xf0, 0x9b, 0x47, 0x93, 0x61, 0xb6, 0x44, 0xe7, 0xac, 0x0b,
	0x11, 0x7a, 0x7c, 0xd0, 0xa0, 0x8b, 0x0d, 0xc8, 0x5e, 0x00, 0xb4, 0x8d, 0xea, 0xa1, 0x1f, 0x24,
	0xed, 0xfd, 0xb4, 0x6b, 0xe4, 0x1d, 0xe2, 0xd9, 0xcf, 0x08, 0xee, 0xcf, 0x1b, 0xd2, 0x87, 0x02,
	0x5d, 0x98, 0x8b, 0xd8, 0x2b, 0x18, 0x5f, 0x5a, 0xa3, 0xbc, 0x75, 0x9d, 0x3f, 0x3a, 0x68, 0x96,
	0xb0, 0xd0, 0x85, 0xbf, 0x99, 0x4c, 0x3a, 0xf6, 0x7f, 0x9c, 0xd0, 0xd3, 0x88, 0x9d, 0xc3, 0x49,
	0xdd, 0x93, 0xe3, 0xe7, 0x12, 0xc9, 0xb3, 0xbf, 0x8f, 0x33, 0x19, 0x75, 0xe0, 0x6a, 0x19, 0x2f,
	0x61, 0x3c, 0x57, 0x24, 0xff, 0x4f, 0x7d, 0x75, 0x14, 0x0e, 0xfb, 0xf9, 0xef, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x9e, 0xbf, 0x8b, 0x75, 0x3d, 0x03, 0x00, 0x00,
}
