package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"snake/models"
	"snake/utils"
)

var gameInterval time.Duration

// Entry point of the game
func main() {
	// Initialize the game configuration
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startCmd.DurationVar(&gameInterval, "speed", time.Millisecond*100, "The interval between frames")

	if len(os.Args) < 2 {
		utils.ExitGame("Please specify a command (start)")
	}

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		startGame()
	default:
		panic("Unknown command")
	}
}

func startGame() {
	fmt.Println("Speed:", gameInterval)
	fmt.Println("Starting game...")
	time.Sleep(time.Second)

	// Remove the cursor
	fmt.Print("\033[?25l")

	// Create a new game
	var gameState = models.NewGame()
	var key = make(chan string)
	var reason string
	var gameOn = true
	var nextFrame bool

	// Listen for key presses with a goroutine
	go utils.ListenForKeyPress(key)

	// The game loop
	// While the game is not over:
	for gameOn {
		gameState.Render()
		nextFrame = false

		// Get input while waiting for the next frame
		for !nextFrame {
			select {
			case newInput := <-key:
				// Only update the direction if the key is valid with respect to current direction
				if gameState.IsValidDirection(newInput) {
					gameState.NewDirection = newInput
				}
			case <-time.After(gameInterval):
				nextFrame = true
			}
		}

		gameOn, reason = gameState.Update()
	}
	utils.ExitGame(reason)
}
