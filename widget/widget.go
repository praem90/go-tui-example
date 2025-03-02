package widget

type Widget interface {
	Render()
	OnKeypress(s string)
	Clear()
	Completed() bool
}
