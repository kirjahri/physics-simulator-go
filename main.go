package main

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Ball struct {
	position raylib.Vector2
	radius   float32
	color    raylib.Color
	velocity float32
}

type Game struct {
	ball Ball
}

func (g *Game) Init() {
	g.ball = Ball{
		position: raylib.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		radius:   20,
		color:    raylib.White,
	}
}

const gravity = 0.098

func (g *Game) Update() {
	g.ball.velocity += gravity
	g.ball.position.Y += g.ball.velocity
	if g.ball.position.Y >= screenHeight {
		g.ball.position.Y = screenHeight
	}
}

func (g *Game) Draw() {
	raylib.DrawCircleV(g.ball.position, g.ball.radius, g.ball.color)
}

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	raylib.InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window")
	defer raylib.CloseWindow()

	raylib.SetTargetFPS(60)

	var game Game
	game.Init()

	for !raylib.WindowShouldClose() {
		game.Update()

		raylib.BeginDrawing()

		raylib.ClearBackground(raylib.Black)
		game.Draw()

		raylib.EndDrawing()
	}
}
