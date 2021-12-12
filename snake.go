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

type Snake struct {
	// The snake body segments
	body []Position
	// The snake's direction (, 2 = down, 3 = left, 4 = right)
	direction string
}

type GameState struct {
	// The snake's state
	Snake Snake
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

	gameOn := true

	// Initialise the frame ticker
	var ticked bool
	tick := time.Tick(time.Second / 2)
	
	// Start the game loop
	for gameOn {
		render(windowSize, &gameState)
		ticked = false

		// Wait for the next tick
		for !ticked {
			select {
				case direction := <- key:
					gameState.Snake.direction = direction
				case <- tick:
					ticked = true
			}
		}
		updateGameState(windowSize, &gameState)
		}
	// Restore the cursor
	fmt.Print("\033[?25h")
}

// Initialise game state
func newGameState(windowSize utils.WindowDimensions) GameState {
	snakeHead := Position{X: windowSize.Cols / 2, Y: windowSize.Rows / 2}
	snakeBody := []Position{snakeHead,{X: snakeHead.X - 1, Y: snakeHead.Y}, {X: snakeHead.X - 2, Y: snakeHead.Y}}
	snake := Snake{body: snakeBody, direction: "right" }
	food := getFoodPosition(windowSize, &snake)

	return GameState{Snake: snake, Food: food, Score: 0}
}

// TODO: Add check for snake collision with wall/food
func updateGameState(windowSize utils.WindowDimensions, gameState *GameState) {
	// Update the snake's position
	snakeBody := gameState.Snake.body[:]
	snakeHead := snakeBody[0]
		switch gameState.Snake.direction {
			case "up": 
				snakeHead.Y--
			case "down": 
				snakeHead.Y++
			case "left": 
				snakeHead.X--
			case "right":
				snakeHead.X++
		}
	// Push the head to the front of the snake
	// create a new array to append the head to the front of the array
		gameState.Snake.body = append([]Position{snakeHead}, snakeBody[:len(snakeBody) - 1]...)

		// Check for boundary collision
		if snakeHead.X < 1 || snakeHead.X > windowSize.Cols - 1 || snakeHead.Y < 1 || snakeHead.Y > windowSize.Rows - 1 {
			fmt.Println("You hit a wall")
			os.Exit(0)
		}

	// Check if the snake has eaten the food
	if snakeHead.X == gameState.Food.X && snakeHead.Y == gameState.Food.Y {
		gameState.Score++
		gameState.Food = getFoodPosition(utils.GetWindowSize(), &gameState.Snake)
		gameState.Snake.body = append(gameState.Snake.body, snakeBody[len(snakeBody) - 1])
	}
}

// Renders all the elements of the game
func render(winSize utils.WindowDimensions, gameState *GameState) {
	Food := gameState.Food
	// Clear the screen
	fmt.Print(utils.CSI + "2J")
	
	renderBorder(winSize)
	renderSnake(&gameState.Snake)

	// Render the food
	SetPosition("F", Food.X, Food.Y)
}	

func renderBorder(winSize utils.WindowDimensions) {
		// Render the top border
		topBorder := "┌" + strings.Repeat("─", winSize.Cols - 2) + "┐"
		SetPosition(topBorder, 1, 1)
	
		// Render Side borders
		for i := 0; i < winSize.Rows-2; i++ {
			SetPosition( "│", 1, i+2)
			SetPosition("│", winSize.Cols, i+2)
		}
		// Render the bottom border
		bottomBorder := "└" + strings.Repeat("─", winSize.Cols - 2) + "┘"
		SetPosition(bottomBorder, 1, winSize.Rows)
}

func renderSnake(snake *Snake) {
	// Render the snake's head
	SetPosition("X", snake.body[0].X, snake.body[0].Y)
	// Render the snake body
	for _, segment := range snake.body[1:] {
		SetPosition("█", segment.X, segment.Y)
	}
}

// Create a random position for the food within the game window
// TODO: avoid food spawning on the snake
func getFoodPosition(winSize utils.WindowDimensions, snake *Snake) Position {
	rand.Seed(time.Now().UnixNano())
	foodPosX := rand.Intn(winSize.Cols - 2) + 1
	foodPosY := rand.Intn(winSize.Rows - 2) + 1

	return Position{X: foodPosX, Y: foodPosY}
}

// Sets the position of a string on the screen
func SetPosition(str string, xPos, yPos int) {
	fmt.Printf(utils.CSI + "%d;%dH%s", yPos, xPos, str)
}
