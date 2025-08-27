package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

type Piece struct {
	shape     [][]int
	x, y      int
	image     *ebiten.Image
	shapeType string
}

var shapes = map[string][][]int{
	"i": {
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 0, 0},
	},
	"j": {
		{0, 1, 0},
		{0, 1, 0},
		{1, 1, 0},
	},
	"l": {
		{0, 1, 0},
		{0, 1, 0},
		{0, 1, 1},
	},
	"o": {
		{1, 1},
		{1, 1},
	},
	"s": {
		{0, 1, 1},
		{1, 1, 0},
		{0, 0, 0},
	},
	"t": {
		{0, 1, 0},
		{1, 1, 1},
		{0, 0, 0},
	},
	"z": {
		{1, 1, 0},
		{0, 1, 1},
		{0, 0, 0},
	},
}

var shapeTypes = []string{"i", "j", "l", "o", "s", "t", "z"}

func (g *Game) newPiece() *Piece {
	rand.Seed(time.Now().UnixNano())
	shapeType := shapeTypes[rand.Intn(len(shapeTypes))]
	shape := shapes[shapeType]
	return &Piece{
		shape:     shape,
		x:         gridWidth/2 - len(shape[0])/2,
		y:         0,
		image:     g.images[shapeType],
		shapeType: shapeType,
	}
}

func (g *Game) movePiece(dx, dy int) bool {
	newX, newY := g.currentPiece.x+dx, g.currentPiece.y+dy
	if !g.checkCollision(&Piece{shape: g.currentPiece.shape, x: newX, y: newY}) {
		g.currentPiece.x = newX
		g.currentPiece.y = newY
		return true
	}
	return false
}

func (g *Game) rotatePiece() {
	newShape := make([][]int, len(g.currentPiece.shape[0]))
	for i := range newShape {
		newShape[i] = make([]int, len(g.currentPiece.shape))
	}
	for i := 0; i < len(g.currentPiece.shape); i++ {
		for j := 0; j < len(g.currentPiece.shape[0]); j++ {
			newShape[j][len(g.currentPiece.shape)-1-i] = g.currentPiece.shape[i][j]
		}
	}
	if !g.checkCollision(&Piece{shape: newShape, x: g.currentPiece.x, y: g.currentPiece.y}) {
		g.currentPiece.shape = newShape
	}
}

func (g *Game) rotatePieceCounterClockwise() {
	newShape := make([][]int, len(g.currentPiece.shape[0]))
	for i := range newShape {
		newShape[i] = make([]int, len(g.currentPiece.shape))
	}
	for i := 0; i < len(g.currentPiece.shape); i++ {
		for j := 0; j < len(g.currentPiece.shape[0]); j++ {
			newShape[len(g.currentPiece.shape[0])-1-j][i] = g.currentPiece.shape[i][j]
		}
	}
	if !g.checkCollision(&Piece{shape: newShape, x: g.currentPiece.x, y: g.currentPiece.y}) {
		g.currentPiece.shape = newShape
	}
}

func (g *Game) checkCollision(p *Piece) bool {
	for i, row := range p.shape {
		for j, cell := range row {
			if cell == 0 {
				continue
			}
			x, y := p.x+j, p.y+i
			if x < 0 || x >= gridWidth || y >= gridHeight || (y >= 0 && g.grid[y][x] != "") {
				return true
			}
		}
	}
	return false
}

func (g *Game) fixPiece() {
	for i, row := range g.currentPiece.shape {
		for j, cell := range row {
			if cell != 0 {
				g.grid[g.currentPiece.y+i][g.currentPiece.x+j] = g.currentPiece.shapeType
			}
		}
	}
}
