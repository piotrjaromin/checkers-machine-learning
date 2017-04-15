package main

import (
	"flag"
	"log"

	"github.com/piotrjaromin/machinelearning/checkers"
)

func main() {

	mode := flag.String("mode", "preview", "determines mode of teach algorithm {teach|preview}")
	drawerName := flag.String("drawer", "donothing", "decides how board will be drawn {opengl|asci|donothing}")
	moveDelay := flag.Int("moveDelay", 400, "in preview mode enabled delay between moves [ms]")
	maxGamesPlayed := flag.Int("maxGames", 1000, "amount of games that should be played during teaching")
	learningRate := flag.Float64("learningRate", 0.0001, "learning rate speed")

	flag.Parse()

	log.Printf("selected mode is %s, drawerName is %s, moveDelay is %d", *mode, *drawerName, *moveDelay)

	drawer := getDrawer(*drawerName, *moveDelay)

	switch *mode {
	case "preview":
		checkers.TeachWithPreview(drawer, *learningRate)
	case "teach":
		endGameAfter := checkers.CreateEndAfterGamesPlayed(*maxGamesPlayed)
		stats := checkers.Teach(endGameAfter, drawer, *learningRate)
		log.Println(stats.String())
	default:
		flag.PrintDefaults()
	}
}


func getDrawer(name string, delay int) checkers.Draw {

	switch name {
	case "asci":
		var drawer checkers.AsciDraw
		return drawer
	case "opengl":
		return checkers.CreateOpenGlDraw(300, 300, 10,10, delay)
	default:
		var emptyDrawer checkers.DoNothingDraw
		return emptyDrawer
	}
}