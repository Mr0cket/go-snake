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

func  Must(err error) {
	if err != nil {
		panic(err)
	}
}

func GetWindowSize() WindowDimensions {

	// get console unix file descriptor
	outputFileDesc := int(os.Stdout.Fd())

	winSize, err := unix.IoctlGetWinsize(outputFileDesc, unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
	}
	
	return WindowDimensions{int(winSize.Row), int(winSize.Col)}
}


func GetChar() (charString string, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	numRead, err := t.Read(bytes)
	if err != nil {
		return
	}
	if (numRead == 1) {
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
						// Right
						charString = "right"
						case 68:
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
		keycode, err := GetChar()
		if err != nil {
			panic(err)
		}
		switch keycode {
			case "ctrl+c": // Escape
				fmt.Print("\033[?25h")
				os.Exit(0)
			case "up":
				key <- keycode
			case "down":
				key <- keycode
			case "left":
				key <- keycode
			case "right":
				key <- keycode
			default:
		}
	}
}