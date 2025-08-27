package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"os"
)

// GameState определяет текущее состояние игры
type GameState int

const (
	StateMenu GameState = iota
	StateGame
	StatePause
	StateCustomMode
	StateSettings
)

// Menu представляет главное меню игры
type Menu struct {
	game          *Game
	buttons       []string
	selectedIndex int
}

// NewMenu создает новое главное меню
func NewMenu(game *Game) *Menu {
	return &Menu{
		game:          game,
		buttons:       []string{"40 линий", "Пользовательский", "Настройки", "Выход"},
		selectedIndex: 0,
	}
}

// Update обновляет состояние главного меню
func (m *Menu) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		m.selectedIndex--
		if m.selectedIndex < 0 {
			m.selectedIndex = len(m.buttons) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		m.selectedIndex++
		if m.selectedIndex >= len(m.buttons) {
			m.selectedIndex = 0
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch m.buttons[m.selectedIndex] {
		case "40 линий":
			m.game.start40Lines()
		case "Пользовательский":
			m.game.state = StateCustomMode
		case "Настройки":
			m.game.state = StateSettings
		case "Выход":
			os.Exit(0)
		}
	}

	// Esc в главном меню ничего не делает
	return nil
}

// Draw отрисовывает главное меню
func (m *Menu) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	if m.game.font != nil {
		drawText(screen, "Zetris", ScreenWidth/2-50, ScreenHeight/2-150, color.White, m.game.font)
	}

	for i, button := range m.buttons {
		y := ScreenHeight/2 - 50 + i*40
		var clr color.Color = color.White
		if i == m.selectedIndex {
			clr = color.RGBA{255, 255, 0, 255}
		}
		if m.game.font != nil {
			drawText(screen, button, ScreenWidth/2-100, y, clr, m.game.font)
		}
	}
}
