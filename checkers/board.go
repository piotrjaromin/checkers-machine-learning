package checkers

import (
        "fmt"
        "io/ioutil"
        "log"
        "math"
        "os"
        "reflect"
)

const DRAW_MOVES = 20

type GAME_RESULT string

var EnableCloneLogging = false

const (
        DRAW GAME_RESULT = "Draw"
        WIN_BLACK GAME_RESULT = "WinBlack"
        WIN_WHITE GAME_RESULT = "WinWhite"
        UNFINISHED GAME_RESULT = "Unfinised"
)

func ( gr GAME_RESULT) isWin(color Color) bool {
        if gr == WIN_BLACK && color == BLACK {
                return true
        }

        if gr == WIN_WHITE && color == WHITE {
                return true
        }

        return false
}

func ( gr GAME_RESULT) isLost(color Color) bool {
        if gr == WIN_WHITE && color == BLACK {
                return true
        }

        if gr == WIN_BLACK && color == WHITE {
                return true
        }

        return false
}


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
        fields              [10][10]*Pawn
        topPlayer           Color //player on top of board ( starting position closer to MAX_X, MAX_Y
        roundsAfterLastKill int
        boardLogger         *log.Logger
        cloneCount          int
}

func NewBoard(topColor Color, logger *log.Logger) Board {

        b := Board{topPlayer: topColor, boardLogger: logger}

        for y := 0; y < 4; y++ {

                for x := 0; x <= MAX_X; x++ {
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

func (b *Board) Move(moves Moves) bool {
        if b.isInvalidMoves(moves) {
                b.boardLogger.Println("invalid move aborting.")
                return false
        }

        if len(moves) == 0 {
                b.boardLogger.Println("No moves to make.")
                return false
        }

        b.boardLogger.Printf("moves amount %d\n", len(moves))

        for _, move := range moves {

                pawn := b.fields[move.From.X][move.From.Y]
                b.boardLogger.Printf("Move %s > %s", pawn.Color.String(), move.String())

                b.fields[move.To.X][move.To.Y] = pawn
                b.fields[move.From.X][move.From.Y] = nil

                isKillMove := func(move Move) bool {
                        return math.Abs(float64(move.From.Y - move.To.Y)) > 1
                }

                if isKillMove(move) {
                        //calculate diff should be either -1 or +1
                        diffY := int(math.Ceil(float64((move.To.Y - move.From.Y) / 2)))
                        diffX := int(math.Ceil(float64((move.To.X - move.From.X) / 2)))

                        destX := move.From.X + diffX
                        destY := move.From.Y + diffY

                        //removes pawn
                        pawnToKill := b.fields[destX][destY]
                        b.boardLogger.Printf("killing > %s at (%d, %d)\n", pawnToKill.Color.String(), destX, destY)
                        b.fields[destX][destY] = nil

                        //reset count after kill
                        b.roundsAfterLastKill = 0
                }
        }

        b.roundsAfterLastKill++

        lastMove := moves[len(moves) - 1]
        pawn := b.getPawn(Position{lastMove.To.X, lastMove.To.Y})
        if b.canBeKingReaching(lastMove.To) {
                b.boardLogger.Printf("Turning pawn %+v into king because it reached %+v", pawn, lastMove)
                pawn.King = true
        }

        return true
}

func (b Board) canBeKingReaching(position Position) bool {

        return position.Y == MAX_Y || position.Y == MIN_Y
}

func (b Board) isInvalidMoves(moves Moves) bool {

        //TODO can the be zero?
        if len(moves) == 0 {
                return false
        }

        validMoves := b.GetValidMovesForPosition(moves[0].From)

        for _, validMove := range validMoves {
                if reflect.DeepEqual(validMove, moves) {
                        return false
                }
        }

        return true
}

func (b Board) getPawn(p Position) *Pawn {

        if b.isOnBoard(p) {
                return b.fields[p.X][p.Y]
        }
        return nil
}

//GetValidMovesForPosition returns list of valid moves for pawn which is posititioned at provided Position
//if there is now pawn at given position return empty array
func (b Board) GetValidMovesForPosition(position Position) []Moves {
        moves := []Moves{}

        if !b.isOnBoard(position) {
                return moves
        }

        pawn := b.getPawn(position)

        if pawn == nil {
                b.boardLogger.Println("[GetValidMovesForPosition] There is no pawn at this position")
                return moves
        }

        //top bottom multipliers (up/down is more of player perspective)
        verticalMoves := func(verticalDirection Direction) {
                dir := int(verticalDirection)
                moves = appendMoves(moves, b.ifPossibleUpToSide(position, TO_RIGHT, dir))
                moves = appendMoves(moves, b.ifPossibleUpToSide(position, TO_LEFT, dir))
                moves = append(moves, b.moveOverOpponent(position, make(Moves, 0), TO_RIGHT, TO_LEFT, dir, pawn.Color)...)
        }

        verticalMoves(pawn.MovementDirection)
        if pawn.King {
                verticalMoves(pawn.MovementDirection.Opposite())
        }

        return moves
}

func (b Board) ifPossibleUpToSide(p Position, toSide int, oneUp int) Moves {
        moves := Moves{}

        if !b.canMove(p, toSide, oneUp) {
                return moves
        }

        getUpToSide := func(p Position) *Pawn {
                return b.fields[p.X + toSide][p.Y + oneUp]
        }

        if getUpToSide(p) == nil {
                moves = append(moves, Move{From: p, To: p.Add(toSide, oneUp)})
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
                jumpOver := b.ifPossibleUpToSide(newRightUp, toRight, oneUp)
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

                jumpOver := b.ifPossibleUpToSide(newLeftUp, toLeft, oneUp)

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
        if !b.isOnBoard(newPos) {
                return nil
        }

        return b.getPawn(newPos)
}

func (b Board) GetFields() [10][10]*Pawn {

        copy := [10][10]*Pawn{}

        b.Iterate(func(x, y int, pawn *Pawn) {
                if pawn != nil {
                        copy[x][y] = &Pawn{pawn.Color, pawn.MovementDirection, pawn.King}
                }
        })

        return copy
}

func (b Board) getUpLeftPawn(p Position, oneUp int) *Pawn {

        newPos := p.Add(-1, oneUp)
        if !b.isOnBoard(newPos) {
                return nil
        }

        return b.getPawn(newPos)
}

func (b Board) GameResult() GAME_RESULT {

        blackLeft := false
        whiteLeft := false

        b.Iterate(func(x, y int, pawn *Pawn) {

                if pawn != nil {
                        switch pawn.Color {
                        case BLACK:
                                blackLeft = true
                        case WHITE:
                                whiteLeft = true
                        }
                }
        })

        if whiteLeft && !blackLeft {
                return WIN_WHITE
        }

        if !whiteLeft && blackLeft {
                return WIN_BLACK
        }

        if b.roundsAfterLastKill > DRAW_MOVES {
                return DRAW
        }

        whiteAllMoves := len(b.GetAllMovesFor(WHITE))
        blackAllMoves := len(b.GetAllMovesFor(BLACK))

        if whiteAllMoves == 0 && blackAllMoves == 0 {
                return DRAW
        }

        if whiteAllMoves == 0 && blackAllMoves > 0 {
                return WIN_BLACK
        }

        if blackAllMoves == 0 && whiteAllMoves > 0 {
                return WIN_WHITE
        }

        return UNFINISHED
}

func (b Board) Iterate(cb func(x, y int, p *Pawn)) {

        for x, row := range b.fields {

                for y, pawn := range row {

                        if pawn != nil {

                                cb(x, y, &Pawn{pawn.Color, pawn.MovementDirection, pawn.King})
                        } else {

                                cb(x, y, nil)
                        }
                }
        }
}

func (b Board) CloneWithLogger(logger *log.Logger) Board {
        cloned := b.Clone()
        cloned.boardLogger = logger
        return cloned
}

//Clone current board and creates new one with the same state
func (b *Board) Clone() Board {

        ioOut := ioutil.Discard
        if EnableCloneLogging {
                ioOut = os.Stdout
        }

        cloned := Board{
                topPlayer:   b.topPlayer,
                boardLogger: log.New(ioOut, fmt.Sprintf("[CLONED-%d]", b.cloneCount), log.LstdFlags),
        }

        b.Iterate(func(x, y int, p *Pawn) {
                if p != nil {
                        cloned.fields[x][y] = &Pawn{p.Color, p.MovementDirection, p.King}
                }
        })

        b.cloneCount++
        return cloned
}

//GetAllMovesFor return all possible moves for player with given color
func (b Board) GetAllMovesFor(color Color) []Moves {

        moves := make([]Moves, 0, 0)

        b.Iterate(func(x, y int, p *Pawn) {

                if p != nil && p.Color == color {
                        moves = append(moves, b.GetValidMovesForPosition(Position{x, y})...)
                }
        })

        return moves
}

func (b Board) String() string {

        str := ""
        currentX := 0
        b.Iterate(func(x int, y int, pawn *Pawn) {

                if currentX != x {
                        str += "\n"
                        currentX = x
                }

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
        })

        return str
}
