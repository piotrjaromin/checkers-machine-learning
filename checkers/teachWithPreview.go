package checkers

import (
	"fmt"
	"log"
	"os"
)

func TeachWithPreview(drawer Draw, learningRate float64) {

	logger := log.New(os.Stdout, fmt.Sprintf("[%s]", "[MAIN]"), log.LstdFlags)
	board := NewBoard(BLACK, logger)

	p1 := NewMLPlayer(BLACK, RandomParams(), learningRate)
	p2 := NewMLPlayer(WHITE, RandomParams(), learningRate)

	state := PlayOneGame(board, p1, p2, drawer)

	logger.Printf("%+v", state)

}
