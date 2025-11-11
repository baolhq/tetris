package core

import (
	"time"

	"github.com/baolhq/tetris/internal/consts"
	"github.com/baolhq/tetris/internal/models"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	blocks      []*models.Block
	activeBlock *models.Block
	elapsed     time.Duration
	prevTime    time.Time
}

func Setup() *Game {
	b := models.NewBlock()

	g := &Game{
		activeBlock: b,
		prevTime:    time.Now(),
	}

	return g
}

func handleInput(g *Game) error {
	Input.Update()

	if Input.WasPressed(InputPause) {
		return ebiten.Termination
	}

	if Input.WasPressed(InputUp) || Input.WasPressed(InputEnter) {
		g.activeBlock.Rotate()
	}

	if Input.WasPressed(InputLeft) {
		g.activeBlock.X -= 1
		g.activeBlock.KeepInBound()
	} else if Input.WasPressed(InputRight) {
		g.activeBlock.X += 1
		g.activeBlock.KeepInBound()
	}

	return nil
}

func (g *Game) Update() error {
	if err := handleInput(g); err != nil {
		return err
	}

	g.elapsed += time.Since(g.prevTime)
	g.prevTime = time.Now()

	if g.elapsed >= time.Second {
		g.elapsed = 0
		g.activeBlock.Y += 1
	}

	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	return consts.ScreenWidth, consts.ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	rows, cols := consts.ScreenWidth/consts.CellSize, consts.ScreenHeight/consts.CellSize

	// Draw background and cell outlines
	screen.Fill(consts.BackgroundColor)
	for y := range cols {
		for x := range rows {
			vector.StrokeRect(
				screen, float32(x*consts.CellSize), float32(y*consts.CellSize),
				consts.CellSize, consts.CellSize,
				1, consts.CellOutlineColor, false,
			)
		}
	}

	// Draw blocks
	g.activeBlock.Draw(screen)
	for _, b := range g.blocks {
		b.Draw(screen)
	}
}
