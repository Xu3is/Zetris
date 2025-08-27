package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
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
	StateEnterName
	StateHighScore
)

// Menu представляет главное меню игры
type Menu struct {
	game          *Game
	buttons       []string
	selectedIndex int
	logoImage     *ebiten.Image
}

// NewMenu создает новое главное меню
func NewMenu(game *Game) *Menu {
	m := &Menu{
		game:          game,
		buttons:       []string{"40 линий", "Пользовательский", "Рекорды", "Настройки", "Выход"},
		selectedIndex: 0,
	}

	// Загрузка изображения логотипа
	path := "src/assets/images/zetris.png"
	if _, err := os.Stat(path); err == nil {
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Printf("Не удалось загрузить zetris.png: %v", err)
		} else {
			m.logoImage = img
		}
	} else {
		log.Printf("Файл zetris.png не найден, используется запасной вариант")
		m.logoImage = ebiten.NewImage(200, -100)
		m.logoImage.Fill(color.RGBA{180, 220, 255, 255})
	}

	return m
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
		case "Рекорды":
			m.game.state = StateHighScore
		case "Настройки":
			m.game.state = StateSettings
		case "Выход":
			os.Exit(0)
		}
	}

	return nil
}

// Draw отрисовывает главное меню
func (m *Menu) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{20, 30, 50, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	// Отрисовка логотипа
	if m.logoImage != nil {
		op := &ebiten.DrawImageOptions{}
		logoWidth, _ := m.logoImage.Size()
		scale := 250.0 / float64(logoWidth) // Масштаб 250 пикселей
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(ScreenWidth/2-125), float64(ScreenHeight/2-250)) // Поднимаем ещё выше
		screen.DrawImage(m.logoImage, op)
	}

	for i, button := range m.buttons {
		y := ScreenHeight/2 - 50 + i*40
		var clr color.Color = color.RGBA{180, 220, 255, 255}
		if i == m.selectedIndex {
			clr = color.RGBA{100, 200, 255, 255}
		}
		if m.game.font != nil {
			drawText(screen, button, ScreenWidth/2-100, y, clr, m.game.font, i == m.selectedIndex)
		}
	}
}
