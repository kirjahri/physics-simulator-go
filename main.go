package main

import (
	"math"
	"math/rand"

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
	balls []Ball
}

const (
	screenWidth  = 800
	screenHeight = 450
	gravity      = 9.8
	multiplier   = 1
	friction     = 0.7
	minimum      = 55
	minBallSize  = 5
	maxBallSize  = 25
)

func (g *Game) Init() {}

func (g *Game) Update() {
	deltaTime := rl.GetFrameTime()

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		g.balls = append(g.balls, Ball{
			position: rl.GetMousePosition(),
			radius:   minBallSize + rand.Float32()*(maxBallSize+1-minBallSize),
			color:    rl.Blue,
		})
	}

	for i, b := range g.balls {
		if !b.isStill {
			threshold := screenHeight - b.radius

			g.balls[i].velocity += gravity * multiplier
			newY := g.balls[i].position.Y + g.balls[i].velocity*deltaTime

			if newY >= threshold {
				// Correct the ball's position by subtracting how far it "penetrates" into the ground
				penetration := newY - threshold
				g.balls[i].position.Y = threshold - penetration

				// If the ball's speed and distance to the ground are low enough, stop it completely
				// This prevents the ball from jittering when close to the ground
				if math.Abs(float64(g.balls[i].velocity)) < minimum && math.Abs(float64(g.balls[i].position.Y-threshold)) < minimum {
					g.balls[i].velocity = 0
					g.balls[i].position.Y = threshold
					g.balls[i].isStill = true
				} else {
					g.balls[i].velocity *= -friction
				}
			} else {
				g.balls[i].position.Y = newY
			}
		}
	}
}

func (g *Game) Draw() {
	for _, b := range g.balls {
		rl.DrawCircleV(b.position, b.radius, b.color)
	}
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

		rl.ClearBackground(rl.White)
		game.Draw()

		rl.EndDrawing()
	}
}
