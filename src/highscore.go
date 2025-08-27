package src

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

// HighScoreScreen представляет экран рекордов
type HighScoreScreen struct {
	game *Game
}

// NewHighScoreScreen создает новый экран рекордов
func NewHighScoreScreen(game *Game) *HighScoreScreen {
	return &HighScoreScreen{
		game: game,
	}
}

// Update обновляет экран рекордов
func (hs *HighScoreScreen) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		hs.game.state = StateMenu
	}
	return nil
}

// Draw отрисовывает экран рекордов
func (hs *HighScoreScreen) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{20, 30, 50, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	if hs.game.font != nil {
		// Центрирование заголовка "Рекорды"
		headerText := "Рекорды"
		w, _ := text.Measure(headerText, hs.game.font, 32)
		drawText(screen, headerText, ScreenWidth/2-int(w/2), ScreenHeight/2-150, color.RGBA{180, 220, 255, 255}, hs.game.font, false)

		classicName := hs.game.classicHighScore.Name
		if classicName == "" {
			classicName = "–"
		}
		customName := hs.game.customHighScore.Name
		if customName == "" {
			customName = "–"
		}

		// Центрирование текста
		line1 := "40 линий:"
		w, _ = text.Measure(line1, hs.game.font, 24)
		drawText(screen, line1, ScreenWidth/2-int(w/2), ScreenHeight/2-50, color.RGBA{180, 220, 255, 255}, hs.game.font, false)

		line2 := fmt.Sprintf("%s: %d", classicName, hs.game.classicHighScore.Score)
		w, _ = text.Measure(line2, hs.game.font, 24)
		drawText(screen, line2, ScreenWidth/2-int(w/2), ScreenHeight/2-20, color.RGBA{180, 220, 255, 255}, hs.game.font, false)

		line3 := "Пользовательский:"
		w, _ = text.Measure(line3, hs.game.font, 24)
		drawText(screen, line3, ScreenWidth/2-int(w/2), ScreenHeight/2+20, color.RGBA{180, 220, 255, 255}, hs.game.font, false)

		line4 := fmt.Sprintf("%s: %d", customName, hs.game.customHighScore.Score)
		w, _ = text.Measure(line4, hs.game.font, 24)
		drawText(screen, line4, ScreenWidth/2-int(w/2), ScreenHeight/2+50, color.RGBA{180, 220, 255, 255}, hs.game.font, false)

	}
}
