package main

import (
	"fmt"
	"image/color"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"

	"snake/model"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Snake",
		Bounds: pixel.R(0, 0, 925, 525),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	game := model.Start()
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))

	score := text.New(pixel.V(10, game.Height-10), atlas)

	last := time.Now()
	for !win.Closed() && !game.Over {
		win.Clear(color.RGBA{0, 43, 54, 0})
		//dt := time.Since(last).Seconds()

		fmt.Fprintln(score, "Score: "+fmt.Sprint(game.Score))
		game.Foods.Shape.Draw(win)
		game.Player.Head.Shape.Draw(win)
		for _, n := range game.Player.Body {
			n.Shape.Draw(win)
		}
		win.Update()

		score.Clear()
		score.Write([]byte("Score: " + fmt.Sprint(game.Score)))
		score.Draw(win, pixel.IM.Scaled(score.Orig, 2))

		if win.JustPressed(pixelgl.KeyUp) || win.JustPressed(pixelgl.KeyW) {
			game.Player.Turn(model.Vector2{X: 0, Y: 1}, game)
		}
		if win.JustPressed(pixelgl.KeyDown) || win.JustPressed(pixelgl.KeyS) {
			game.Player.Turn(model.Vector2{X: 0, Y: -1}, game)
		}
		if win.JustPressed(pixelgl.KeyLeft) || win.JustPressed(pixelgl.KeyA) {
			game.Player.Turn(model.Vector2{X: -1, Y: 0}, game)
		}
		if win.JustPressed(pixelgl.KeyRight) || win.JustPressed(pixelgl.KeyD) {
			game.Player.Turn(model.Vector2{X: 1, Y: 0}, game)
		}
		if win.JustPressed(pixelgl.KeySpace) {
			game.Player.AddNode()
		}

		if time.Since(last).Seconds() >= game.Speed {
			game.Player.Move(game)
			last = time.Now()
		}
	}
}
