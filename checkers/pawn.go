package checkers

type Color int

const (
        BLACK Color = 1
        WHITE Color = 0
)

type Direction int

const (
        UP Direction = 1
        DOWN Direction = -1
)

func (d Direction) Opposite() Direction {

        if d == UP {
                return DOWN
        }

        return UP
}

type Pawn struct {
        Color             Color
        MovementDirection Direction
        King              bool
}

func (c Color) String() string {
        if c == BLACK {
                return "black"
        }
        return "white"
}

func (self Pawn) isOpponent(other *Pawn) bool {

        if other != nil {
                return self.Color != other.Color
        }
        return false
}