package checkers

import "log"
import "fmt"
import (
        "os"
        "io/ioutil"
)

var logger = log.New(os.Stdout, "[Teach]", log.LstdFlags)

type EndTeachingFunc func(b Board, p1 PlayerML, p2 PlayerML) bool

func CreateEndAfterGamesPlayed(maxGamesPlayed int) EndTeachingFunc {

        count := 0
        return func(b Board, p1 PlayerML, p2 PlayerML) bool {
                count++
                return count > maxGamesPlayed
        }
}

func Teach(etf EndTeachingFunc, drawer Draw, learningRate float64) Stats {

        logger := log.New(ioutil.Discard, fmt.Sprintf("[%s]", "[MAIN]"), log.LstdFlags)
        board := NewBoard(BLACK, logger)

        p1 := NewMLPlayer(BLACK, RandomParams(), learningRate)
        p2 := NewMLPlayer(WHITE, RandomParams(), learningRate)

        fmt.Printf("P1 starting with: %+v\n", p1.GetState())
        fmt.Printf("P2 starting with: %+v\n", p2.GetState())

        stats := Stats{}

        for !etf(board, p1, p2) {
                gs := PlayOneGame(board, p1, p2, drawer)
                stats.update(gs)
        }

        fmt.Printf("P1 ending with: %+v\n", p1.GetState())
        fmt.Printf("P2 ending with: %+v\n", p2.GetState())

        return stats
}

func PlayOneGame(b Board, p1 PlayerML, p2 PlayerML, drawer Draw) GAME_RESULT {

        //Teach algorithm
        itr := 0
        for UNFINISHED == b.GameResult() {

                movePlayer(&b, p1, p2)

                if b.GameResult() != UNFINISHED {
                        movePlayer(&b, p2, p1)
                }

                logger.Printf("[%d]State %s", itr, calculateState(b).String())
                drawer.Draw(b)
        }

        //End of Teach algorithm

        return b.GameResult()
}

func movePlayer(b *Board, player PlayerML, opponent PlayerML) {

        stateBeforeMoves := calculateState(*b)
        stateValueBeforeMove := calculateStateValue(player, stateBeforeMoves, b.GameResult())
        player.GetMoves(b.Clone())

        //simulate other player move
        tempBoard := b.Clone()
        opponentMoves := opponent.GetMoves(tempBoard.Clone())
        tempBoard.Move(opponentMoves)

        //simulate first player response
        firstPlayerResponseMoves := player.GetMoves(tempBoard.Clone())
        tempBoard.Move(firstPlayerResponseMoves)

        //evaluate new player state
        stateAfterNextMove := calculateStateValue(player, calculateState(tempBoard), tempBoard.GameResult())
        stateError := stateAfterNextMove - stateValueBeforeMove
        player.EvaluateState(stateError, stateBeforeMoves)
}

type StatePlayer struct {
        AllCount     int
        KingCount    int
        AttacksCount int
}

func (sp StatePlayer) ToArray() []int {
        return []int{sp.AllCount, sp.AttacksCount, sp.KingCount}
}

func (sp StatePlayer) String() string {
        return fmt.Sprintf("All %d, Kings %d, attacks %d", sp.AllCount, sp.KingCount, sp.AttacksCount)
}

type GameState struct {
        Black StatePlayer
        White StatePlayer
}

func (gs GameState) ToArray() []int {
        arr := gs.Black.ToArray()
        return append(arr, gs.White.ToArray()...)
}

func (gs GameState) String() string {
        return fmt.Sprintf("Black %s | White %s", gs.Black.String(), gs.White.String())
}

func calculateState(b Board) GameState {

        gameState := GameState{}

        b.Iterate(func(x, y int, p *Pawn) {

                if p != nil {

                        player := &gameState.White
                        if p.Color == BLACK {
                                player = &gameState.Black
                        }
                        player.AllCount++

                        if p.King {
                                player.KingCount++
                        }

                        player.AttacksCount += attacksCount(b, x, y)
                }
        })

        return gameState
}

func attacksCount(b Board, x, y int) int {

        pawnPos := Position{x, y}
        pawn := b.getPawn(pawnPos)

        moves := b.moveOverOpponent(pawnPos, make(Moves, 0), TO_RIGHT, TO_LEFT, int(pawn.MovementDirection), pawn.Color)
        if pawn.King {
                moves = append(moves, b.moveOverOpponent(pawnPos, make(Moves, 0), TO_RIGHT, TO_LEFT, int(pawn.MovementDirection.Opposite()), pawn.Color)...)
        }

        return len(moves)
}

func calculateStateValue(player PlayerML, state GameState, gameResult GAME_RESULT) float64 {

        if gameResult == DRAW {
                return 0
        }

        if player.Color() == BLACK {
                if gameResult == WIN_BLACK {
                        return 100
                } else if gameResult == WIN_WHITE {
                        return -100
                }
        }

        if player.Color() == WHITE {
                if gameResult == WIN_WHITE {
                        return 100
                } else if gameResult == WIN_BLACK {
                        return -100
                }
        }

        params := player.GetState()
        calcForPlayer := func(startIndex int, p StatePlayer) float64 {

                return float64(p.AllCount) * params[startIndex] + float64(p.AttacksCount) * params[startIndex + 1] + float64(p.KingCount) * params[startIndex + 2]
        }

        return calcForPlayer(1, state.Black) + calcForPlayer(4, state.Black)
}