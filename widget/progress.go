package widget

import (
	"fmt"
	"strings"
	"time"
)

type ProgressWidget struct {
	CurLine  int
	Lines    int
	Hint     string
	State    string
	Progress int
	MaxWidth int
	Infinite bool
	Channel  chan string
}

func (i *ProgressWidget) Render() {
	width := i.MaxWidth - 8

	if i.Infinite {
		i.CurLine = 1
		i.Lines = 1
		spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		fmt.Printf("%s %s", spinners[0], i.Hint)

		go func() {
			idx := 0
			for {
				if i.Completed() {
					break
				}
				fmt.Printf("\033[G%s\033[%dC", spinners[idx%len(spinners)], len(i.Hint)+1)
				idx++
				time.Sleep(time.Millisecond * 100)
			}
		}()

	} else {
		fmt.Printf("%s\n", i.Hint)
		fmt.Printf("[\033[%dG] %d%s\n", width+2, i.Progress, "%")
		bars := (i.Progress * width) / 100
		i.CurLine = 2
		if bars > 0 {
			fmt.Printf("\033[1A\033[2G%s\n", strings.Repeat("-", bars))
		}
		i.Lines = 2
	}
}

func (i *ProgressWidget) Clear() {
	if i.CurLine > 0 {
		fmt.Printf("\033[%dA", i.CurLine)
	}
	fmt.Print("\033[G")  // Move to the beginning of the line
	fmt.Print("\033[0J") // Move to the beginning of the line
}

func (i *ProgressWidget) OnKeypress(s string) {
	if i.Infinite {
		if s == "\n" {
			i.State = "completed"
		}
	} else {

		if s == "\x1b[A" {
			i.Progress = min(100, i.Progress+5)
		}

		if s == "\x1b[B" {
			i.Progress = max(0, i.Progress-5)
		}
		if i.Progress == 100 {
			i.State = "completed"
		}

		i.Clear()
		i.Render()
	}
}

func (i *ProgressWidget) Completed() bool {
	return i.State == "completed"
}

func NewProgressWidget(hint string) *ProgressWidget {
	return &ProgressWidget{
		Hint:     hint,
		State:    "progress",
		MaxWidth: 50,
		Lines:    3,
		Progress: 0,
	}
}
