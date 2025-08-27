package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"time"
)

// PauseMenu представляет меню паузы
type PauseMenu struct {
	game          *Game
	buttons       []string
	selectedIndex int
}

// NewPauseMenu создает новое меню паузы
func NewPauseMenu(game *Game) *PauseMenu {
	return &PauseMenu{
		game:          game,
		buttons:       []string{"Продолжить", "Перезапустить", "В меню"},
		selectedIndex: 0,
	}
}

// Update обновляет меню паузы
func (pm *PauseMenu) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		pm.selectedIndex--
		if pm.selectedIndex < 0 {
			pm.selectedIndex = len(pm.buttons) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		pm.selectedIndex++
		if pm.selectedIndex >= len(pm.buttons) {
			pm.selectedIndex = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch pm.buttons[pm.selectedIndex] {
		case "Продолжить":
			pm.game.state = StateGame
			pm.game.isPaused = false
		case "Перезапустить":
			pm.game.grid = [gridHeight][gridWidth]string{}
			pm.game.currentPiece = pm.game.newPiece()
			pm.game.nextPiece = pm.game.newPiece()
			pm.game.score = 0
			if !pm.game.isCustomSpeed {
				pm.game.fallSpeed = speedLevels[0].fallSpeed
			}
			pm.game.isGameOver = false
			pm.game.isPaused = false
			pm.game.isPieceGrounded = false
			pm.game.keyLastAction = make(map[ebiten.Key]time.Time)
			pm.game.keyPressStart = make(map[ebiten.Key]time.Time)
			pm.game.clearedLines = 0
			pm.game.state = StateGame
		case "В меню":
			pm.game.state = StateMenu
			pm.game.grid = [gridHeight][gridWidth]string{}
			pm.game.currentPiece = pm.game.newPiece()
			pm.game.nextPiece = pm.game.newPiece()
			pm.game.score = 0
			pm.game.fallSpeed = speedLevels[0].fallSpeed
			pm.game.isGameOver = false
			pm.game.isPaused = false
			pm.game.isPieceGrounded = false
			pm.game.keyLastAction = make(map[ebiten.Key]time.Time)
			pm.game.keyPressStart = make(map[ebiten.Key]time.Time)
			pm.game.clearedLines = 0
			pm.game.isCustomSpeed = false
			pm.game.isLimitedTo40Lines = false
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		pm.game.state = StateGame
		pm.game.isPaused = false
	}

	return nil
}

// Draw отрисовывает меню паузы
func (pm *PauseMenu) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{20, 30, 50, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	if pm.game.font != nil {
		// Центрирование заголовка "Пауза"
		headerText := "Пауза"
		w, _ := text.Measure(headerText, pm.game.font, 32)
		drawText(screen, headerText, ScreenWidth/2-int(w/2), ScreenHeight/2-150, color.RGBA{180, 220, 255, 255}, pm.game.font, false)
	}

	for i, button := range pm.buttons {
		y := ScreenHeight/2 - 50 + i*40
		var clr color.Color = color.RGBA{180, 220, 255, 255}
		if i == pm.selectedIndex {
			clr = color.RGBA{100, 200, 255, 255}
		}
		if pm.game.font != nil {
			// Центрирование текста кнопок
			w, _ := text.Measure(button, pm.game.font, 24)
			drawText(screen, button, ScreenWidth/2-int(w/2), y, clr, pm.game.font, i == pm.selectedIndex)
		}
	}
}
