package src

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
)

// SettingsMenu представляет экран настроек
type SettingsMenu struct {
	game          *Game
	selectedIndex int
	volume        float64
	resolutions   [][2]int
	resIndex      int
	elements      []string
}

// NewSettingsMenu создает новое меню настроек
func NewSettingsMenu(game *Game) *SettingsMenu {
	return &SettingsMenu{
		game:          game,
		selectedIndex: 0,
		volume:        0.2, // Начальная громкость соответствует menuPlayer
		resolutions: [][2]int{
			{600, 600},
			{800, 600},
			{1024, 768},
			{1280, 720},
			{1280, 1024},
			{1600, 900},
		},
		resIndex: 0,
		elements: []string{"Громкость", "Разрешение", "Применить"},
	}
}

// Update обновляет меню настроек
func (sm *SettingsMenu) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		sm.selectedIndex--
		if sm.selectedIndex < 0 {
			sm.selectedIndex = len(sm.elements) - 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		sm.selectedIndex++
		if sm.selectedIndex >= len(sm.elements) {
			sm.selectedIndex = 0
		}
	}

	if sm.selectedIndex == 0 { // Громкость
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			sm.volume -= 0.05
			if sm.volume < 0 {
				sm.volume = 0
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			sm.volume += 0.05
			if sm.volume > 1 {
				sm.volume = 1
			}
		}
	}

	if sm.selectedIndex == 1 { // Разрешение
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			sm.resIndex--
			if sm.resIndex < 0 {
				sm.resIndex = len(sm.resolutions) - 1
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			sm.resIndex++
			if sm.resIndex >= len(sm.resolutions) {
				sm.resIndex = 0
			}
		}
	}

	// Применение настроек
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && sm.selectedIndex == 2 {
		sm.applySettings()
		sm.game.state = StateMenu
	}

	// Выход в главное меню по Esc
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		sm.game.state = StateMenu
	}

	return nil
}

// Draw отрисовывает меню настроек
func (sm *SettingsMenu) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	if sm.game.font != nil {
		drawText(screen, "Настройки", ScreenWidth/2-70, ScreenHeight/2-150, color.White, sm.game.font)
	}

	for i, element := range sm.elements {
		y := ScreenHeight/2 - 50 + i*40
		var clr color.Color = color.White
		if i == sm.selectedIndex {
			clr = color.RGBA{255, 255, 0, 255}
		}
		var text string
		switch element {
		case "Громкость":
			text = fmt.Sprintf("Громкость: %.0f%%", sm.volume*100)
		case "Разрешение":
			text = fmt.Sprintf("Разрешение: %dx%d", sm.resolutions[sm.resIndex][0], sm.resolutions[sm.resIndex][1])
		case "Применить":
			text = "Применить"
		}
		if sm.game.font != nil {
			drawText(screen, text, ScreenWidth/2-100, y, clr, sm.game.font)
		}
	}
}

// applySettings применяет выбранные настройки
func (sm *SettingsMenu) applySettings() {
	for _, player := range []*audio.Player{sm.game.menuPlayer, sm.game.gamePlayer, sm.game.customPlayer} {
		if player != nil {
			player.SetVolume(float64(sm.volume))
			log.Printf("Установлена громкость %.2f для плеера", sm.volume)
		}
	}
	width, height := sm.resolutions[sm.resIndex][0], sm.resolutions[sm.resIndex][1]
	ebiten.SetWindowSize(width, height)
}
