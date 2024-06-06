package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "tiktaktoe/game_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (field *FieldType) check() (bool, pb.Player) {
	for i := 0 ; i < 3; i++ {
		if field[i][0] != pb.Player_NONE && field[i][0] == field[i][1] && field[i][1] == field[i][2]{
			return true, field[i][0]
		}
		if field[0][i] != pb.Player_NONE && field[0][i] == field[1][i] && field[1][i] == field[2][i] {
			return true, field[0][i]
		}
	}
	if field[0][0] != pb.Player_NONE && field[0][0] == field[1][1] && field[1][1] == field[2][2] {
		return true, field[0][0]
	}
	if field[0][2] != pb.Player_NONE && field[0][2] == field[1][1] && field[1][1] == field[2][0] {
		return true, field[0][2]
	}

	for _, row := range field {
		for _, el := range row {
			if el == pb.Player_NONE {
				return false, pb.Player_NONE
			}
		}
	}
	return true, pb.Player_NONE
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
		<-game.waitCross
	}()
	go func(){
		defer wg.Done()
		<-game.waitZero
	}()

	move, err := stream.Recv()
	if err != nil {
		log.Fatalf("InitPlayers: %v\n", err)
	}

	if move.Who == pb.Player_CROSS {
		log.Println("init cross stream")
		game.cross_stream = stream
		close(game.waitCross)
	} else {
		log.Println("init zero stream")
		game.zero_stream = stream
		close(game.waitZero)
	}
	wg.Wait()
}

func (game *Game) ApplyMove(move *pb.Move) (bool, error) {
	fmt.Println(move.Message)

	row, col := move.Row, move.Col
	win := false
	if game.Field[row][col] != pb.Player_NONE {
		// нельзя поставить, значит посылаем запрос автору на переделку
		move.Row, move.Col = -1, -1
		move.Message = "Wrong row, col"
		if move.Who == pb.Player_CROSS {
			move.Who = pb.Player_ZERO
			game.cross_stream.Send(move)
		} else {
			move.Who = pb.Player_CROSS
			game.zero_stream.Send(move)
		}
		return false, nil
	} else {
		game.Field[row][col] = move.Who
		move.Finish, move.Winner = game.Field.check()
		win = move.Finish
	}

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
	return win, nil
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
		} else if status.Code(err) == codes.Canceled {
			break
		} else if err != nil {
			log.Fatalf("Recv error: %v\n", err)
			continue
		}

		finish, err := s.game.ApplyMove(move)
		if finish {
			break
		}
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