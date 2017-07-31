package checkers

type PlayerStruct struct {
	*PlayerMLStruct
}

func NewPlayer(color Color, stateParams []float64, learningRate float64) PlayerML {

	mlPlayer := NewMLPlayer(color, stateParams, learningRate).(*PlayerMLStruct)
	return &PlayerStruct{mlPlayer}
}

func (p *PlayerStruct) EvaluateState(stateError float64, state GameState) {

	//do nothing for player struct
}
