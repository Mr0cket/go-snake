package utils

import (
	"fmt"
	"os"

	"github.com/pkg/term"
	"golang.org/x/sys/unix"
)

const CSI = "\033["

type WindowDimensions struct {
	Rows, Cols int
}

type Result interface{}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func GetWindowSize() WindowDimensions {

	// get console unix file descriptor
	outputFileDesc := int(os.Stdout.Fd())

	winSize, err := unix.IoctlGetWinsize(outputFileDesc, unix.TIOCGWINSZ)
	Must(err)

	return WindowDimensions{int(winSize.Row), int(winSize.Col)}
}

func GetChar() (charString string) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	numRead, err := t.Read(bytes)
	Must(err)

	if numRead == 1 {
		switch bytes[0] {
		case 3:
			charString = "ctrl+c"
		default:
			charString = string(bytes[0])
		}
	} else if bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequences, beginning with "ESC".
		switch bytes[2] {
		case 65: // Up
			charString = "up"
		case 66: // Down
			charString = "down"
		case 67: // Right
			charString = "right"
		case 68: // Left
			charString = "left"
		default: // Ignore any other keypresses
		}
	}
	t.Restore()
	t.Close()
	return
}

func ListenForKeyPress(key chan string) {
	for {
		keycode := GetChar()

		switch keycode {
		case "ctrl+c": // Escape
			ExitGame("Game exited")

		// Send the direction to the channel
		case "up", "down", "left", "right":
			key <- keycode
		default:
		}
	}
}

func ExitGame(reason string) {
	fmt.Println(reason)

	// Restore pointer
	fmt.Print("\033[?25h")

	// Exit the process
	os.Exit(0)
}

func IsValidDirection(newDirection, currentDirection string) bool {
	switch newDirection {
	case "up":
		return currentDirection != "down"
	case "down":
		return currentDirection != "up"
	case "left":
		return currentDirection != "right"
	case "right":
		return currentDirection != "left"
	default:
		return false
	}
}

// Sets the position of a string on the screen
func SetPosition(str string, xPos, yPos int) {
	fmt.Printf(CSI+"%d;%dH%s", yPos, xPos, str)
}
