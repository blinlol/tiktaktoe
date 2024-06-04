package main

import (
	"fmt"
	"sync"
)

type Field [][]string

type Response struct {
	win bool
	who string
}

var Zero string = "0"
var Cross string = "X"
var Space string = "-"

var wg sync.WaitGroup

func CreateField() Field {
	f := make(Field, 3)
	for i := 0; i < 3; i++ {
		f[i] = []string{Space, Space, Space}
	}
	return f
}

func checkRow(out chan <- Response, field *Field, row int) {
	defer wg.Done()
	who := (*field)[row][0]
	if who != Space && (*field)[row][1] == (*field)[row][2] && (*field)[row][1] == who {
		out <- Response{win: true, who: who}
		close(out)
	}
}

func checkCol(out chan <- Response, field *Field, col int) {
	defer wg.Done()
	who := (*field)[0][col]
	if who != Space && (*field)[1][col] == (*field)[2][col] && (*field)[1][col] == who {
		out <- Response{win: true, who: who}
		close(out)
	}
}

// return is anyone win and who win
func Check(field *Field) Response {
	wg.Add(6)
	responses := make(chan Response, 9)
	for i := 0; i < 3; i++ {
		go checkCol(responses, field, i)
		go checkRow(responses, field, i)
	}
	wg.Wait()

	res := Response{win: false, who: ""}
	select {
		case res = <- responses:
		default:
	}
	return res
}

func (field *Field) Print() {
	for _, row := range *field {
		for _, s := range row {
			fmt.Print(s)
		}
		fmt.Println()
	}
}