package main

import (
	"log"

	"github.com/baolhq/tetris/internal/consts"
	"github.com/baolhq/tetris/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(consts.ScreenWidth, consts.ScreenHeight)
	ebiten.SetWindowTitle("Tetris - Ebiten v2")

	if err := ebiten.RunGame(core.Setup()); err != nil {
		log.Fatal(err)
	}
}
