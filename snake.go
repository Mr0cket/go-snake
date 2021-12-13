package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"snake/utils"
)

type Position struct {
	X, Y int
}

type Snake []Position
type GameState struct {
	// The snake's body
	Snake Snake
	// Direction of the snake
	direction string
	// The food's state
	Food Position
	// The amount of food the snake has eaten
	Score int
}

// Entry point of the game
func main() {
	switch os.Args[1] {
	case "start":
		startGame()
	default:
		panic("Unknown command")
	}

}

func startGame() {
	fmt.Println("Starting game")

	// Remove the cursor
	fmt.Print("\033[?25l")

	// Get dimensions of the terminal window
	windowSize := utils.GetWindowSize()

	gameState := newGameState(windowSize)

	key := make(chan string)
	go utils.ListenForKeyPress(key)

	var reason string
	var gameOn = true

	// Initialise the frame ticker
	var tick bool
	ticker := time.Tick(time.Second / 10)

	// The game loop
	for gameOn {
		render(windowSize, &gameState)
		tick = false

		// Wait for the next tick
		for !tick {
			select {
			case direction := <-key:
				// Only update the direction if the key is a valid direction
				if utils.IsValidDirection(direction, gameState.direction) {
					gameState.direction = direction
				}
			case <-ticker:
				tick = true
			}
		}

		gameOn, reason = updateGameState(windowSize, &gameState)
	}
	utils.ExitGame(reason)
}

// Initialise game state
func newGameState(windowSize utils.WindowDimensions) GameState {
	snakeHead := Position{X: windowSize.Cols / 2, Y: windowSize.Rows / 2}
	Snake := Snake{snakeHead, {X: snakeHead.X - 1, Y: snakeHead.Y}, {X: snakeHead.X - 2, Y: snakeHead.Y}}

	Food := getFoodPosition(windowSize, Snake)

	return GameState{Snake: Snake, Food: Food, direction: "right", Score: 0}
}

// TODO: Add check for snake collision with wall/food
func updateGameState(windowSize utils.WindowDimensions, gameState *GameState) (bool, string) {
	// Update the snake's position
	snakeBody := gameState.Snake[:]
	snakeHead := snakeBody[0]
	switch gameState.direction {
	case "up":
		snakeHead.Y--
	case "down":
		snakeHead.Y++
	case "left":
		snakeHead.X--
	case "right":
		snakeHead.X++
	}
	// create a new array to append the head to the front of the array
	gameState.Snake = append([]Position{snakeHead}, snakeBody[:len(snakeBody)-1]...)

	/* Game Rules */

	// Snake - Boundary Collision (game over)
	if snakeHead.X < 1 || snakeHead.X > windowSize.Cols-1 || snakeHead.Y < 1 || snakeHead.Y > windowSize.Rows-1 {
		return false, "You hit a wall"
	}

	// Snake - food collision (body grows)
	if snakeHead.X == gameState.Food.X && snakeHead.Y == gameState.Food.Y {
		gameState.Score++
		// Update the food position
		gameState.Food = getFoodPosition(utils.GetWindowSize(), gameState.Snake)
		gameState.Snake = append(gameState.Snake, snakeBody[len(snakeBody)-1])
	}

	// Snake - Self Collision (game over)
	for _, segment := range gameState.Snake[1:] {
		if snakeHead.X == segment.X && snakeHead.Y == segment.Y {
			fmt.Println("You ate yourself")
			return false, "You ate yourself"
		}
	}
	return true, ""
}

// Renders all the elements of the game
func render(winSize utils.WindowDimensions, gameState *GameState) {
	Food := gameState.Food
	// Clear the screen
	fmt.Print(utils.CSI + "2J")

	renderBorder(winSize)
	renderSnake(gameState.Snake)

	// Render the food
	utils.SetPosition("F", Food.X, Food.Y)

	// Render the score
	utils.SetPosition(fmt.Sprintf("Score: %d", gameState.Score), 1, 1)
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

func renderSnake(snake Snake) {
	// Render the snake's head
	utils.SetPosition("X", snake[0].X, snake[0].Y)
	// Render the snake body
	for _, segment := range snake[1:] {
		utils.SetPosition("█", segment.X, segment.Y)
	}
}

// Create a random position for the food within the game window
// TODO: avoid food spawning on the snake
func getFoodPosition(winSize utils.WindowDimensions, snake Snake) Position {
	rand.Seed(time.Now().UnixNano())
	foodPosX := rand.Intn(winSize.Cols-2) + 1
	foodPosY := rand.Intn(winSize.Rows-2) + 1

	// Check if the food is on the snake
	for _, segment := range snake {
		if foodPosX == segment.X && foodPosY == segment.Y {
			return getFoodPosition(winSize, snake)
		}
	}

	return Position{X: foodPosX, Y: foodPosY}
}
