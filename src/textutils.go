package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

// drawText отрисовывает текст с заданным шрифтом и цветом
func drawText(screen *ebiten.Image, str string, x, y int, clr color.Color, font *text.GoTextFace) {
	if font != nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		var colorScale ebiten.ColorScale
		r, g, b, a := clr.RGBA()
		colorScale.SetR(float32(r) / 0xffff)
		colorScale.SetG(float32(g) / 0xffff)
		colorScale.SetB(float32(b) / 0xffff)
		colorScale.SetA(float32(a) / 0xffff)
		op.ColorScale = colorScale
		op.LineSpacing = 24
		text.Draw(screen, str, font, op)
	} else {
		// Запасной вариант: прямоугольник с цветом
		img := ebiten.NewImage(200, 24)
		img.Fill(clr)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(img, op)
	}
}
