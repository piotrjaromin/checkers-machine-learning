package checkers

import (
        "fmt"
)

type Draw interface {
        Draw(b Board)
}

type AsciDraw int

func (a AsciDraw) Draw(b Board) {
        fmt.Print("\n===============\n")
        fmt.Printf(b.String())
        fmt.Print("\n===============\n")
}

type DoNothingDraw int

func (dnd DoNothingDraw) Draw(b Board) {

}