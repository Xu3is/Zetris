package src

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"time"
)

// CustomMode представляет пользовательский режим
type CustomMode struct {
	game          *Game
	elements      []string
	selectedIndex int
	isLimited     bool
	speedLevel    int
}

// NewCustomMode создает новый пользовательский режим
func NewCustomMode(game *Game) *CustomMode {
	return &CustomMode{
		game:          game,
		elements:      []string{"Ограничение линий", "Скорость", "Начать"},
		selectedIndex: 0,
		isLimited:     false,
		speedLevel:    0,
	}
}

// Update обновляет пользовательский режим
func (cm *CustomMode) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		cm.selectedIndex--
		if cm.selectedIndex < 0 {
			cm.selectedIndex = len(cm.elements) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		cm.selectedIndex++
		if cm.selectedIndex >= len(cm.elements) {
			cm.selectedIndex = 0
		}
	}

	if cm.selectedIndex == 0 { // Ограничение линий
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			cm.isLimited = !cm.isLimited
		}
	}

	if cm.selectedIndex == 1 { // Скорость
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			cm.speedLevel--
			if cm.speedLevel < 0 {
				cm.speedLevel = len(speedLevels) - 1
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			cm.speedLevel++
			if cm.speedLevel >= len(speedLevels) {
				cm.speedLevel = 0
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && cm.selectedIndex == 2 {
		cm.game.grid = [gridHeight][gridWidth]string{}
		cm.game.currentPiece = cm.game.newPiece()
		cm.game.nextPiece = cm.game.newPiece()
		cm.game.score = 0
		cm.game.fallSpeed = speedLevels[cm.speedLevel].fallSpeed
		cm.game.isGameOver = false
		cm.game.isPaused = false
		cm.game.isPieceGrounded = false
		cm.game.keyLastAction = make(map[ebiten.Key]time.Time)
		cm.game.keyPressStart = make(map[ebiten.Key]time.Time)
		cm.game.isLimitedTo40Lines = cm.isLimited
		cm.game.isCustomSpeed = true
		cm.game.clearedLines = 0
		cm.game.state = StateGame
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		cm.game.state = StateMenu
	}

	return nil
}

// Draw отрисовывает пользовательский режим
func (cm *CustomMode) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	if cm.game.font != nil {
		drawText(screen, "Пользовательский режим", ScreenWidth/2-120, ScreenHeight/2-150, color.White, cm.game.font)
	}

	for i, element := range cm.elements {
		y := ScreenHeight/2 - 50 + i*40
		var clr color.Color = color.White
		if i == cm.selectedIndex {
			clr = color.RGBA{255, 255, 0, 255}
		}
		var text string
		switch element {
		case "Ограничение линий":
			if cm.isLimited {
				text = "Ограничение линий: 40"
			} else {
				text = "Ограничение линий: Нет"
			}
		case "Скорость":
			text = fmt.Sprintf("Скорость: Уровень %d", speedLevels[cm.speedLevel].level)
		case "Начать":
			text = "Начать"
		}
		if cm.game.font != nil {
			drawText(screen, text, ScreenWidth/2-100, y, clr, cm.game.font)
		}
	}

	if cm.game.font != nil {
		drawText(screen, "Стрелки: Выбор/Изменение, Enter: Начать, Esc: В меню", ScreenWidth/2-200, ScreenHeight/2+150, color.White, cm.game.font)
	}
}
