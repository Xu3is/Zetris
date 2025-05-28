package main

import (
	"github.com/Xu3is/Zetris/src"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	game, err := src.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(src.ScreenWidth, src.ScreenHeight)
	ebiten.SetWindowTitle("Zetris")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
