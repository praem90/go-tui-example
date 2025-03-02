package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/praem90/tui/widget"
)

func main() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "-echo").Run()
	b := make([]byte, 128)
	todo := []string{}
	currText := ""

	renderTextInput()
	renderEmailInput()
	renderSelectInput()
	renderProgressInput()
	renderInfiniteProgressInput()

	for range 11 {
		println(os.Getenv("LINES"))
	}

	for {
		n, _ := os.Stdin.Read(b)
		key := string(b[:n])

		if len(key) > 0 {
			if key == "\n" {
				if len(currText) > 0 {
					if currText == ":q" {
						fmt.Print("\033[11A")
						fmt.Print("\033[G")
						fmt.Print("\033[J")
						reset()
						os.Exit(1)
					}
					todo = append(todo, currText)
					currText = ""
				}
			} else {
				if strings.HasPrefix(key, "\x1b[") {
					continue
				} else {
					switch b[0] {
					case 127:
						if len(currText) > 0 {
							currText = currText[:len(currText)-1]
						}
					default:
						currText += key
					}
				}
			}
		}

		fmt.Print("\033[11A")
		fmt.Print("\033[G")
		fmt.Print("\033[J")

		for i := 9; i >= 0; i-- {
			if len(todo)-1-i >= 0 {
				println(todo[len(todo)-1-i])
			} else {
				println()
			}
		}

		width := max(40, len(currText)+1)

		fmt.Printf("╭%s╮\n", strings.Repeat("─", width))
		fmt.Printf("│%s│\n", strings.Repeat(" ", width))
		fmt.Printf("╰%s╯", strings.Repeat("─", width))

		fmt.Printf("\033[1A\033[2G%s", currText)
		if len(currText) == 0 {
			fmt.Print("\033[38;5;242mEnter your message\033[0m\033[2G")
		}
	}
}

func reset() {
	println("Restting")
	exec.Command("stty", "-F", "/dev/tty", "-cbreak", "echo").Run()
}

func renderTextInput() {
	textInput := widget.NewTextInputWidget("What is your name?", "Mohan Raj")
	RenderWidget(textInput)
}

func renderEmailInput() {
	textInput := widget.NewEmailInputWidget("What is your email?", "you@example.com")
	RenderWidget(textInput)
}

func renderSelectInput() {
	textInput := widget.NewSelectInputWidget(
		"Languages",
		[]string{"Tamil", "English", "PHP", "Golang", "Js", "Rust"},
		true,
	)

	RenderWidget(textInput)
}

func renderProgressInput() {
	textInput := widget.NewProgressWidget(
		"Downloading data...",
	)
	RenderWidget(textInput)
}
func renderInfiniteProgressInput() {
	textInput := widget.NewProgressWidget(
		"Processing...",
	)
	textInput.Infinite = true
	RenderWidget(textInput)
}

func RenderWidget(w widget.Widget) {
	w.Render()
	b := make([]byte, 128)
	for {
		n, _ := os.Stdin.Read(b)
		key := string(b[:n])
		if n > 0 {
			w.OnKeypress(key)
			if w.Completed() {
				w.Clear()
				break
			}
		}
	}
}
