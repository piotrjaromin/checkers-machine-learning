package checkers

import (
        _ "image/jpeg"

        "github.com/hajimehoshi/ebiten"

        "image/color"
        "github.com/fogleman/gg"
        "time"
        "log"
)

type OpenGlDraw struct {
        xFieldWidth   int
        yFieldWidth   int
        xFieldsAmount int
        yFieldsAmount int
        gridWidth     int
        dc            *gg.Context
        board         Board
        delay         int
}

var red = color.RGBA{0xff, 0x00, 0x00, 0xff}
var black = color.RGBA{0xff, 0xff, 0xff, 0xff}

func CreateOpenGlDraw(screenWidth, screenHeight, xFieldsAmount, yFieldsAmount int, delay int) Draw {

        dc := gg.NewContext(screenWidth, screenHeight)

        ogd := OpenGlDraw{
                xFieldWidth: dc.Width() / xFieldsAmount,
                yFieldWidth: dc.Height() / yFieldsAmount,
                xFieldsAmount: xFieldsAmount,
                yFieldsAmount: yFieldsAmount,
                dc: dc,
                gridWidth: 2,
                delay: delay,
        }

        if err := ebiten.Run(ogd.update, screenWidth, screenHeight, 2, "Checkers"); err != nil {
                log.Fatal(err)
        }

        return &ogd
}

func (ogd OpenGlDraw) update(screen *ebiten.Image) error {

        //ogd.drawGrid()
        //
        ////println(ogd.board.String())
        //ogd.board.Iterate(func(x, y int, pawn *Pawn) {
        //
        //        if pawn != nil {
        //                println("itterating pawns")
        //                clr := black
        //                if pawn.Color == WHITE {
        //                        clr = red
        //                }
        //                ogd.drawPawn(x, y, clr, pawn.King)
        //        }
        //})
        //
        //ogd.dc.SetRGB(0, 0, 0)
        //img := ogd.dc.Image()
        //ogd.dc.Fill()
        //screen.Fill(color.White)
        //
        //eimg, _ := ebiten.NewImageFromImage(img, ebiten.FilterNearest)
        //
        //screen.DrawImage(eimg, nil)

        return nil
}

func (ogd OpenGlDraw) drawGrid() {

        for x := 1; x < ogd.xFieldsAmount; x++ {
                x1 := x * ogd.xFieldWidth
                ogd.dc.SetRGB(0, 0, 0)
                ogd.dc.SetLineWidth(float64(ogd.gridWidth))
                ogd.dc.DrawLine(float64(x1), 0, float64(x1), float64(ogd.dc.Height()))
                ogd.dc.Stroke()
        }

        for y := 1; y < ogd.yFieldsAmount; y++ {
                y1 := y * ogd.yFieldWidth
                ogd.dc.SetRGB(0, 0, 0)
                ogd.dc.SetLineWidth(float64(ogd.gridWidth))
                ogd.dc.DrawLine(0, float64(y1), float64(ogd.dc.Width()), float64(y1))
                ogd.dc.Stroke()
        }
}

func (ogd OpenGlDraw) drawPawn(posX, posY int, color color.Color, isKing bool) {

        r := ogd.yFieldWidth / 2

        x := posX * ogd.xFieldWidth + ogd.xFieldWidth / 2
        y := posY * ogd.yFieldWidth + ogd.yFieldWidth / 2

        ogd.dc.DrawCircle(float64(x), float64(y), float64(r))
        ogd.dc.SetColor(color)
        ogd.dc.Fill()

        if isKing {
                innerR := ogd.yFieldWidth / 5
                ogd.dc.DrawCircle(float64(x), float64(y), float64(innerR))
                ogd.dc.SetRGB(0, 0, 0)
                ogd.dc.Fill()
        }
}

func (ogd *OpenGlDraw) Draw(b Board) {
        ogd.board = b

        println("///start////")
        println(b.String())
        println("///end////")

        if ogd.delay > 0 {
                time.Sleep(time.Duration(ogd.delay) * time.Millisecond)
        }
}