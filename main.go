package main

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	game "github.com/technologyfreak/hnefatafl/game"
)

const (
	totalBlackPawns = 24
	totalWhitePawns = 12
)

func main() {
	game := game.Game{}
	game.Init()

	raylib.InitWindow(int32(game.ScreenWidth), int32(game.ScreenHeight), "Hnefatafl")
	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		game.Update()
		game.Draw()
	}

	raylib.CloseWindow()
}
