package main

import (
	"fmt"
)

func main() {
	field := CreateField()
	Move(&field, 1, 1, Cross)
	Move(&field, 2, 1, Zero)
	Move(&field, 2, 0, Zero)
	Move(&field, 2, 2, Zero)
	field.Print()
	fmt.Println(Check(&field))
}
