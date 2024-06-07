package game


import (
	"errors"
	"fmt"

	pb "tiktaktoe/game_proto"
)


type GameField [3][3]pb.Player


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

func (f *GameField) check(row, col int) (end bool, winner pb.Player) {
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

func (field *GameField) ApplyMove(row, col int, who pb.Player) (is_end bool, winner pb.Player, err error) {
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
