package checkers

import (
	"log"
	"math/rand"
	"time"
)

const STATE_PARAMS = 7

type PlayerML interface {
	GetMoves(b Board) Moves
	EvaluateState(stateError float64, state GameState)
	GetState() []float64
	Color() Color
}

type PlayerMLStruct struct {
	color        Color
	stateParams  []float64
	LearningRate float64
}

func NewMLPlayer(color Color, stateParams []float64, learningRate float64) PlayerML {

	return &PlayerMLStruct{color, stateParams, learningRate}
}

func (p PlayerMLStruct) GetMoves(b Board) Moves {

	allMoves := b.GetAllMovesFor(p.color)

	if len(allMoves) == 0 {
		return Moves{}
	}

	var highestStateValue *float64
	var bestMove Moves
	for _, moves := range allMoves {
		tempBoard := b.Clone()

		if tempBoard.Move(moves) {
			gameResultAfterMove := tempBoard.GameResult()
			if gameResultAfterMove.isWin(p.color) {
				bestMove = moves
				break
			}

			gameState := calculateState(tempBoard)
			stateValue := calculateStateValue(PlayerML(&p), gameState, gameResultAfterMove)
			if highestStateValue == nil || stateValue >= *highestStateValue {
				highestStateValue = &stateValue
				bestMove = moves
			}
		} else {
			log.Println("Could not perform move for player ", p.color)
		}
	}

	return bestMove
}

func (p PlayerMLStruct) GetState() []float64 {
	return p.stateParams
}

func (p PlayerMLStruct) Color() Color {
	return p.color
}

func (p *PlayerMLStruct) EvaluateState(stateError float64, state GameState) {

	//we are omitting first parameter because it is treated as constant  eg. stateVar*x^0 which is stateVar
	p.stateParams[0] += p.LearningRate * stateError
	stateLen := len(p.stateParams)
	for i := 1; i < stateLen; i++ {
		p.stateParams[i] += p.LearningRate * stateError * float64(state.ToArray()[i-1])
	}
}

func RandomParams() []float64 {

	rand.Seed(time.Now().UnixNano())

	params := make([]float64, STATE_PARAMS, STATE_PARAMS)
	for i := 0; i < STATE_PARAMS; i++ {
		params[i] = rand.Float64() * 1
	}

	return params
}
