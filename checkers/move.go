package checkers

import (
	"fmt"
)

type Position struct {
	X int
	Y int
}

func (p Position) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

type Move struct {
	From Position
	To   Position
}

func (m Move) String() string {
	return fmt.Sprintf("%s -> %s", m.From.String(), m.To.String())
}

type Moves []Move

func (m Moves) clone() Moves {

	cloned := make(Moves, 0)

	for _, el := range m {
		cloned = append(cloned, el)
	}

	return cloned
}

func (m Moves) String() string {
	str := ""
	for _, el := range m {
		str += fmt.Sprintf("%s | ", el.String())
	}
	return str
}

func getOpponent(color Color) Color {

	if color == BLACK {
		return WHITE
	}

	return BLACK
}
