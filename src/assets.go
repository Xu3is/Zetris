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
	"log"
	"math/rand"
	"os"
)

// loadAssets загружает изображения, шрифт и аудиофайлы
func (g *Game) loadAssets() error {
	// Загрузка изображений фигур
	for _, shape := range []string{"i", "j", "l", "o", "s", "t", "z"} {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("src/assets/images/%s.png", shape))
		if err != nil {
			log.Printf("Ошибка загрузки %s.png: %v", shape, err)
			img = ebiten.NewImage(cellSize, cellSize)
			img.Fill(color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
		}
		g.images[shape] = img
	}

	// Загрузка изображения клетки поля
	img, _, err := ebitenutil.NewImageFromFile("src/assets/images/boardcell.png")
	if err != nil {
		log.Printf("Ошибка загрузки boardcell.png: %v", err)
		img = ebiten.NewImage(cellSize, cellSize)
		img.Fill(color.RGBA{128, 128, 128, 255})
	}
	g.images["boardcell"] = img

	// Загрузка шрифта Times New Roman
	ttfData, err := os.ReadFile("src/assets/Times New Roman.ttf")
	if err != nil {
		log.Printf("Ошибка загрузки шрифта Times New Roman.ttf: %v. Используется запасной вариант.", err)
		g.font = nil
	} else {
		faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(ttfData))
		if err != nil {
			log.Printf("Ошибка создания шрифта Times New Roman: %v. Используется запасной вариант.", err)
			g.font = nil
		} else {
			g.font = &text.GoTextFace{
				Source: faceSource,
				Size:   24,
			}
		}
	}

	// Создание аудиоконтекста
	audioContext := audio.NewContext(44100)

	// Загрузка аудиофайлов
	for _, audioFile := range []string{"menu.wav", "game.wav", "custom.wav"} {
		path := fmt.Sprintf("src/assets/%s", audioFile)
		// Чтение файла в память
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Ошибка чтения %s: %v", path, err)
			continue
		}

		reader := bytes.NewReader(data)
		d, err := wav.DecodeWithSampleRate(audioContext.SampleRate(), reader)
		if err != nil {
			log.Printf("Ошибка декодирования %s: %v", path, err)
			continue
		}

		player, err := audioContext.NewPlayer(d)
		if err != nil {
			log.Printf("Ошибка создания плеера %s: %v", path, err)
			continue
		}

		// Установка разной начальной громкости для разных плееров
		switch audioFile {
		case "menu.wav":
			player.SetVolume(0.2) // Громкость меню
			g.menuPlayer = player
			log.Printf("menuPlayer успешно создан для %s, громкость: 0.2", path)
		case "game.wav":
			player.SetVolume(0.05) // Сильно уменьшенная громкость для классической игры
			g.gamePlayer = player
			log.Printf("gamePlayer успешно создан для %s, громкость: 0.05", path)
		case "custom.wav":
			player.SetVolume(0.1) // Уменьшенная громкость для пользовательского режима
			g.customPlayer = player
			log.Printf("customPlayer успешно создан для %s, громкость: 0.1", path)
		}
	}

	// Проверка инициализации menuPlayer
	if g.menuPlayer == nil {
		log.Println("Предупреждение: menuPlayer не инициализирован, музыка в меню не будет воспроизводиться")
	}

	return nil
}
