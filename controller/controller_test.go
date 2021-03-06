package controller

import (
	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/state"
	"testing"
)

func TestControllerStart(t *testing.T) {
	s := state.NewState()
	je := make(chan common.HatEvent)
	se := make(chan *common.DisplayMessage)
	hat := &hatMock{
		je: je,
		se: se,
	}

	done := make(chan bool)
	defer close(done)

	c := &Controller{
		hat:            hat,
		joystickEvents: je,
		screenEvents:   se,
		done:           done,
		state:          s,
	}

	c.Start()

	y := s.Cursor.Y
	x := s.Cursor.X

	hat.MoveDown()
	msg := <-se
	if msg.CursorY != y+1 {
		t.Errorf("msg.CursorY should be %d but it's %d", y+1, msg.CursorY)
	}

	hat.MoveUp()
	msg = <-se
	if msg.CursorY != y {
		t.Errorf("msg.CursorY should be %d but it's %d", y, msg.CursorY)
	}

	hat.MoveRight()
	msg = <-se
	if msg.CursorX != x+1 {
		t.Errorf("msg.CursorX should be %d but it's %d", x+1, msg.CursorY)
	}

	hat.MoveLeft()
	msg = <-se
	if msg.CursorY != x {
		t.Errorf("msg.CursorX should be %d but it's %d", x, msg.CursorY)
	}

	hat.Press()
	msg = <-se
	if !msg.Screen[y][x] {
		t.Errorf("msg.Screen[%d][%d] should be set", y, x)
	}
}
