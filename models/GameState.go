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
	Direction    string
	NewDirection string
	// The food's state
	Food utils.Position
	// The amount of food the snake has eaten
	Score int
}

// Get dimensions of the terminal window
var windowSize = utils.GetWindowSize()

// Initialise game state
func NewGame() *GameState {
	// Get dimensions of the terminal window

	Snake := NewSnake()

	g := GameState{Snake: *Snake, Direction: "right", NewDirection: "right", Score: 0}
	g.newFood()
	return &g
}

// Update the game state
func (g *GameState) Update() (bool, string) {
	g.Direction = g.NewDirection
	oldSnake := g.Snake.body
	g.Snake.Move(g.Direction)

	/* Check Game Rules */
	// Snake - Boundary Collision (game over)
	if g.Snake.HitWall(windowSize) {
		return false, "You hit a wall"
	}

	// Snake - food collision (body grows)
	if g.Snake.HitPoint(g.Food) {
		g.Score++
		// Update the food position
		g.newFood()
		g.Snake.Append(oldSnake[len(oldSnake)-1])
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
func (g *GameState) newFood() {
	rand.Seed(time.Now().UnixNano())
	foodPosX := rand.Intn(windowSize.Cols-2) + 1
	foodPosY := rand.Intn(windowSize.Rows-2) + 1

	// Check if the food is on the snake
	for _, segment := range g.Snake.body {
		if foodPosX == segment.X && foodPosY == segment.Y {
			g.newFood()
			return
		}
	}

	g.Food = utils.Position{X: foodPosX, Y: foodPosY}
}

// Renders all the elements of the game
func (g GameState) Render() {
	// Clear the screen
	fmt.Print(utils.CSI + "2J")

	// Render the border
	renderBorder()

	// Render the snake
	g.Snake.Render()

	// Render the food
	utils.SetPosition("F", g.Food.X, g.Food.Y)

	// Render the score
	utils.SetPosition(fmt.Sprintf("Score: %d", g.Score), 1, 1)
}

func renderBorder() {
	// Render the top border
	topBorder := "┌" + strings.Repeat("─", windowSize.Cols-2) + "┐"
	utils.SetPosition(topBorder, 1, 1)

	// Render Side borders
	for i := 0; i < windowSize.Rows-2; i++ {
		utils.SetPosition("│", 1, i+2)
		utils.SetPosition("│", windowSize.Cols, i+2)
	}

	// Render the bottom border
	bottomBorder := "└" + strings.Repeat("─", windowSize.Cols-2) + "┘"
	utils.SetPosition(bottomBorder, 1, windowSize.Rows)
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
