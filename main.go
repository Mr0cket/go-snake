package main

import (
	"fmt"
	"os"
	"time"

	"snake/models"
	"snake/utils"
)

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
	gameState := models.NewGame(windowSize)

	key := make(chan string)
	go utils.ListenForKeyPress(key)

	var reason string
	var gameOn = true

	// Initialise the frame ticker
	var nextFrame bool
	var interval = time.Millisecond * 100

	// The game loop
	// While the game is not over:
	for gameOn {
		gameState.Render(windowSize)
		nextFrame = false

		// Get input while waiting for the next frame
		for !nextFrame {
			select {
			case newInput := <-key:
				// Only update the direction if the key is valid with respect to current direction
				if gameState.IsValidDirection(newInput) {
					gameState.NewDirection = newInput
				}
			case <-time.After(interval):
				nextFrame = true
			}
		}

		gameOn, reason = gameState.Update(windowSize)
	}
	utils.ExitGame(reason)
}
