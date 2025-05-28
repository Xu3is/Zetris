package src

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"time"
)

const (
	ScreenWidth       = 600
	ScreenHeight      = 600
	gridWidth         = 10
	gridHeight        = 20
	cellSize          = 24
	keyRepeatDelay    = 150 * time.Millisecond
	keyRepeatInterval = 50 * time.Millisecond
	lockDelayDefault  = 500 * time.Millisecond
	lockDelayLimit    = 5 * time.Second
)

type SpeedLevel struct {
	fallSpeed float64
	level     int
}

var speedLevels = []SpeedLevel{
	{fallSpeed: 0.1, level: 1},
	{fallSpeed: 0.0667, level: 2},
	{fallSpeed: 0.05, level: 3},
	{fallSpeed: 0.04, level: 4},
	{fallSpeed: 0.0333, level: 5},
}

type Game struct {
	grid            [gridHeight][gridWidth]string
	currentPiece    *Piece
	nextPiece       *Piece
	score           int
	fallTime        float64
	fallSpeed       float64
	isPaused        bool
	isGameOver      bool
	images          map[string]*ebiten.Image
	font            *text.GoTextFace
	lastUpdate      time.Time
	keyLastAction   map[ebiten.Key]time.Time
	keyPressStart   map[ebiten.Key]time.Time
	lockDelayStart  time.Time
	isPieceGrounded bool
}

func NewGame() (*Game, error) {
	g := &Game{
		fallSpeed:     speedLevels[0].fallSpeed,
		images:        make(map[string]*ebiten.Image),
		lastUpdate:    time.Now(),
		keyLastAction: make(map[ebiten.Key]time.Time),
		keyPressStart: make(map[ebiten.Key]time.Time),
	}

	err := g.loadAssets()
	if err != nil {
		return nil, err
	}

	g.currentPiece = g.newPiece()
	g.nextPiece = g.newPiece()
	return g, nil
}

func (g *Game) Update() error {
	if g.isGameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.grid = [gridHeight][gridWidth]string{}
			g.currentPiece = g.newPiece()
			g.nextPiece = g.newPiece()
			g.score = 0
			g.fallSpeed = speedLevels[0].fallSpeed
			g.isGameOver = false
			g.isPaused = false
			g.isPieceGrounded = false
			g.keyLastAction = make(map[ebiten.Key]time.Time)
			g.keyPressStart = make(map[ebiten.Key]time.Time)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			return ebiten.Termination
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.isPaused = !g.isPaused
	}

	if g.isPaused {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.movePiece(-1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.movePiece(1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.movePiece(0, 1)
	}

	keys := []ebiten.Key{ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyDown}
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			now := time.Now()
			if _, exists := g.keyPressStart[key]; !exists {
				g.keyPressStart[key] = now
			}

			if now.Sub(g.keyPressStart[key]) >= keyRepeatDelay {
				lastAction, exists := g.keyLastAction[key]
				if !exists || now.Sub(lastAction) >= keyRepeatInterval {
					switch key {
					case ebiten.KeyLeft:
						g.movePiece(-1, 0)
					case ebiten.KeyRight:
						g.movePiece(1, 0)
					case ebiten.KeyDown:
						g.movePiece(0, 1)
					}
					g.keyLastAction[key] = now
				}
			}
		} else {
			delete(g.keyLastAction, key)
			delete(g.keyPressStart, key)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.rotatePieceCounterClockwise()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		g.rotatePiece()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		for g.movePiece(0, 1) {
		}
		g.fixPiece()
		g.clearLines()
		g.currentPiece = g.nextPiece
		g.nextPiece = g.newPiece()
		g.isPieceGrounded = false
		if g.checkCollision(g.currentPiece) {
			g.isGameOver = true
		}
		g.lastUpdate = time.Now()
	}

	if time.Since(g.lastUpdate).Seconds() >= g.fallSpeed {
		if !g.movePiece(0, 1) {
			if !g.isPieceGrounded {
				g.isPieceGrounded = true
				g.lockDelayStart = time.Now()
			}
		} else {
			g.isPieceGrounded = false
		}
		g.lastUpdate = time.Now()
	}

	if g.isPieceGrounded {
		now := time.Now()
		isMoving := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyDown)
		lockDelay := lockDelayDefault
		if isMoving && now.Sub(g.lockDelayStart) < lockDelayLimit {
			lockDelay = lockDelayLimit
		}
		if now.Sub(g.lockDelayStart) >= lockDelay {
			g.fixPiece()
			g.clearLines()
			g.currentPiece = g.nextPiece
			g.nextPiece = g.newPiece()
			g.isPieceGrounded = false
			if g.checkCollision(g.currentPiece) {
				g.isGameOver = true
			}
		}
	}

	scoreThreshold := g.score / 5000
	if scoreThreshold >= len(speedLevels) {
		scoreThreshold = len(speedLevels) - 1
	}
	g.fallSpeed = speedLevels[scoreThreshold].fallSpeed

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	offsetX := (ScreenWidth - gridWidth*cellSize) / 2
	offsetY := (ScreenHeight - gridHeight*cellSize) / 2

	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(j*cellSize+offsetX), float64(i*cellSize+offsetY))
			if g.grid[i][j] != "" {
				screen.DrawImage(g.images[g.grid[i][j]], op)
			} else {
				screen.DrawImage(g.images["boardcell"], op)
			}
		}
	}

	for i, row := range g.currentPiece.shape {
		for j, cell := range row {
			if cell != 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64((g.currentPiece.x+j)*cellSize+offsetX), float64((g.currentPiece.y+i)*cellSize+offsetY))
				screen.DrawImage(g.currentPiece.image, op)
			}
		}
	}

	drawText := func(str string, x, y int, clr color.Color) {
		if g.font != nil {
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			var colorScale ebiten.ColorScale
			switch clr {
			case color.White:
				colorScale.SetR(1.0)
				colorScale.SetG(1.0)
				colorScale.SetB(1.0)
				colorScale.SetA(1.0)
			default:
				r, g, b, a := clr.RGBA()
				colorScale.SetR(float32(r) / 0xffff)
				colorScale.SetG(float32(g) / 0xffff)
				colorScale.SetB(float32(b) / 0xffff)
				colorScale.SetA(float32(a) / 0xffff)
			}
			op.ColorScale = colorScale
			op.LineSpacing = 24
			text.Draw(screen, str, g.font, op)
		} else {
			img := ebiten.NewImage(100, 24)
			img.Fill(clr)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(img, op)
		}
	}

	if g.isGameOver {
		overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
		overlay.Fill(color.RGBA{0, 0, 0, 192})
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(overlay, op)
	}

	if g.isPaused {
		overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
		overlay.Fill(color.RGBA{0, 0, 0, 192})
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(overlay, op)
	}

	if !g.isGameOver && !g.isPaused {
		drawText(fmt.Sprintf("Счёт: %d", g.score), ScreenWidth-150, 30, color.White)
		scoreThreshold := g.score / 5000
		if scoreThreshold >= len(speedLevels) {
			scoreThreshold = len(speedLevels) - 1
		}
		drawText(fmt.Sprintf("Скорость: %d", speedLevels[scoreThreshold].level), ScreenWidth-150, 60, color.White)
		drawText("Для паузы нажмите Esc", 10, 10, color.White)
	}

	if g.isPaused {
		drawText("Пауза", ScreenWidth/2-30, ScreenHeight/2-30, color.White)
	}

	if g.isGameOver {
		drawText("Вы проиграли!", ScreenWidth/2-80, ScreenHeight/2-100, color.White)
		drawText(fmt.Sprintf("Итоговый счёт: %d", g.score), ScreenWidth/2-80, ScreenHeight/2-60, color.White)
		scoreThreshold := g.score / 5000
		if scoreThreshold >= len(speedLevels) {
			scoreThreshold = len(speedLevels) - 1
		}
		drawText(fmt.Sprintf("Итоговая скорость: %d", speedLevels[scoreThreshold].level), ScreenWidth/2-80, ScreenHeight/2-20, color.White)
		drawText("R: Перезапустить", ScreenWidth/2-80, ScreenHeight/2+20, color.White)
		drawText("Q: Выйти", ScreenWidth/2-80, ScreenHeight/2+60, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
