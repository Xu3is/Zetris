package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"strings"
)

// EnterNameScreen представляет экран ввода имени
type EnterNameScreen struct {
	game         *Game
	input        string
	nextState    GameState
	isCustomMode bool
}

// NewEnterNameScreen создает новый экран ввода имени
func NewEnterNameScreen(game *Game) *EnterNameScreen {
	return &EnterNameScreen{
		game:         game,
		input:        "",
		nextState:    StateGame,
		isCustomMode: false,
	}
}

// Update обновляет экран ввода имени
func (ens *EnterNameScreen) Update() error {
	var keys []ebiten.Key
	inpututil.AppendJustPressedKeys(keys[:0])

	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		if key >= ebiten.KeyA && key <= ebiten.KeyZ {
			if len(ens.input) < 10 { // Ограничение длины имени
				ens.input += strings.ToUpper(string('A' + key - ebiten.KeyA))
			}
		}
		if key == ebiten.KeyBackspace && len(ens.input) > 0 {
			ens.input = ens.input[:len(ens.input)-1]
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(ens.input) > 0 {
		if ens.isCustomMode {
			ens.game.customPlayerName = ens.input
			ens.game.customNameEntered = true
			ens.game.state = StateCustomMode
		} else {
			ens.game.classicPlayerName = ens.input
			ens.game.classicNameEntered = true
			ens.game.start40Lines()
		}
		ens.input = ""
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ens.game.state = StateMenu
		ens.input = ""
	}

	return nil
}

// Draw отрисовывает экран ввода имени
func (ens *EnterNameScreen) Draw(screen *ebiten.Image) {
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{20, 30, 50, 192})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	if ens.game.font != nil {
		drawText(screen, "Введите имя:", ScreenWidth/2-80, ScreenHeight/2-100, color.RGBA{180, 220, 255, 255}, ens.game.font, false)
		drawText(screen, ens.input+"_", ScreenWidth/2-80, ScreenHeight/2-20, color.RGBA{180, 220, 255, 255}, ens.game.font, false)
		drawText(screen, "Enter: Подтвердить, Esc: В меню", ScreenWidth/2-180, ScreenHeight/2+100, color.RGBA{180, 220, 255, 255}, ens.game.font, false)
	}
}
