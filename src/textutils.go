package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

// drawText отрисовывает текст с заданным шрифтом, цветом и тенью
func drawText(screen *ebiten.Image, str string, x, y int, clr color.Color, font *text.GoTextFace, isSelected bool) {
	if font != nil {
		// Отрисовка тени (смещение на 2 пикселя вниз и вправо)
		shadowOp := &text.DrawOptions{}
		shadowOp.GeoM.Translate(float64(x+2), float64(y+2))
		shadowColor := color.RGBA{0, 0, 0, 128} // Полупрозрачная чёрная тень
		var shadowColorScale ebiten.ColorScale
		r, g, b, a := shadowColor.RGBA()
		shadowColorScale.SetR(float32(r) / 0xffff)
		shadowColorScale.SetG(float32(g) / 0xffff)
		shadowColorScale.SetB(float32(b) / 0xffff)
		shadowColorScale.SetA(float32(a) / 0xffff)
		shadowOp.ColorScale = shadowColorScale
		shadowOp.LineSpacing = 24
		text.Draw(screen, str, font, shadowOp)

		// Отрисовка основного текста
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		var colorScale ebiten.ColorScale
		r, g, b, a = clr.RGBA()
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
