// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: game_proto/game.proto

package game_proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GameClient is the client API for Game service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameClient interface {
	// клиенты подключаются к серверу, сервер создает игру.
	// клиенты посылают свои желаемые ходы, а сервер подтверждает,
	// что ход возможен (пересылает этот же ход)
	// и затем посылает ход следующего игрока.
	// т.е. при получении хода сервер его выполняет, и посылает обратно обоим клиентам.
	StartGame(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error)
	MakeMove(ctx context.Context, opts ...grpc.CallOption) (Game_MakeMoveClient, error)
	Test(ctx context.Context, in *MessageTime, opts ...grpc.CallOption) (*MessageTime, error)
	// двунаправленный поток общения клиента и сервера, отвечает за
	// 1. создание и подключение к игре
	// 2. определение кто за кого играет
	// 3. отправка сервером информации кто ходит следующим
	// 4. отправка клиентом хода и получение его на сервере
	// валидация хода происходит на сервере и клиенте, отвечает за это отдельный пакет
	GameStream(ctx context.Context, opts ...grpc.CallOption) (Game_GameStreamClient, error)
}

type gameClient struct {
	cc grpc.ClientConnInterface
}

func NewGameClient(cc grpc.ClientConnInterface) GameClient {
	return &gameClient{cc}
}

func (c *gameClient) StartGame(ctx context.Context, in *StartRequest, opts ...grpc.CallOption) (*StartResponse, error) {
	out := new(StartResponse)
	err := c.cc.Invoke(ctx, "/Game/StartGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) MakeMove(ctx context.Context, opts ...grpc.CallOption) (Game_MakeMoveClient, error) {
	stream, err := c.cc.NewStream(ctx, &Game_ServiceDesc.Streams[0], "/Game/MakeMove", opts...)
	if err != nil {
		return nil, err
	}
	x := &gameMakeMoveClient{stream}
	return x, nil
}

type Game_MakeMoveClient interface {
	Send(*Move) error
	Recv() (*Move, error)
	grpc.ClientStream
}

type gameMakeMoveClient struct {
	grpc.ClientStream
}

func (x *gameMakeMoveClient) Send(m *Move) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gameMakeMoveClient) Recv() (*Move, error) {
	m := new(Move)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gameClient) Test(ctx context.Context, in *MessageTime, opts ...grpc.CallOption) (*MessageTime, error) {
	out := new(MessageTime)
	err := c.cc.Invoke(ctx, "/Game/Test", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameClient) GameStream(ctx context.Context, opts ...grpc.CallOption) (Game_GameStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Game_ServiceDesc.Streams[1], "/Game/GameStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &gameGameStreamClient{stream}
	return x, nil
}

type Game_GameStreamClient interface {
	Send(*GameClientMessage) error
	Recv() (*GameServerMessage, error)
	grpc.ClientStream
}

type gameGameStreamClient struct {
	grpc.ClientStream
}

func (x *gameGameStreamClient) Send(m *GameClientMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gameGameStreamClient) Recv() (*GameServerMessage, error) {
	m := new(GameServerMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GameServer is the server API for Game service.
// All implementations must embed UnimplementedGameServer
// for forward compatibility
type GameServer interface {
	// клиенты подключаются к серверу, сервер создает игру.
	// клиенты посылают свои желаемые ходы, а сервер подтверждает,
	// что ход возможен (пересылает этот же ход)
	// и затем посылает ход следующего игрока.
	// т.е. при получении хода сервер его выполняет, и посылает обратно обоим клиентам.
	StartGame(context.Context, *StartRequest) (*StartResponse, error)
	MakeMove(Game_MakeMoveServer) error
	Test(context.Context, *MessageTime) (*MessageTime, error)
	// двунаправленный поток общения клиента и сервера, отвечает за
	// 1. создание и подключение к игре
	// 2. определение кто за кого играет
	// 3. отправка сервером информации кто ходит следующим
	// 4. отправка клиентом хода и получение его на сервере
	// валидация хода происходит на сервере и клиенте, отвечает за это отдельный пакет
	GameStream(Game_GameStreamServer) error
	mustEmbedUnimplementedGameServer()
}

// UnimplementedGameServer must be embedded to have forward compatible implementations.
type UnimplementedGameServer struct {
}

func (UnimplementedGameServer) StartGame(context.Context, *StartRequest) (*StartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartGame not implemented")
}
func (UnimplementedGameServer) MakeMove(Game_MakeMoveServer) error {
	return status.Errorf(codes.Unimplemented, "method MakeMove not implemented")
}
func (UnimplementedGameServer) Test(context.Context, *MessageTime) (*MessageTime, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Test not implemented")
}
func (UnimplementedGameServer) GameStream(Game_GameStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GameStream not implemented")
}
func (UnimplementedGameServer) mustEmbedUnimplementedGameServer() {}

// UnsafeGameServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameServer will
// result in compilation errors.
type UnsafeGameServer interface {
	mustEmbedUnimplementedGameServer()
}

func RegisterGameServer(s grpc.ServiceRegistrar, srv GameServer) {
	s.RegisterService(&Game_ServiceDesc, srv)
}

func _Game_StartGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).StartGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Game/StartGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).StartGame(ctx, req.(*StartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_MakeMove_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GameServer).MakeMove(&gameMakeMoveServer{stream})
}

type Game_MakeMoveServer interface {
	Send(*Move) error
	Recv() (*Move, error)
	grpc.ServerStream
}

type gameMakeMoveServer struct {
	grpc.ServerStream
}

func (x *gameMakeMoveServer) Send(m *Move) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gameMakeMoveServer) Recv() (*Move, error) {
	m := new(Move)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Game_Test_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageTime)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameServer).Test(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Game/Test",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameServer).Test(ctx, req.(*MessageTime))
	}
	return interceptor(ctx, in, info, handler)
}

func _Game_GameStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GameServer).GameStream(&gameGameStreamServer{stream})
}

type Game_GameStreamServer interface {
	Send(*GameServerMessage) error
	Recv() (*GameClientMessage, error)
	grpc.ServerStream
}

type gameGameStreamServer struct {
	grpc.ServerStream
}

func (x *gameGameStreamServer) Send(m *GameServerMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gameGameStreamServer) Recv() (*GameClientMessage, error) {
	m := new(GameClientMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Game_ServiceDesc is the grpc.ServiceDesc for Game service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Game_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Game",
	HandlerType: (*GameServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartGame",
			Handler:    _Game_StartGame_Handler,
		},
		{
			MethodName: "Test",
			Handler:    _Game_Test_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MakeMove",
			Handler:       _Game_MakeMove_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "GameStream",
			Handler:       _Game_GameStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "game_proto/game.proto",
}
