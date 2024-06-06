package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "tiktaktoe/game_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/credentials/insecure"
)

type FieldType [3][3]pb.Player

type Game struct {
	Id int
	Field FieldType
    Iam pb.Player
}

var game Game

var playerToString map[pb.Player]string = map[pb.Player]string{
    pb.Player_ZERO: "0",
    pb.Player_CROSS: "X",
    pb.Player_NONE: "-",
}

func (field *FieldType) Print(){
    for _, row := range field {
        for _, symb := range row {
            fmt.Print(playerToString[symb])
        }
        fmt.Println()
    }
}

func GetRowCol() (row, col int32) {
    var err error

    fmt.Print("Enter row: ")
    _, err = fmt.Scanf("%d", &row)
    for ; err != nil || 0 > row || row > 2 ; {
        fmt.Printf("\nFail read row: %v\n", err)
        _, err = fmt.Scanf("%d", &row)
    }

    fmt.Print("Enter col: ")
    _, err = fmt.Scanf("%d", &col)
    for ; err != nil || 0 > col || col > 2 ; {
        fmt.Printf("\nFail read col: %v\n", err)
        _, err = fmt.Scanf("%d", &col)
    }

    return
}

func (game *Game) SendMove(stream pb.Game_MakeMoveClient) error {
    row, col := GetRowCol()
    move := pb.Move{Row: row, Col: col, Who: game.Iam}
    err := stream.Send(&move)
    if err != nil {
        log.Printf("Fail to send: %v\n", err)
        return err
    }
    return nil
}

func (field *FieldType) ApplyMove(move *pb.Move) {
    fmt.Println(move.Message)

    row, col := move.Row, move.Col

    if row < 0 || col < 0 {
        return
    }
    if field[row][col] != pb.Player_NONE {
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
    game.Id = int(response.GetGameId())
    game.Iam = response.GetIam()

    stream, err := client.MakeMove(ctx)
    if err != nil {
        log.Fatalf("client.RouteChat failed: %v", err)
    }

    // first message only for synchronization
    err = stream.Send(&pb.Move{Row: -1, Col: -1, Who: game.Iam, Message: "Want to play"})
    if err != nil {
        log.Printf("%v\n", err)
        return err
    }

    game.Field.Print()
    if response.Iam == pb.Player_CROSS {
        err = game.SendMove(stream)
        if err != nil {
            log.Printf("Failed to send move%v\n", err)
            return err
        }
    }
    for {
        move, err := stream.Recv()
        if status.Code(err) == codes.Unavailable {
            log.Printf("%v\n", err)
            return err
        }
        if err == io.EOF {
            log.Println("Game end!")
            fmt.Println(move.Message)
            break
        }
        game.Field.ApplyMove(move)
        if move.Finish {
            fmt.Printf("Winner: %s\n", move.Winner)
            break
        }
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
        log.Fatalf("main: %v\n", err)
    }
}


// сделать так, чтобы начинали играть и выводить поле только после того, как нашелся второй игрок