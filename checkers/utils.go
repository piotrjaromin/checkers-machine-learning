package checkers


func appendMoves(totalMoves []Moves, moves Moves) []Moves {
	if len(moves) >0 {
		return append(totalMoves, moves)
	}

	return totalMoves
}