package core

import (
	"image/color"
	"slices"
	"time"

	"github.com/baolhq/tetris/internal/consts"
	"github.com/baolhq/tetris/internal/models"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	occupied    [][]color.RGBA
	activeBlock *models.Block
	elapsed     time.Duration
	prevTime    time.Time
	baseTimer   time.Duration
	delayTimer  time.Duration
	accelTimer  time.Duration
	moveDir     int // -1=left, 0=none, 1=right
}

func Setup() *Game {
	b := models.NewBlock()

	grid := make([][]color.RGBA, consts.GameRows)
	for y := range grid {
		grid[y] = make([]color.RGBA, consts.GameCols)
		for x := range grid[y] {
			grid[y][x] = consts.BackgroundColor
		}
	}

	g := &Game{
		activeBlock: b,
		occupied:    grid,
		prevTime:    time.Now(),
		baseTimer:   time.Millisecond * 1000,
		delayTimer:  time.Millisecond * 1000,
		accelTimer:  time.Millisecond * 20,
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

	// --- query input once per frame ---
	leftDown, rightDown := Input.IsDown(InputLeft), Input.IsDown(InputRight)
	leftWasRep, rightWasRep := Input.WasRepeated(InputLeft), Input.WasRepeated(InputRight)
	leftWasPressed, rightWasPressed := Input.WasPressed(InputLeft), Input.WasPressed(InputRight)

	moveLeft := leftWasPressed || leftWasRep
	moveRight := rightWasPressed || rightWasRep

	// --- update direction lock ---
	if g.moveDir == 0 {
		if leftWasRep {
			g.moveDir = -1
		} else if rightWasRep {
			g.moveDir = 1
		}
	}

	// Unlock if locked key is released
	switch g.moveDir {
	case -1:
		if !leftDown {
			g.moveDir = 0
			if rightDown && rightWasRep {
				g.moveDir = 1
			}
		}
	case 1:
		if !rightDown {
			g.moveDir = 0
			if leftDown && leftWasRep {
				g.moveDir = -1
			}
		}
	}

	// Cancel opposite key if a direction is locked
	switch g.moveDir {
	case -1:
		moveRight = false
	case 1:
		moveLeft = false
	}

	// Cancel both if simultaneous press without lock
	if moveLeft && moveRight && g.moveDir == 0 {
		moveLeft, moveRight = false, false
	}

	// --- apply movement ---
	switch {
	case moveLeft:
		g.activeBlock.X--
		g.activeBlock.KeepInBound()
		if checkCollision(g) {
			g.activeBlock.X++
		}
	case moveRight:
		g.activeBlock.X++
		g.activeBlock.KeepInBound()
		if checkCollision(g) {
			g.activeBlock.X--
		}
	}

	// --- handle soft drop ---
	if Input.IsDown(InputDown) {
		g.baseTimer = g.accelTimer
	} else if Input.WasReleased(InputDown) {
		g.baseTimer = g.delayTimer
	}

	return nil
}

func checkCollision(g *Game) bool {
	for _, cell := range g.activeBlock.Shape {
		x := g.activeBlock.X + cell[0]
		y := g.activeBlock.Y + cell[1]

		if y >= consts.GameRows {
			return true
		}

		if x < 0 || x > consts.GameCols || y < 0 {
			continue
		}

		if g.occupied[y][x] != consts.BackgroundColor {
			return true
		}
	}

	return false
}

func checkComplete(g *Game) []bool {
	comp := make([]bool, len(g.occupied))

	for y := range len(g.occupied) {
		lineComp := true

		for x := range len(g.occupied[y]) {
			if g.occupied[y][x] == consts.BackgroundColor {
				lineComp = false
			}
		}

		comp[y] = lineComp
	}

	return comp
}

func shiftDown(g *Game, comp []bool) {
	for y := range g.occupied {
		if comp[y] {
			// Shift all rows above this one down by one
			for py := y; py > 0; py-- {
				copy(g.occupied[py], g.occupied[py-1])
			}
			// Clear the top row
			for x := range g.occupied[0] {
				g.occupied[0][x] = consts.BackgroundColor
			}
		}
	}
}

func (g *Game) Update() error {
	if err := handleInput(g); err != nil {
		return err
	}

	g.elapsed += time.Since(g.prevTime)
	g.prevTime = time.Now()

	if g.elapsed >= g.baseTimer {
		g.elapsed = 0
		g.activeBlock.Y++

		if checkCollision(g) {
			g.activeBlock.Y--

			for _, c := range g.activeBlock.Shape {
				x := g.activeBlock.X + c[0]
				y := g.activeBlock.Y + c[1]
				g.occupied[y][x] = g.activeBlock.Color
			}

			g.activeBlock = models.NewBlock()
		}

		comp := checkComplete(g)
		if slices.Contains(comp, true) {
			shiftDown(g, comp)
			g.baseTimer -= time.Millisecond * 100
			g.delayTimer -= time.Millisecond * 100
		}
	}

	return nil
}

func (g *Game) Layout(w, h int) (int, int) {
	return consts.ScreenWidth, consts.ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background and cell outlines
	screen.Fill(consts.BackgroundColor)
	for y := range consts.GameRows {
		for x := range consts.GameCols {
			vector.StrokeRect(
				screen, float32(x*consts.CellSize), float32(y*consts.CellSize),
				consts.CellSize, consts.CellSize,
				1, consts.CellOutlineColor, false,
			)
		}
	}

	// Draw blocks
	models.DrawBlock(screen, g.activeBlock)
	for y := range g.occupied {
		for x := range g.occupied[y] {
			clr := g.occupied[y][x]
			if clr == consts.BackgroundColor {
				continue
			}

			vector.FillRect(
				screen, float32(x*consts.CellSize), float32(y*consts.CellSize),
				consts.CellSize, consts.CellSize, clr, false,
			)
		}
	}
}
