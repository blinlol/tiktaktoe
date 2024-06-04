package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "tiktaktoe/game_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FieldType [3][3]pb.Player

type Game struct {
	Id int
	Field FieldType
    Iam pb.Player
}

var game Game


func (field *FieldType) Print(){
    for _, row := range field {
        fmt.Println(row)
    }
}

func GetRowCol() (row, col int32, err error) {
    fmt.Print("Enter row: ")
    _, err = fmt.Scanf("%d", &row)
    if err != nil {
        log.Printf("\nFail read row: %v\n", err)
        return
    }

    fmt.Print("Enter col: ")
    _, err = fmt.Scanf("%d", &col)
    if err != nil {
        log.Printf("\nFail read col: %v\n", err)
        return
    }

    return
}

func (game *Game) SendMove(stream pb.Game_MakeMoveClient) error {
    row, col, err := GetRowCol()
    if err != nil {
        return err
    }
    move := pb.Move{Row: row, Col: col, Who: game.Iam}
    err = stream.Send(&move)
    if err != nil {
        log.Printf("Fail to send: %v\n", err)
        return err
    }
    return nil
}

func (field *FieldType) ApplyMove(move *pb.Move) {
    fmt.Println(move.Message)

    row, col := move.Row, move.Col
    if field[row][col] != pb.Player_NONE {
        return
    }
    if row < 0 || col < 0 {
        return
    }
    field[row][col] = move.Who
    field.Print()
}

func Play(ctx context.Context, client pb.GameClient) error {
    response, err := client.StartGame(ctx, &pb.StartRequest{})
    if err != nil {
        return err
    }
    log.Printf("Play::response = %s\n", response.String())
    game.Field.Print()
    game.Id = int(response.GetGameId())
    game.Iam = response.GetIam()

    stream, err := client.MakeMove(ctx)
    if err != nil {
        log.Fatalf("client.RouteChat failed: %v", err)
    }

    err = stream.Send(&pb.Move{Row: -1, Col: -1, Who: game.Iam})
    if err != nil {
        log.Fatalf("%v\n", err)
    }

    if response.Iam == pb.Player_CROSS {
        game.SendMove(stream)
    }
    for {
        move, err := stream.Recv()
        if err == io.EOF {
            log.Println("Game end!")
            break
        } else if err != nil {
            log.Printf("%v\n", err)
            continue
        }
        game.Field.ApplyMove(move)
        if move.Who != game.Iam {
            game.SendMove(stream)
        }
    }
    return nil
}

func main(){
    conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Didn't connenct: %v\n", err)
    }
    defer conn.Close()
    client := pb.NewGameClient(conn)

    ctx := context.Background()

    err = Play(ctx, client)
    if err != nil {
        log.Fatalf("%v\n", err)
    }
}