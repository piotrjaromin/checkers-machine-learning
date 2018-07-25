package checkers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Draw interface {
	Draw(b Board)
}

type AsciDraw struct{}

func (a AsciDraw) Draw(b Board) {
	fmt.Print("\n===============\n")
	fmt.Printf(b.String())
	fmt.Print("\n===============\n")
}

type DoNothingDraw struct{}

func (dnd DoNothingDraw) Draw(b Board) {

}

type HttpDraw struct {
	board     Board
	boardChan chan Board
	delay     time.Duration
}

func NewHttpDraw(delay int) *HttpDraw {
	h := &HttpDraw{
		delay:     time.Duration(delay),
		boardChan: make(chan Board),
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)

		// b := h.board

		b := <-h.boardChan

		json.NewEncoder(w).Encode(struct {
			Board [10][10]*Pawn
		}{
			b.fields,
		})
	}

	http.HandleFunc("/", handler)
	go func() { log.Fatal(http.ListenAndServe(":8080", nil)) }()
	return h
}

func (h *HttpDraw) Draw(b Board) {
	h.board = b
	h.boardChan <- b
	// time.Sleep(h.delay * time.Millisecond)
}
