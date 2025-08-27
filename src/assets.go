package src

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"math/rand"
	"os"
)

// loadAssets загружает изображения, шрифт и аудиофайлы
func (g *Game) loadAssets() error {
	// Загрузка изображений фигур
	for _, shape := range []string{"i", "j", "l", "o", "s", "t", "z"} {
		path := fmt.Sprintf("src/assets/images/%s.png", shape)
		if _, err := os.Stat(path); err == nil {
			img, _, err := ebitenutil.NewImageFromFile(path)
			if err != nil {
				img = ebiten.NewImage(cellSize, cellSize)
				img.Fill(color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
			}
			g.images[shape] = img
		} else {
			img := ebiten.NewImage(cellSize, cellSize)
			img.Fill(color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
			g.images[shape] = img
		}
	}

	// Загрузка изображения клетки поля
	path := "src/assets/images/boardcell.png"
	if _, err := os.Stat(path); err == nil {
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			img = ebiten.NewImage(cellSize, cellSize)
			img.Fill(color.RGBA{128, 128, 128, 255})
		}
		g.images["boardcell"] = img
	} else {
		img := ebiten.NewImage(cellSize, cellSize)
		img.Fill(color.RGBA{128, 128, 128, 255})
		g.images["boardcell"] = img
	}

	// Загрузка шрифта Times New Roman
	ttfPath := "src/assets/Times New Roman.ttf"
	if _, err := os.Stat(ttfPath); err == nil {
		ttfData, err := os.ReadFile(ttfPath)
		if err != nil {
			g.font = nil
		} else {
			faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(ttfData))
			if err != nil {
				g.font = nil
			} else {
				g.font = &text.GoTextFace{
					Source: faceSource,
					Size:   24,
				}
			}
		}
	} else {
		g.font = nil
	}

	// Создание аудиоконтекста
	audioContext := audio.NewContext(44100)

	// Загрузка аудиофайлов
	for _, audioFile := range []string{"menu.wav", "game.wav", "custom.wav"} {
		path := fmt.Sprintf("src/assets/%s", audioFile)
		if _, err := os.Stat(path); err != nil {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		reader := bytes.NewReader(data)
		d, err := wav.DecodeWithSampleRate(audioContext.SampleRate(), reader)
		if err != nil {
			continue
		}

		player, err := audioContext.NewPlayer(d)
		if err != nil {
			continue
		}

		// Установка разной начальной громкости для разных плееров
		switch audioFile {
		case "menu.wav":
			player.SetVolume(0.1) // Громкость меню
			g.menuPlayer = player
		case "game.wav":
			player.SetVolume(0.1) // Громкость для классической игры
			g.gamePlayer = player
		case "custom.wav":
			player.SetVolume(0.1) // Громкость для пользовательского режима
			g.customPlayer = player
		}
	}

	return nil
}
