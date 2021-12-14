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
	gameState := models.GameState.New(models.GameState{}, windowSize)

	key := make(chan string)
	go utils.ListenForKeyPress(key)

	var reason string
	var gameOn = true

	// Initialise the frame ticker
	var tick bool
	ticker := time.Tick(time.Second * 4)

	// The game loop
	// While the game is not over:
	for gameOn {
		gameState.Render(windowSize)
		tick = false

		// Get input while waiting for the next frame
		for !tick {
			select {
			case newInput := <-key:
				// Only update the direction if the key is valid with respect to current direction
				if gameState.IsValidDirection(newInput) {
					gameState.Direction = newInput
				}
			case <-ticker:
				tick = true
			}
		}

		gameOn, reason = gameState.Update(windowSize)
	}
	utils.ExitGame(reason)
}
