package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	todo := []string{}
	currText := ""

	ch := make(chan []byte)

	go func() {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "-echo").Run()
		b := make([]byte, 128)
		ch <- []byte("")
		for {
			n, _ := os.Stdin.Read(b)
			ch <- b[:n]
		}

	}()

	for range 11 {
		println(os.Getenv("LINES"))
	}

	for {
		select {
		case b := <-ch:
			key := string(b)

			if len(key) > 0 {
				if key == "\n" {
					if len(currText) > 0 {
						if currText == ":q" {
							fmt.Print("\033[11A")
							fmt.Print("\033[G")
							fmt.Print("\033[J")
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
				fmt.Print("\033[3mEnter your message\033[23m\033[2G")
			}
		}
	}
}
