package main

import (
	"github.com/piotrjaromin/machinelearning/checkers"
	"fmt"
)

func main() {

	board := checkers.NewBoard(checkers.BLACK)

	fmt.Printf(board.String())

	fmt.Print("\n===============\n")

	for i := 0; i < 10; i++ {
		p := board.GetFields()[i][0]
		if p != nil {
			switch p.Color {
			case checkers.BLACK:
				fmt.Printf(" B ")
			case checkers.WHITE:
				fmt.Printf(" W ")
			}
		} else {
			fmt.Printf(" _ ")
		}
	}
}
