package main

import "fmt"

type Field [][]string

var Zero string = "0"
var Cross string = "X"
var Space string = "-"

func CreateField() Field {
	f := make(Field, 3)
	for i := 0; i < 3; i++ {
		f[i] = []string{Space, Space, Space}
	}
	return f
}

func Move(field *Field, x, y int, val string) {
	(*field)[y][x] = val
	// checks
}

func checkRow(field *Field, row int) (bool, string) {
	who := (*field)[row][0]
	if who != Space && (*field)[row][1] == (*field)[row][2] && (*field)[row][1] == who {
		return true, who
	}
	return false, ""
}

func checkCol(field *Field, col int) (bool, string) {
	who := (*field)[0][col]
	if who != Space && (*field)[1][col] == (*field)[2][col] && (*field)[1][col] == who {
		return true, who
	}
	return false, ""
}

// return is anyone win and who win
func Check(field *Field) (bool, string) {
	for i := 0; i < 3; i++ {
		win, who := checkCol(field, i)
		if win {
			return win, who
		}

		win, who = checkRow(field, i)
		if win {
			return win, who
		}
	}
	return false, ""
}

func (field *Field) Print() {
	for _, row := range *field {
		for _, s := range row {
			fmt.Print(s)
		}
		fmt.Println()
	}
}