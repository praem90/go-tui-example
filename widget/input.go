package widget

import (
	"fmt"
	"net/mail"
	"strings"
)

type InputWidget struct {
	Lines     int
	CurLine   int
	CurCol    int
	Question  string
	Hint      string
	State     string
	Input     string
	MaxWidth  int
	Validator func(*string) bool
}

func (i *InputWidget) Render() {
	width := i.MaxWidth - 2
	i.renderBorder(width)

	fmt.Printf("\033[3A\033[3G %s ", i.Question)

	i.Lines = 3
	i.CurLine = 1

	if i.Validator != nil && len(i.Input) > 0 {
		fmt.Print("\033[2B\033[2G")

		if i.Validator(&i.Input) {
			fmt.Print("\033[38;5;82m")
			fmt.Printf("\033[3G ✓ Valid ")
			i.State = "valid"
		} else {
			fmt.Print("\033[38;5;220m")
			fmt.Print("\033[3G ⚠ InValid ")
			i.State = "warning"
		}
		fmt.Print("\033[0m")
		fmt.Print("\033[2A\033[2G")
	}

	fmt.Printf("\033[1B\033[2G%s", i.Input)
	if len(i.Input) == 0 {
		fmt.Printf("\033[38;5;242m %s \033[0m\033[2G", i.Hint)
	}
}

func (i *InputWidget) renderBorder(width int) {
	if i.State == "warning" {
		fmt.Print("\033[38;5;220m")
	}
	if i.State == "answered" || i.State == "valid" {
		fmt.Print("\033[38;5;82m")
	}
	fmt.Printf("╭%s╮\n", strings.Repeat("─", width))
	fmt.Printf("│%s│\n", strings.Repeat(" ", width))
	fmt.Printf("╰%s╯\n", strings.Repeat("─", width))

	fmt.Print("\033[0m")
}

func (i *InputWidget) Clear() {
	if i.CurLine > 0 {
		fmt.Printf("\033[%dA", i.CurLine)
	}
	fmt.Print("\033[G")  // Move to the beginning of the line
	fmt.Print("\033[0J") // Move to the beginning of the line
}

func (i *InputWidget) OnKeypress(s string) {
	if s == "\n" {
		i.Clear()
		i.Render()
		i.State = "answered"
		return
	}

	if strings.HasPrefix(s, "\x1b") {
		return
	}

	i.Input += s
	i.Clear()
	i.Render()
}

func (i *InputWidget) Completed() bool {
	return i.State == "answered"
}

func NewTextInputWidget(question string, hint string) *InputWidget {
	return &InputWidget{
		Question: question,
		Hint:     hint,
		State:    "",
		Input:    "",
		MaxWidth: 40,
		Lines:    3,
	}
}

func NewEmailInputWidget(question string, hint string) *InputWidget {
	return &InputWidget{
		Question: question,
		Hint:     hint,
		State:    "",
		Input:    "",
		MaxWidth: 40,
		Lines:    3,
		Validator: func(s *string) bool {
			if _, err := mail.ParseAddress(*s); err != nil {
				return false
			}

			return true
		},
	}
}
