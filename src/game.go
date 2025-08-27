package src

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
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
	{fallSpeed: 0.1, level: 1}, // Самая медленная
	{fallSpeed: 0.0667, level: 2},
	{fallSpeed: 0.05, level: 3},
	{fallSpeed: 0.04, level: 4},
	{fallSpeed: 0.0333, level: 5}, // Самая быстрая
}

type Game struct {
	settingsMenu       *SettingsMenu
	grid               [gridHeight][gridWidth]string
	currentPiece       *Piece
	nextPiece          *Piece
	score              int
	fallTime           float64
	fallSpeed          float64
	isPaused           bool
	isGameOver         bool
	images             map[string]*ebiten.Image
	font               *text.GoTextFace
	lastUpdate         time.Time
	keyLastAction      map[ebiten.Key]time.Time
	keyPressStart      map[ebiten.Key]time.Time
	lockDelayStart     time.Time
	isPieceGrounded    bool
	state              GameState
	lastState          GameState
	menu               *Menu
	pauseMenu          *PauseMenu
	customMode         *CustomMode
	isLimitedTo40Lines bool
	isCustomSpeed      bool
	clearedLines       int
	menuPlayer         *audio.Player
	customPlayer       *audio.Player
	gamePlayer         *audio.Player
}

func NewGame() (*Game, error) {
	g := &Game{
		fallSpeed:     speedLevels[0].fallSpeed,
		images:        make(map[string]*ebiten.Image),
		lastUpdate:    time.Now(),
		keyLastAction: make(map[ebiten.Key]time.Time),
		keyPressStart: make(map[ebiten.Key]time.Time),
		state:         StateMenu,
		lastState:     StateMenu,
	}
	g.settingsMenu = NewSettingsMenu(g)
	err := g.loadAssets()
	if err != nil {
		return nil, err
	}

	g.menu = NewMenu(g)
	g.pauseMenu = NewPauseMenu(g)
	g.customMode = NewCustomMode(g)
	g.currentPiece = g.newPiece()
	g.nextPiece = g.newPiece()
	return g, nil
}

func (g *Game) Update() error {
	// Управление музыкой при смене состояния
	if g.state != g.lastState {
		// Остановка всех плееров
		if g.menuPlayer != nil && g.menuPlayer.IsPlaying() {
			g.menuPlayer.Pause()
		}
		if g.gamePlayer != nil && g.gamePlayer.IsPlaying() {
			g.gamePlayer.Pause()
		}
		if g.customPlayer != nil && g.customPlayer.IsPlaying() {
			g.customPlayer.Pause()
		}

		// Запуск соответствующего плеера
		switch g.state {
		case StateMenu:
			if g.menuPlayer != nil {
				g.menuPlayer.Rewind()
				g.menuPlayer.Play()
			}
		case StateGame:
			if g.isCustomSpeed {
				if g.customPlayer != nil {
					g.customPlayer.Rewind()
					g.customPlayer.Play()
				}
			} else {
				if g.gamePlayer != nil {
					g.gamePlayer.Rewind()
					g.gamePlayer.Play()
				}
			}
		case StateCustomMode, StatePause, StateSettings:
			// Музыка не играет
		}
	}

	// Обновление меню настроек
	if g.state == StateSettings {
		err := g.settingsMenu.Update()
		if err != nil {
			return err
		}
		return nil
	}

	// Обновление главного меню
	if g.state == StateMenu {
		err := g.menu.Update()
		if err != nil {
			return err
		}
		return nil
	}

	// Обновление паузы
	if g.state == StatePause {
		err := g.pauseMenu.Update()
		if err != nil {
			return err
		}
		return nil
	}

	// Обновление пользовательского режима
	if g.state == StateCustomMode {
		err := g.customMode.Update()
		if err != nil {
			return err
		}
		return nil
	}

	// Обновление игры
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
			g.clearedLines = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.state = StateMenu
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
			g.clearedLines = 0
			g.isCustomSpeed = false
			g.isLimitedTo40Lines = false
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if !g.isPaused {
			g.isPaused = true
			g.state = StatePause
		} else {
			g.state = StateMenu
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
			g.clearedLines = 0
			g.isCustomSpeed = false
			g.isLimitedTo40Lines = false
		}
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

	if !g.isCustomSpeed {
		scoreThreshold := g.score / 5000
		if scoreThreshold >= len(speedLevels) {
			scoreThreshold = len(speedLevels) - 1
		}
		g.fallSpeed = speedLevels[scoreThreshold].fallSpeed
	}

	g.lastState = g.state
	return nil
}

func (g *Game) clearLines() {
	linesCleared := 0
	for i := gridHeight - 1; i >= 0; i-- {
		filled := true
		for j := 0; j < gridWidth; j++ {
			if g.grid[i][j] == "" {
				filled = false
				break
			}
		}
		if filled {
			linesCleared++
			for j := i; j > 0; j-- {
				g.grid[j] = g.grid[j-1]
			}
			g.grid[0] = [gridWidth]string{}
			i++
		}
	}
	g.score += linesCleared * 100
	g.clearedLines += linesCleared

	if g.isLimitedTo40Lines && g.clearedLines >= 40 {
		g.isGameOver = true
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == StateSettings {
		g.settingsMenu.Draw(screen)
		return
	}

	if g.state == StateMenu {
		g.menu.Draw(screen)
		return
	}
	if g.state == StatePause {
		g.pauseMenu.Draw(screen)
		return
	}
	if g.state == StateCustomMode {
		g.customMode.Draw(screen)
		return
	}

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

	if g.isGameOver {
		overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
		overlay.Fill(color.RGBA{0, 0, 0, 192})
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(overlay, op)

		if g.isLimitedTo40Lines && g.clearedLines >= 40 {
			if g.font != nil {
				drawText(screen, "Вы выиграли!", ScreenWidth/2-80, ScreenHeight/2-100, color.White, g.font)
			}
		} else {
			if g.font != nil {
				drawText(screen, "Вы проиграли!", ScreenWidth/2-80, ScreenHeight/2-100, color.White, g.font)
			}
		}
		if g.font != nil {
			drawText(screen, fmt.Sprintf("Итоговый счёт: %d", g.score), ScreenWidth/2-80, ScreenHeight/2-60, color.White, g.font)
			drawText(screen, fmt.Sprintf("Очищено линий: %d", g.clearedLines), ScreenWidth/2-80, ScreenHeight/2-20, color.White, g.font)
			drawText(screen, "R: Перезапустить", ScreenWidth/2-80, ScreenHeight/2+20, color.White, g.font)
			drawText(screen, "Q: В меню", ScreenWidth/2-80, ScreenHeight/2+60, color.White, g.font)
		}
	} else if !g.isPaused {
		if g.font != nil {
			drawText(screen, fmt.Sprintf("Счёт: %d", g.score), ScreenWidth-150, 30, color.White, g.font)
			drawText(screen, fmt.Sprintf("Очищено линий: %d", g.clearedLines), ScreenWidth-150, 60, color.White, g.font)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) start40Lines() {
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
	g.isLimitedTo40Lines = true
	g.isCustomSpeed = false
	g.clearedLines = 0
	g.state = StateGame
}
