package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
		buttons:       []string{"В меню"},
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && pm.buttons[pm.selectedIndex] == "В меню" {
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
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

	return nil
}

// Draw отрисовывает меню паузы
func (pm *PauseMenu) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	if pm.game.font != nil {
		drawText(screen, "Пауза", ScreenWidth/2-50, ScreenHeight/2-100, color.White, pm.game.font)
	}

	for i, button := range pm.buttons {
		y := ScreenHeight/2 - 20 + i*40
		var clr color.Color = color.White
		if i == pm.selectedIndex {
			clr = color.RGBA{255, 255, 0, 255}
		}
		if pm.game.font != nil {
			drawText(screen, button, ScreenWidth/2-50, y, clr, pm.game.font)
		}
	}

	if pm.game.font != nil {
		drawText(screen, "Стрелки: Выбор, Enter: Подтвердить, Esc: В меню", ScreenWidth/2-180, ScreenHeight/2+100, color.White, pm.game.font)
	}
}
