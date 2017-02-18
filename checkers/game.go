package checkers

type Position struct {
	X int
	Y int
}

type Move struct {
	From Position
	To   Position
}

type Moves []Move

func (m Moves) clone() Moves {

	cloned := make(Moves, 0)

	for _, el := range m {
		cloned = append(cloned, el)
	}

	return cloned
}

func getOpponent(color Color) Color {

	if color == BLACK {
		return WHITE
	}

	return BLACK
}