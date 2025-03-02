package widget

import (
	"fmt"
	"strings"
)

type SelectWidget struct {
	Lines          int
	CurLine        int
	Question       string
	Options        []Option
	State          string
	HighlighedItem int
	MaxWidth       int
	Multiple       bool
}

type Option struct {
	label    string
	value    int
	selected bool
}

func (i *SelectWidget) Render() {
	width := i.MaxWidth - 2
	fmt.Printf("╭ %s %s╮\n", i.Question, strings.Repeat("─", width-2-len(i.Question)))
	for idx, option := range i.Options {
		selected := " "
		if option.selected {
			selected = "✓"
			fmt.Printf("\033[48;5;28m")
			fmt.Printf("\033[38;5;235m")
			if idx == i.HighlighedItem {
				fmt.Printf("\033[38;5;248m")
			}
		} else if idx == i.HighlighedItem {
			fmt.Printf("\033[48;5;235m")
			fmt.Printf("\033[38;5;28m")
		}
		fmt.Printf(
			"│%s %s %s│\n",
			selected,
			option.label,
			strings.Repeat(" ", width-3-len(option.label)),
		)
		fmt.Printf("\033[0m")
	}
	fmt.Printf("╰%s╯\n", strings.Repeat("─", width))

	i.Lines = len(i.Options) + 2
	i.CurLine = 0
}

func (i *SelectWidget) Clear() {
	fmt.Printf("\033[%dA", i.Lines)
	fmt.Print("\033[G") // Move to the beginning of the line
	fmt.Print("\033[0J")
}

func (i *SelectWidget) OnKeypress(s string) {
	if s == "\n" {
		i.State = "answered"
		return
	}

	if s == " " {
		for idx, _ := range i.Options {
			if i.HighlighedItem == idx {
				if i.Options[idx].selected {
					i.Options[idx].selected = false
				} else {
					i.Options[idx].selected = true
				}
			} else if i.Multiple == false {
				i.Options[idx].selected = false
			}
		}
	}

	if s == "\x1b[A" {
		i.HighlighedItem = max(0, i.HighlighedItem-1)
	}

	if s == "\x1b[B" {
		i.HighlighedItem = min(len(i.Options)-1, i.HighlighedItem+1)
	}

	i.Clear()
	i.Render()
}

func (i *SelectWidget) Completed() bool {
	return i.State == "answered"
}

func NewSelectInputWidget(question string, options []string, multiple bool) *SelectWidget {
	size := len(options)
	opts := make([]Option, size)
	for i, option := range options {
		opts[i] = Option{
			label:    option,
			value:    i,
			selected: false,
		}
	}
	return &SelectWidget{
		Question: question,
		Options:  opts,
		State:    "unanswered",
		MaxWidth: 40,
		Lines:    len(options) + 2,
		Multiple: multiple,
	}
}
