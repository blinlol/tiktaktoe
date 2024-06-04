package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"sync"

	pb "tiktaktoe/game_proto"

	"google.golang.org/grpc"
)


type FieldType [3][3]pb.Player

type Game struct {
	Id int
	Field FieldType
	cross_stream pb.Game_MakeMoveServer
	zero_stream pb.Game_MakeMoveServer
	waitZero chan int
	waitCross chan int
}

type server struct {
	pb.UnimplementedGameServer

	game Game
}


func (game *Game) InitPlayers(stream pb.Game_MakeMoveServer){
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func(){
		defer wg.Done()
		log.Println(1)
		<-game.waitCross
		log.Println(11)
	}()
	go func(){
		defer wg.Done()
		log.Println(2)
		<-game.waitZero
		log.Println(22)
	}()

	move, err := stream.Recv()
	if err != nil {
		log.Fatalf("InitPlayers: %v\n", err)
	}

	if move.Who == pb.Player_CROSS {
		log.Println("cross stream")
		game.cross_stream = stream
		close(game.waitCross)
	} else {
		log.Println("zero stream")
		game.zero_stream = stream
		close(game.waitZero)
	}
	wg.Wait()
}


func (game *Game) ApplyMove(move *pb.Move) error {
	fmt.Println(move.Message)
	row, col := move.Row, move.Col

	if game.Field[row][col] != pb.Player_NONE {
		log.Fatalln("Wrong row col")
		return nil
	}
	game.Field[row][col] = move.Who
	
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(){
		defer wg.Done()
		game.cross_stream.Send(move)
	}()

	go func(){
		defer wg.Done()
		game.zero_stream.Send(move)
	}()

	wg.Wait()
	return nil
}

func (*server) Test(ctx context.Context, in *pb.MessageTime) (*pb.MessageTime, error) {
	log.Printf("server.Test: %s\n", in.String())
	answer := pb.MessageTime{
		Message: "from server",
		Time: time.Now().String()}
	return &answer, nil
}

func (s *server) StartGame(ctx context.Context, request *pb.StartRequest) (*pb.StartResponse, error) {
	if s.game.waitCross == nil {
		s.game.waitZero = make(chan int)
		s.game.waitCross = make(chan int)
	}
	var response pb.StartResponse
	if s.game.Id == 0 {
		response.Iam = pb.Player_ZERO
		response.GameId = 1
		s.game.Id = 1
	} else {
		response.Iam = pb.Player_CROSS
		response.GameId = 1
	}
	log.Printf("StartGame response: %s", response.String())
	return &response, nil
}

func (s *server) MakeMove(stream pb.Game_MakeMoveServer) error {
	s.game.InitPlayers(stream)

	for {
		move, err := stream.Recv()
		if err == io.EOF {
			log.Println("End Game")
			break
		} else if err != nil {
			log.Printf("Recv error: %v\n", err)
			continue
		}

		err = s.game.ApplyMove(move)
		if err != nil {
			log.Printf("ApplyMove %v\n", err)
		}
	}
	return nil
}

func main(){
	lis, err := net.Listen("tcp", ":50052")
	if err != nil{
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterGameServer(s, &server{})
	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Serve error: %v\n", err)
	}
}