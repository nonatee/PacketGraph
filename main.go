package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := Init()
	ebiten.RunGame(game)
	close(game.packetChan)
}
