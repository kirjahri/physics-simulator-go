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
}

type Game struct {
	balls []Ball
}

const (
	screenWidth  = 800
	screenHeight = 450
	gravity      = 9.8
	multiplier   = 100
	frictionY    = 0.7
	frictionX    = 0.98
	minBallMass  = 1
	maxBallMass  = 5
	minBallVel   = 1
	maxBallVel   = 5
)

func (b *Ball) Update(dt float32) {
	b.velocity = rl.Vector2Add(b.velocity, rl.Vector2Scale(b.acceleration, dt))
	b.position = rl.Vector2Add(b.position, rl.Vector2Scale(b.velocity, dt))

	if b.position.Y > screenHeight-b.radius {
		b.position.Y = screenHeight - b.radius
		b.velocity.Y *= -frictionY
	}

	b.velocity.X *= frictionX
}

func (b *Ball) ResolveEdgeCollision() {
	if b.position.X > screenWidth-b.radius {
		b.position.X = screenWidth - b.radius
		b.velocity.X *= -1
	} else if b.position.X < b.radius {
		b.position.X = b.radius
		b.velocity.X *= -1
	}

	if b.position.Y < b.radius {
		b.position.Y = b.radius
		b.velocity.Y *= -1
	}
}

func (b *Ball) ResolveBallCollision(otherBall *Ball) {
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

func (g *Game) Init() {
	g.balls = make([]Ball, 0)
}

func (g *Game) Update() {
	dt := rl.GetFrameTime()

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mass := minBallMass + rand.Float32()*(maxBallMass+1-minBallMass)
		g.balls = append(g.balls, Ball{
			position:     rl.GetMousePosition(),
			radius:       float32(math.Sqrt(float64(mass))) * 10,
			color:        rl.Blue,
			mass:         mass,
			acceleration: rl.NewVector2(0, gravity*multiplier),
		})
	}

	for i := range g.balls {
		g.balls[i].Update(dt)
		g.balls[i].ResolveEdgeCollision()

		for j := range g.balls {
			if g.balls[i] == g.balls[j] {
				continue
			}

			g.balls[i].ResolveBallCollision(&g.balls[j])
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
