package models

import (
	"fmt"
	"math/rand"
	"snake/utils"
	"strings"
	"time"
)

type GameState struct {
	// The snake's body
	Snake Snake
	// Direction of the snake
	Direction string
	// The food's state
	Food utils.Position
	// The amount of food the snake has eaten
	Score int
}

// Initialise game state
func (g GameState) New(windowSize utils.WindowDimensions) GameState {
	Snake := Snake.New(Snake{}, windowSize)

	g = GameState{Snake: Snake, Direction: "right", Score: 0}
	g.Food = g.newFood(windowSize)
	return g
}

// Update the game state
func (g GameState) Update(windowSize utils.WindowDimensions) (bool, string) {
	oldSnake := g.Snake.body
	g.Snake.Move(g.Direction)

	/* Game Rules */
	// Snake - Boundary Collision (game over)
	if g.Snake.HitWall(windowSize) {
		return false, "You hit a wall"
	}

	// Snake - food collision (body grows)
	if g.Snake.HitPoint(g.Food) {
		g.Score++
		// Update the food position
		g.Food = g.newFood(windowSize)
		g.Snake.Add(oldSnake[len(oldSnake)-1])
	}

	// Snake - Self Collision (game over)
	for _, segment := range oldSnake[1:] {
		if g.Snake.HitPoint(segment) {
			return false, "You ate yourself"
		}
	}
	return true, ""
}

// Create a random position for the food within the game window
func (g GameState) newFood(winSize utils.WindowDimensions) utils.Position {
	rand.Seed(time.Now().UnixNano())
	foodPosX := rand.Intn(winSize.Cols-2) + 1
	foodPosY := rand.Intn(winSize.Rows-2) + 1

	// Check if the food is on the snake
	for _, segment := range g.Snake.body {
		if foodPosX == segment.X && foodPosY == segment.Y {
			return g.newFood(winSize)
		}
	}

	return utils.Position{X: foodPosX, Y: foodPosY}
}

// Renders all the elements of the game
func (g GameState) Render(winSize utils.WindowDimensions) {
	// Clear the screen
	fmt.Print(utils.CSI + "2J")

	// Render the border
	renderBorder(winSize)

	// Render the snake
	g.Snake.Render()

	// Render the food
	utils.SetPosition("F", g.Food.X, g.Food.Y)

	// Render the score
	utils.SetPosition(fmt.Sprintf("Score: %d", g.Score), 1, 1)
}

func renderBorder(winSize utils.WindowDimensions) {
	// Render the top border
	topBorder := "┌" + strings.Repeat("─", winSize.Cols-2) + "┐"
	utils.SetPosition(topBorder, 1, 1)

	// Render Side borders
	for i := 0; i < winSize.Rows-2; i++ {
		utils.SetPosition("│", 1, i+2)
		utils.SetPosition("│", winSize.Cols, i+2)
	}

	// Render the bottom border
	bottomBorder := "└" + strings.Repeat("─", winSize.Cols-2) + "┘"
	utils.SetPosition(bottomBorder, 1, winSize.Rows)
}

func (g GameState) IsValidDirection(newDirection string) bool {
	switch newDirection {
	case "up":
		return g.Direction != "down"
	case "down":
		return g.Direction != "up"
	case "left":
		return g.Direction != "right"
	case "right":
		return g.Direction != "left"
	default:
		return false
	}
}
