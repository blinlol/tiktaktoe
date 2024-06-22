package game


import (
	"errors"
	"fmt"
	"log"

	pb "tiktaktoe/game_proto"

)


type GameField [3][3]pb.Player


type Game struct {
	Id int32
	Field GameField
	Streams map[pb.Player]pb.Game_GameStreamServer
	WaitChan chan int
	NextMovePlayer pb.Player
}


var playerToString map[pb.Player]string = map[pb.Player]string{
    pb.Player_ZERO: "0",
    pb.Player_CROSS: "X",
    pb.Player_NONE: "-",
}


func (field *GameField) Print(){
    for _, row := range field {
        for _, symb := range row {
            fmt.Print(playerToString[symb])
        }
        fmt.Println()
    }
}

func (f *GameField) check(row, col int32) (end bool, winner pb.Player) {
	if f[row][0] == f[row][1] && f[row][1] == f[row][2] {
		return true, f[row][col]
	}
	if f[0][col] == f[1][col] && f[1][col] == f[2][col] {
		return true, f[row][col]
	}
	if f[0][0] == f[1][1] && f[1][1] == f[2][2] && f[0][0] != pb.Player_NONE {
		return true, f[0][0]
	}
	if f[0][2] == f[1][1] && f[1][1] == f[2][0] && f[2][0] != pb.Player_NONE {
		return true, f[0][0]
	}

	// draft game
	for _, row := range f {
		for _, el := range row {
			if el == pb.Player_NONE {
				return
			}
		}
	}
	end = true
	return
}

func (f *GameField) IsMoveValid(row, col int32) bool {
    if row < 0 || row > 2 {
		return false
    }
	if col < 0 || col > 2 {
		return false
    }
    if f[row][col] != pb.Player_NONE {
		return false
    }
	return true
}

func (field *GameField) ApplyMove(row, col int32, who pb.Player) (is_end bool, winner pb.Player, err error) {
	// is_end == true: someone win or draft game
	// winner == pb.Player && is_end: draft game

    if row < 0 || row > 2 {
		err = errors.New("wrong row")
        return 
    }
	if col < 0 || col > 2 {
		err = errors.New("wrong col")
        return 
    }
    if field[row][col] != pb.Player_NONE {
		err = errors.New("cell already taken")
        return 
    }

    field[row][col] = who
	is_end, winner = field.check(row, col)
	return
}


func (game *Game) InitFirstPlayer(stream pb.Game_GameStreamServer) (pb.Player, error) {
	who := pb.Player_CROSS
	game.Streams[who] = stream
	<- game.WaitChan

	message := pb.GameServerMessage{
		Message: "Hello!",
		GameId: game.Id,
		Type: pb.ServerMessageType_INIT_PLAYER,
		Who: who}
	
	return who, stream.Send(&message)
}

func (game *Game) InitSecondPlayer(stream pb.Game_GameStreamServer) (pb.Player, error) {
	who := pb.Player_ZERO
	game.Streams[who] = stream
	close(game.WaitChan)

	message := pb.GameServerMessage{
		Message: "Hello!",
		GameId: game.Id,
		Type: pb.ServerMessageType_INIT_PLAYER,
		Who: who}
	return who, stream.Send(&message)
}

func (game *Game) ApplyMove(row, col int32) (end bool, winner pb.Player, err error) {
	log.Println("in game.ApplyMove")
	end, winner, err = game.Field.ApplyMove(row, col, game.NextMovePlayer)
	if err != nil {
		return
	}

	m := pb.GameServerMessage{
		GameId: game.Id,
		Type: pb.ServerMessageType_SAVE_MOVE,
		Who: game.NextMovePlayer,
		Row: row,
		Col: col,
	}
	log.Println(&m)
	err = game.Streams[pb.Player_CROSS].Send(&m)
	if err != nil {
		return
	}
	err = game.Streams[pb.Player_ZERO].Send(&m)

	if game.NextMovePlayer == pb.Player_CROSS {
		game.NextMovePlayer = pb.Player_ZERO
	} else {
		game.NextMovePlayer = pb.Player_CROSS
	}
	return
}