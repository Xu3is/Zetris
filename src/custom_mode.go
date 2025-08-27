package src

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
		if !cm.game.customNameEntered {
			cm.game.state = StateEnterName
			cm.game.enterName.nextState = StateGame
			cm.game.enterName.isCustomMode = true
			return nil
		}
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
	overlay.Fill(color.RGBA{20, 30, 50, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	if cm.game.font != nil {
		// Создаём новый GoTextFace для заголовка с большим размером шрифта
		headerFont := &text.GoTextFace{
			Source: cm.game.font.Source,
			Size:   25, // Увеличиваем размер шрифта для заголовка
		}
		headerText := "Пользовательский режим"
		headerY := ScreenHeight/2 - 100 // Аналогично highscore.go (Y=200)
		drawText(screen, headerText, ScreenWidth/2-100, headerY, color.RGBA{180, 220, 255, 255}, headerFont, false)

		// Отрисовка элементов меню под заголовком
		for i, element := range cm.elements {
			y := headerY + 80 + i*40 // Начинаем с отступа 40 пикселей от заголовка (Y=240, 280, 320)
			var clr color.Color = color.RGBA{180, 220, 255, 255}
			if i == cm.selectedIndex {
				clr = color.RGBA{100, 200, 255, 255}
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
			drawText(screen, text, ScreenWidth/2-100, y, clr, cm.game.font, i == cm.selectedIndex)
		}
	} else {
		// Запасной вариант, если шрифт не загружен
		img := ebiten.NewImage(200, 24)
		img.Fill(color.RGBA{180, 220, 255, 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(ScreenWidth/2-100), float64(ScreenHeight/2-100))
		screen.DrawImage(img, op)
	}
}
