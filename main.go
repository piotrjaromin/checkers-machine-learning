package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/piotrjaromin/checkers-machine-learning/checkers"
)

func main() {

	mode := flag.String("mode", "preview", `determines mode of teach algorithm {teach|preview|play}
- teach - plays multiple games and teaches model
- preview - plays one game with randomized model parameters
- play - plays game on already trained model, requires playerParamsPath argument`)

	drawerName := flag.String("drawer", "donothing", `decides how board will be drawn {opengl|asci|donothing}
- donothing - do not draw board
- asci - draws asci board for game
- opengl - not implements, ui with graphic
	`)

	moveDelay := flag.Int("moveDelay", 400, "in preview mode enabled delay between moves [ms]")
	maxGamesPlayed := flag.Int("maxGames", 1000, "amount of games that should be played during teaching")
	learningRate := flag.Float64("learningRate", 0.00001, "learning rate speed, tells how fast model parameters should change")
	playersParamsPath := flag.String("playerParamsPath", "./params.json", "path to json file containing player parameters")

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
	case "play":
		endGameAfter := checkers.CreateEndAfterGamesPlayed(*maxGamesPlayed)
		config := readPlayerParams(*playersParamsPath)
		stats := checkers.PlayGames(endGameAfter, drawer, config.P1, config.P2)
		log.Println(stats.String())
	default:
		flag.PrintDefaults()
	}
}

func readPlayerParams(path string) playerConfig {

	file, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}

	var playerConfig playerConfig
	json.Unmarshal(file, &playerConfig)
	return playerConfig
}

func getDrawer(name string, delay int) checkers.Draw {

	switch name {
	case "asci":
		var drawer checkers.AsciDraw
		return drawer
	default:
		var emptyDrawer checkers.DoNothingDraw
		return emptyDrawer
	}
}

type playerConfig struct {
	P1 []float64 `json:"p1"`
	P2 []float64 `json:"p2"`
}
