package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ball struct {
	position     rl.Vector2
	radius       float32
	color        rl.Color
	mass         float32
	acceleration rl.Vector2
	velocity     rl.Vector2
	isStill      bool
}

type Game struct {
	balls []Ball
}

const (
	screenWidth  = 800
	screenHeight = 450
	gravity      = 9.8
	multiplier   = 100
	friction     = 0.7
	minimum      = 55
	minBallMass  = 1
	maxBallMass  = 5
	minBallVel   = 1
	maxBallVel   = 5
)

func (b *Ball) Update() {
	dt := rl.GetFrameTime()
	b.velocity = rl.Vector2Add(b.velocity, rl.Vector2Scale(b.acceleration, dt))
	b.position = rl.Vector2Add(b.position, rl.Vector2Scale(b.velocity, dt))
}

func (b *Ball) CheckEdgeCollision() {
	if b.position.X > screenWidth-b.radius {
		b.position.X = screenWidth - b.radius
		b.velocity.X *= -1
	} else if b.position.X < b.radius {
		b.position.X = b.radius
		b.velocity.X *= -1
	}

	if b.position.Y > screenHeight-b.radius {
		b.position.Y = screenHeight - b.radius
		b.velocity.Y *= -1
	} else if b.position.Y < b.radius {
		b.position.Y = b.radius
		b.velocity.Y *= -1
	}
}

func (b *Ball) CheckBallCollision(otherBall *Ball) {
	impactVector := rl.Vector2Subtract(otherBall.position, b.position)
	distance := rl.Vector2Length(impactVector)

	if distance < b.radius+otherBall.radius {
		// Correct the overlap when the balls collide
		overlap := distance - (b.radius + otherBall.radius)
		direction := rl.Vector2Scale(rl.Vector2Normalize(impactVector), overlap*0.5)
		b.position = rl.Vector2Add(b.position, direction)
		otherBall.position = rl.Vector2Subtract(otherBall.position, direction)

		// Correct the distance
		distance = b.radius + otherBall.radius
		impactVector = rl.Vector2Scale(rl.Vector2Normalize(impactVector), distance)

		massSum := b.mass + otherBall.mass
		velDiff := rl.Vector2Subtract(otherBall.velocity, b.velocity)
		num := rl.Vector2DotProduct(velDiff, impactVector)
		den := massSum * distance * distance

		// Ball 0 (b)
		deltaVel0 := rl.Vector2Scale(impactVector, 2*otherBall.mass*num/den)
		b.velocity = rl.Vector2Add(b.velocity, deltaVel0)

		// Ball 1 (otherBall)
		deltaVel1 := rl.Vector2Scale(impactVector, -2*b.mass*num/den)
		otherBall.velocity = rl.Vector2Add(otherBall.velocity, deltaVel1)

		// speed0 := rl.Vector2Length(b.velocity)
		// speed1 := rl.Vector2Length(otherBall.velocity)
		// kin0 := 0.5 * b.mass * speed0 * speed0
		// kin1 := 0.5 * otherBall.mass * speed1 * speed1
		// fmt.Println(kin0 + kin1)
	}
}

func (g *Game) Init() {}

func (g *Game) Update() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mass := minBallMass + rand.Float32()*(maxBallMass+1-minBallMass)
		g.balls = append(g.balls, Ball{
			position: rl.GetMousePosition(),
			radius:   float32(math.Sqrt(float64(mass))) * 10,
			color:    rl.Blue,
			mass:     mass,
			velocity: rl.NewVector2((minBallVel+rand.Float32()*(maxBallVel+1-minBallVel))*100, (minBallVel+rand.Float32()*(maxBallVel+1-minBallVel))*100),
		})
	}

	for i := range g.balls {
		// if !b.isStill {
		// 	threshold := screenHeight - b.radius
		//
		// 	g.balls[i].velocity.Y += gravity * multiplier * dt
		// 	newY := g.balls[i].position.Y + g.balls[i].velocity.Y*dt
		//
		// 	if newY >= threshold {
		// 		// Correct the ball's position by subtracting how far it "penetrates" into the ground
		// 		penetration := newY - threshold
		// 		g.balls[i].position.Y = threshold - penetration
		//
		// 		// If the ball's speed and distance to the ground are low enough, stop it completely
		// 		// This prevents the ball from jittering when close to the ground
		// 		if math.Abs(float64(g.balls[i].velocity.Y)) < minimum && math.Abs(float64(g.balls[i].position.Y-threshold)) < minimum {
		// 			g.balls[i].velocity.Y = 0
		// 			g.balls[i].position.Y = threshold
		// 			g.balls[i].isStill = true
		// 		} else {
		// 			g.balls[i].velocity.Y *= -friction
		// 		}
		// 	} else {
		// 		g.balls[i].position.Y = newY
		// 	}
		// }

		g.balls[i].Update()
		g.balls[i].CheckEdgeCollision()

		for j := range g.balls {
			if g.balls[i] == g.balls[j] {
				continue
			}

			g.balls[i].CheckBallCollision(&g.balls[j])
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
