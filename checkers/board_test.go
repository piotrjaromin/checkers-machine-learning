package checkers

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBoard(t *testing.T) {

	logger := log.New(os.Stdout, fmt.Sprintf("[%s]", "[MAIN]"), log.LstdFlags)

	Convey("Board should", t, func() {

		Convey("for new board with black at top place white pawn at x=5, y=3", func() {
			b := NewBoard(BLACK, logger)

			p := Position{5, 3}
			pawn := b.getPawn(p)

			So(pawn, ShouldNotBeNil)
			So(pawn.Color, ShouldEqual, WHITE)
		})

		Convey("Allow two moves when pawn has no obstacles", func() {

			b := NewBoard(BLACK, logger)
			from := Position{5, 3}

			moves := b.GetValidMovesForPosition(from)

			So(moves, ShouldHaveLength, 2)
			So(moves, should.Resemble, []Moves{
				{{from, Position{6, 4}}},
				{{from, Position{4, 4}}},
			})
		})

		Convey("Allow one move when pawn is on right side of board", func() {

			b := NewBoard(BLACK, logger)
			from := Position{9, 3}

			moves := b.GetValidMovesForPosition(from)

			So(moves, ShouldHaveLength, 1)
			So(moves, should.Resemble, []Moves{
				{{from, Position{8, 4}}},
			})
		})

		Convey("Allow one move when pawn is on left side of board", func() {

			b := NewBoard(WHITE, logger)
			from := Position{9, 3}

			moves := b.GetValidMovesForPosition(from)

			So(moves, ShouldHaveLength, 1)
			So(moves, should.Resemble, []Moves{
				{{from, Position{8, 4}}},
			})
		})

		Convey("Allow jump over one enemy and remove it", func() {

			b := NewBoard(BLACK, logger)

			//place back for white to be able to jump
			So(b.Move(Moves{{Position{1, 3}, Position{2, 4}}}), should.BeTrue)
			So(b.Move(Moves{{Position{2, 4}, Position{3, 5}}}), should.BeTrue)

			from := Position{4, 6}
			moves := b.GetValidMovesForPosition(from)

			So(moves, ShouldHaveLength, 2)
			So(moves, should.Resemble, []Moves{
				{{from, Position{5, 5}}}, //free place (normal move)
				{{from, Position{2, 4}}}, //jump to position
			})

			So(b.getPawn(Position{3, 6}), ShouldBeNil)
		})

		Convey("Allow chain jump over two enemies  || TODO FIX ME UP I d not have chain jump", func() {

			b := NewBoard(BLACK, logger)

			//place back for white to be able to jump
			So(b.Move(Moves{{Position{1, 3}, Position{2, 4}}}), should.BeTrue)
			So(b.Move(Moves{{Position{2, 4}, Position{3, 5}}}), should.BeTrue)

			from := Position{4, 6}
			moves := b.GetValidMovesForPosition(from)

			So(moves, ShouldHaveLength, 2)
			So(moves, should.Resemble, []Moves{
				{{from, Position{5, 5}}}, //free place (normal move)
				{{from, Position{2, 4}}}, //jump to position
			})

			So(b.getPawn(Position{3, 6}), ShouldBeNil)
		})

		Convey("Not allow invalid move", func() {

			b := NewBoard(BLACK, logger)
			So(b.Move(Moves{{Position{1, 3}, Position{7, 4}}}), should.BeFalse)
		})
	})

}
