package checkers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCalculateState(t *testing.T) {

	Convey("Should calculate proper state", t, func() {

		Convey("for new game", func() {

			b := NewBoard(BLACK, logger)

			state := calculateState(b)

			So(state.Black, ShouldResemble, state.White)
			So(state.Black.AllCount, ShouldEqual, 20)
			So(state.Black.AttacksCount, ShouldEqual, 0)
			So(state.Black.KingCount, ShouldEqual, 0)
		})

		Convey("for game with attack", func() {

			b := NewBoard(BLACK, logger)

			//place back for white to be able to jump
			b.Move(Moves{{Position{1, 3}, Position{2, 4}}})
			b.Move(Moves{{Position{2, 4}, Position{3, 5}}})

			state := calculateState(b)

			So(state.Black, ShouldNotResemble, state.White)
			So(state.Black.AllCount, ShouldEqual, 20)
			So(state.Black.AttacksCount, ShouldEqual, 2)
			So(state.Black.KingCount, ShouldEqual, 0)
		})
	})

}
