package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ball struct {
	position rl.Vector2
	radius   float32
	color    rl.Color
	velocity float32
	isStill  bool
}

type Game struct {
	ball Ball
}

const (
	screenWidth  = 800
	screenHeight = 450
	gravity      = 9.8
	multiplier   = 1
	friction     = 0.85
	minimum      = 55
)

func (g *Game) Init() {
	g.ball = Ball{
		position: rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		radius:   20,
		color:    rl.White,
	}
}

func (g *Game) Update() {
	if !g.ball.isStill {
		deltaTime := rl.GetFrameTime()
		threshold := screenHeight - g.ball.radius

		g.ball.velocity += gravity * multiplier
		newY := g.ball.position.Y + g.ball.velocity*deltaTime

		if newY >= threshold {
			// Correct the ball's position by subtracting how far it "penetrates" into the ground
			penetration := newY - threshold
			g.ball.position.Y = threshold - penetration

			// If the ball's speed and distance to the ground are low enough, stop it completely
			// This prevents the ball from jittering when close to the ground
			if math.Abs(float64(g.ball.velocity)) < minimum && math.Abs(float64(g.ball.position.Y-threshold)) < minimum {
				g.ball.velocity = 0
				g.ball.position.Y = threshold
				g.ball.isStill = true
			} else {
				g.ball.velocity *= -friction
			}
		} else {
			g.ball.position.Y = newY
		}
	}
}

func (g *Game) Draw() {
	rl.DrawCircleV(g.ball.position, g.ball.radius, g.ball.color)
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowHighdpi)

	rl.InitWindow(screenWidth, screenHeight, "raylib")
	defer rl.CloseWindow()

	rl.SetTargetFPS(160)

	var game Game
	game.Init()

	for !rl.WindowShouldClose() {
		game.Update()

		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		game.Draw()

		rl.EndDrawing()
	}
}
