package models

import (
	"image/color"
	"math/rand"

	"github.com/baolhq/tetris/internal/consts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BlockShape [][]int

var (
	BlockO BlockShape = [][]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	BlockI BlockShape = [][]int{{0, 0}, {0, 1}, {0, 2}, {0, 3}}
	BlockT BlockShape = [][]int{{0, 0}, {1, 0}, {2, 0}, {1, 1}}
	BlockS BlockShape = [][]int{{1, 0}, {2, 0}, {1, 1}, {0, 1}}
	BlockZ BlockShape = [][]int{{0, 0}, {1, 0}, {1, 1}, {2, 1}}
	BlockL BlockShape = [][]int{{0, 0}, {0, 1}, {0, 2}, {1, 2}}
	BlockJ BlockShape = [][]int{{1, 0}, {1, 1}, {1, 2}, {0, 2}}
)

type Block struct {
	X, Y          int
	Width, Height int
	Shape         BlockShape
	Color         color.RGBA
}

func pickRandomBlock() (BlockShape, color.RGBA) {
	blocks := []struct {
		shape BlockShape
		color color.RGBA
	}{
		{BlockO, consts.Red},
		{BlockI, consts.Green},
		{BlockT, consts.Blue},
		{BlockS, consts.Yellow},
		{BlockZ, consts.Purple},
		{BlockL, consts.Orange},
		{BlockJ, consts.Aqua},
	}

	idx := rand.Intn(len(blocks))
	return blocks[idx].shape, blocks[idx].color
}

func NewBlock() *Block {
	shape, color := pickRandomBlock()
	b := &Block{
		Shape: shape,
		Color: color,
	}

	b.UpdateDimensions()
	maxX := consts.ScreenWidth/consts.CellSize - b.Width
	b.X = rand.Intn(maxX + 1)
	b.Y = -b.Height

	// Random rotations
	for range 3 {
		if rand.Intn(2) == 0 {
			break
		}
		b.Rotate()
	}

	return b
}

func (b *Block) UpdateDimensions() {
	maxX, maxY := 0, 0
	for _, p := range b.Shape {
		if p[0] > maxX {
			maxX = p[0]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}
	b.Width = maxX + 1
	b.Height = maxY + 1
}

func (b *Block) Rotate() {
	newShape := make(BlockShape, len(b.Shape))
	minX, minY := 1<<31-1, 1<<31-1

	for i, p := range b.Shape {
		x, y := p[0], p[1]
		rx, ry := y, -x
		newShape[i] = []int{rx, ry}

		if rx < minX {
			minX = rx
		}
		if ry < minY {
			minY = ry
		}
	}

	// Shift shape back to positive space
	for i, p := range newShape {
		newShape[i][0] = p[0] - minX
		newShape[i][1] = p[1] - minY
	}

	b.Shape = newShape
	b.UpdateDimensions()
	b.KeepInBound()
}

func (b *Block) KeepInBound() {
	w, h := b.Width, b.Height
	gridW := consts.ScreenWidth / consts.CellSize
	gridH := consts.ScreenHeight / consts.CellSize

	if b.X+w > gridW {
		b.X = gridW - w
	} else if b.X < 0 {
		b.X = 0
	}

	if b.Y+h > gridH {
		b.Y = gridH - h
	}
}

func DrawBlock(screen *ebiten.Image, block *Block) {
	for _, p := range block.Shape {
		x := (block.X + p[0]) * consts.CellSize
		y := (block.Y + p[1]) * consts.CellSize

		vector.FillRect(
			screen,
			float32(x), float32(y),
			consts.CellSize,
			consts.CellSize,
			block.Color,
			false,
		)
	}
}
