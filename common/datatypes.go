package common

const (
	// the Sense HAT display is 8X8 matrix
	WindowSize = 8
)

// Joystick events
type HatEvent uint8

const (
	Pressed HatEvent = iota
	MoveUp
	MoveLeft
	MoveDown
	MoveRight
)

// Color is the Color of one pixel in the Canvas
type Color bool

// HAT display events
type DisplayMessage struct {
	Screen  [][]Color
	CursorX uint8
	CursorY uint8
}

func NewDisplayMessage(mat [][]Color, x, y uint8) *DisplayMessage {
	return &DisplayMessage{
		Screen:  mat,
		CursorX: x,
		CursorY: y,
	}
}
