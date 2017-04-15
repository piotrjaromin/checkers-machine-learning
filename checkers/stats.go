package checkers

import "fmt"

type Stats struct {
        WhiteWins   int
        BlackWins   int
        Draws       int
        GamesPlayed int
}

func (s *Stats) update(gs GAME_RESULT) {

        if gs == UNFINISHED {
                return
        }

        switch gs {
        case WIN_BLACK:
                s.BlackWins++
        case WIN_WHITE:
                s.WhiteWins++
        case DRAW:
                s.Draws++
        }

        s.GamesPlayed++
}

func (s Stats) String() string {
        return fmt.Sprintf("Games palyed: %d\nBlack wins: %d\nWhite wins: %d\nDraws: %d", s.GamesPlayed, s.BlackWins, s.WhiteWins, s.Draws)
}