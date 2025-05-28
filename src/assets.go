package src

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
	"math/rand"
	"os"
)

// loadAssets загружает изображения и шрифт
func (g *Game) loadAssets() error {
	// Загрузка изображений фигур
	for _, shape := range []string{"i", "j", "l", "o", "s", "t", "z"} {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("src/assets/images/%s.png", shape))
		if err != nil {
			log.Printf("Ошибка загрузки %s.png: %v", shape, err)
			// Запасной вариант: создать изображение с цветным прямоугольником
			img = ebiten.NewImage(cellSize, cellSize)
			img.Fill(color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
		}
		g.images[shape] = img
	}

	// Загрузка изображения клетки поля
	img, _, err := ebitenutil.NewImageFromFile("src/assets/images/boardcell.png")
	if err != nil {
		log.Printf("Ошибка загрузки boardcell.png: %v", err)
		// Запасной вариант: создать серый прямоугольник
		img = ebiten.NewImage(cellSize, cellSize)
		img.Fill(color.RGBA{128, 128, 128, 255})
	}
	g.images["boardcell"] = img

	// Загрузка шрифта Times New Roman
	ttfData, err := os.ReadFile("src/assets/Times New Roman.ttf")
	if err != nil {
		log.Printf("Ошибка загрузки шрифта Times New Roman.ttf: %v. Используется запасной вариант.", err)
		g.font = nil
		return nil
	}
	faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(ttfData))
	if err != nil {
		log.Printf("Ошибка создания шрифта Times New Roman: %v. Используется запасной вариант.", err)
		g.font = nil
		return nil
	}
	g.font = &text.GoTextFace{
		Source: faceSource,
		Size:   24,
	}

	return nil
}
