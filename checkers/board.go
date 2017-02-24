package checkers

import (
        "log"
        "fmt"
        "reflect"
)

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

type GAME_RESULT int

const (
        DRAW GAME_RESULT = 0
        WIN_BLACK GAME_RESULT = 1
        WIN_WHITE GAME_RESULT = 2
        UNFINISHED GAME_RESULT = 3
)

const (
        MAX_X = 9
        MAX_Y = 9
        MIN_X = 0
        MIN_Y = 0
)

const (
        TO_LEFT = -1
        TO_RIGHT = 1
)

type Board struct {
        fields    [10][10]*Pawn
        topPlayer Color //player on top of board ( starting position closer to MAX_X, MAX_Y
}

func NewBoard(topColor Color) Board {

        b := Board{topPlayer: topColor}

        for y := 0; y < 4; y++ {

                for x := 0; x <= MAX_X; x ++ {
                        if x % 2 == 0 && y % 2 == 0 {
                                b.fields[x][y] = &Pawn{getOpponent(topColor), UP, false}
                                b.fields[x][y + 6] = &Pawn{topColor, DOWN, false}
                        } else if x % 2 == 1 && y % 2 == 1 {
                                b.fields[x][y] = &Pawn{getOpponent(topColor), UP, false}
                                b.fields[x][y + 6] = &Pawn{topColor, DOWN, false}
                        }
                }

        }

        return b
}

func (b*Board) Move(moves Moves) bool {

        if b.isInvalidMoves(moves) {
                return false
        }

        for _, move := range moves {

                b.fields[move.To.X][move.To.Y] = b.fields[move.From.X][move.From.Y]
                b.fields[move.From.X][move.From.Y] = nil
        }

        return true
}

func (b Board) isInvalidMoves(moves Moves) bool {

        validMoves := b.GetValidMovesForPosition(moves[0].From)

        for _, validMove := range validMoves {
                if reflect.DeepEqual(validMove, moves) {
                        return false
                }
        }
        return true
}

func (b Board) getPawn(p Position) *Pawn {
        return b.fields[p.X][p.Y]
}

func (b Board) GetValidMovesForPosition(position Position) []Moves {
        moves := []Moves{}

        if ! b.isOnBoard(position) {
                return moves
        }

        pawn := b.getPawn(position)

        if pawn == nil {
                log.Println("[GetValidMovesForPosition] There is no pawn at this position")
                return moves
        }

        //top bottom multipliers (up/down is more of player perspective)
        verticalMoves := func(verticalDirection Direction) {
                dir := int(verticalDirection)
                moves = appendMoves(moves, b.ifPossibleRightUp(position, TO_RIGHT, dir))
                moves = appendMoves(moves, b.ifPossibleLeftUp(position, TO_LEFT, dir))
                moves = append(moves, b.moveOverOpponent(position, make(Moves, 0), TO_RIGHT, TO_LEFT, dir, pawn.Color)...)
        }

        verticalMoves(pawn.MovementDirection)
        if pawn.King {
                verticalMoves(pawn.MovementDirection.Opposite())
        }

        return moves
}

func (b Board) ifPossibleRightUp(p Position, toRight int, oneUp int) Moves {
        moves := Moves{}

        if ! b.canMove(p, toRight, oneUp) {
                return moves
        }

        getUpRight := func(p Position) *Pawn {
                return b.fields[p.X + 1][ p.Y + oneUp]
        }

        if getUpRight(p) == nil {
                moves = append(moves, Move{From: p, To: p.Add(toRight, oneUp)})
        }
        return moves
}

func (b Board) ifPossibleLeftUp(p Position, toLeft int, oneUp int) Moves {
        moves := Moves{}

        if ! b.canMove(p, toLeft, oneUp) {
                return moves
        }

        getUpLeft := func(p Position) *Pawn {
                return b.fields[p.X - 1][p.Y + oneUp]
        }

        if getUpLeft(p) == nil {
                moves = append(moves, Move{From: p, To: p.Add(toLeft, oneUp)})
        }
        return moves
}

func (b Board) moveOverOpponent(position Position, currentMoves Moves, toRight int, toLeft int, oneUp int, color Color) []Moves {
        totalMoves := make([]Moves, 0, 2)
        rightUp := b.getUpRightPawn(position, oneUp)

        if rightUp != nil && rightUp.Color == getOpponent(color) {

                //this is position of opponent
                newRightUp := position.Add(toRight, oneUp)
                //so this will add jumped over position
                jumpOver := b.ifPossibleRightUp(newRightUp, toRight, oneUp)
                //we can jump over
                if len(jumpOver) > 0 {
                        // add as partial move
                        currentMoves := append(currentMoves, Move{position, jumpOver[0].To})

                        totalMoves = append(totalMoves, currentMoves.clone())

                        //continue(chain jumps)
                        totalMoves = append(totalMoves, b.moveOverOpponent(jumpOver[0].To, currentMoves, toRight, toLeft, oneUp, color)...)
                }
        }

        leftUp := b.getUpLeftPawn(position, oneUp)
        if leftUp != nil && leftUp.Color == getOpponent(color) {
                newLeftUp := position.Add(toLeft, oneUp)

                jumpOver := b.ifPossibleLeftUp(newLeftUp, toLeft, oneUp)

                if len(jumpOver) > 0 {
                        currentMoves := append(currentMoves, Move{position, jumpOver[0].To})

                        totalMoves = append(totalMoves, currentMoves.clone())
                        totalMoves = append(totalMoves, b.moveOverOpponent(jumpOver[0].To, currentMoves, toRight, toLeft, oneUp, color)...)
                }

        }

        return totalMoves
}

func (b Board) canMove(position Position, direction int, oneUp int) bool {

        return b.isOnBoard(position.Add(direction, oneUp))
}

func (b Board) isOnBoard(position Position) bool {

        if (position.Y > MAX_Y) || (position.Y < MIN_Y) {
                return false
        }

        if (position.X > MAX_X) || (position.X < MIN_X) {
                return false
        }

        return true
}

func (p Position) Add(x int, y int) Position {
        return Position{
                p.X + x,
                p.Y + y,
        }
}

func (b Board) getUpRightPawn(p Position, oneUp int) *Pawn {

        newPos := p.Add(1, oneUp)
        if ! b.isOnBoard(newPos) {
                return nil
        }

        return b.getPawn(newPos)
}

func (b Board) GetFields() [10][10]*Pawn {

        copy := [10][10]*Pawn{}

        for x, row := range b.fields {

                for y := range row {
                        pawn := b.fields[x][y]
                        if pawn != nil {
                                copy[x][y] = &Pawn{pawn.Color, pawn.MovementDirection, pawn.King}
                        }
                }
        }
        return copy
}

func (b Board) getUpLeftPawn(p Position, oneUp int) *Pawn {

        newPos := p.Add(-1, oneUp)
        if ! b.isOnBoard(newPos) {
                return nil
        }

        return b.getPawn(newPos)
}

func (b Board) GameResult() GAME_RESULT {

        blackLeft := false
        whiteLeft := false

        for x, row := range b.fields {

                for y := range row {
                        pawn := b.fields[x][y]
                        if pawn != nil {
                                switch pawn.Color {
                                case BLACK:
                                        blackLeft = true
                                case WHITE:
                                        whiteLeft = true
                                }
                        }
                }
        }

        if whiteLeft && !blackLeft {
                return WIN_WHITE
        }

        if !whiteLeft && blackLeft {
                return WIN_BLACK
        }

        return UNFINISHED
}

func (b Board) String() string {

        str := ""
        for x, row := range b.fields {

                for y, pawn := range row {

                        if pawn != nil {
                                switch pawn.Color {
                                case BLACK:
                                        str += fmt.Sprintf(" B(%d,%d) ", x, y)
                                case WHITE:
                                        str += fmt.Sprintf(" W(%d,%d) ", x, y)
                                }
                        } else {
                                str += "   _    "
                        }
                }

                str += "\n"
        }

        return str
}