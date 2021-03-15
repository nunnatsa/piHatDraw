package controller

import (
	"testing"

	"github.com/nunnatsa/piHatDraw/hat"

	"github.com/nunnatsa/piHatDraw/state"
)

func TestControllerStart(t *testing.T) {
	s := state.NewState()
	je := make(chan hat.Event)
	se := make(chan hat.DisplayMessage)
	hatMock := &hatMock{
		je: je,
		se: se,
	}

	done := make(chan bool)
	defer close(done)

	c := &Controller{
		hat:            hatMock,
		joystickEvents: je,
		screenEvents:   se,
		done:           done,
		state:          s,
	}

	y := s.Cursor.Y
	x := s.Cursor.X

	c.Start()
	msg := <-se
	if msg.CursorX != x {
		t.Errorf("msg.CursorX should be %d but it's %d", x, msg.CursorX)
	}
	if msg.CursorY != y {
		t.Errorf("msg.CursorY should be %d but it's %d", y, msg.CursorY)
	}

	hatMock.MoveDown()
	msg = <-se
	if msg.CursorY != y+1 {
		t.Errorf("msg.CursorY should be %d but it's %d", y+1, msg.CursorY)
	}

	hatMock.MoveUp()
	msg = <-se
	if msg.CursorY != y {
		t.Errorf("msg.CursorY should be %d but it's %d", y, msg.CursorY)
	}

	hatMock.MoveRight()
	msg = <-se
	if msg.CursorX != x+1 {
		t.Errorf("msg.CursorX should be %d but it's %d", x+1, msg.CursorY)
	}

	hatMock.MoveLeft()
	msg = <-se
	if msg.CursorY != x {
		t.Errorf("msg.CursorX should be %d but it's %d", x, msg.CursorY)
	}

	hatMock.Press()
	msg = <-se
	if !msg.Screen[y][x] {
		t.Errorf("msg.Screen[%d][%d] should be set", y, x)
	}
}
