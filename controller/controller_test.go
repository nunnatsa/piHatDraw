package controller

import (
	"encoding/json"
	"testing"

	"github.com/nunnatsa/piHatDraw/notifier"
	"github.com/nunnatsa/piHatDraw/webapp"

	"github.com/nunnatsa/piHatDraw/hat"

	"github.com/nunnatsa/piHatDraw/state"
)

func TestControllerStart(t *testing.T) {
	s := state.NewState(40, 24)
	je := make(chan hat.Event)
	se := make(chan hat.DisplayMessage)
	n := notifier.NewNotifier()
	ce := make(chan webapp.ClientEvent)

	hatMock := &hatMock{
		je: je,
		se: se,
	}

	reg1 := make(chan []byte)
	reg2 := make(chan []byte)

	client1 := n.Subscribe(reg1)
	client2 := n.Subscribe(reg2)

	done := make(chan bool)
	defer close(done)

	c := &Controller{
		hat:            hatMock,
		joystickEvents: je,
		screenEvents:   se,
		done:           done,
		state:          s,
		notifier:       n,
		clientEvents:   ce,
	}

	y := s.Cursor.Y
	x := s.Cursor.X

	c.Start()

	msg := <-se
	if msg.CursorX != x-s.Window.X {
		t.Errorf("msg.CursorX should be %d but it's %d", x-s.Window.X, msg.CursorX)
	}

	ce <- webapp.ClientEventRegistered(client1)
	<-checkNotifications(t, reg1, x, y)
	ce <- webapp.ClientEventRegistered(client2)
	<-checkNotifications(t, reg2, x, y)

	if msg.CursorY != y-s.Window.Y {
		t.Errorf("msg.CursorY should be %d but it's %d", y-s.Window.Y, msg.CursorY)
	}

	hatMock.MoveDown()
	msg = <-se
	if msg.CursorY != y+1-s.Window.Y {
		t.Errorf("msg.CursorY should be %d but it's %d", y+1-s.Window.Y, msg.CursorY)
	}
	<-checkNotifications(t, reg1, x, y+1)
	<-checkNotifications(t, reg2, x, y+1)

	hatMock.MoveUp()
	msg = <-se
	if msg.CursorY != y-s.Window.Y {
		t.Errorf("msg.CursorY should be %d but it's %d", y-s.Window.Y, msg.CursorY)
	}
	<-checkNotifications(t, reg1, x, y)
	<-checkNotifications(t, reg2, x, y)

	hatMock.MoveRight()
	msg = <-se
	if msg.CursorX != x+1-s.Window.X {
		t.Errorf("msg.CursorX should be %d but it's %d", x+1-s.Window.X, msg.CursorY)
	}
	<-checkNotifications(t, reg1, x+1, y)
	<-checkNotifications(t, reg2, x+1, y)

	hatMock.MoveLeft()
	msg = <-se
	if msg.CursorX != x-s.Window.X {
		t.Errorf("msg.CursorX should be %d but it's %d", x-s.Window.X, msg.CursorX)
	}
	<-checkNotifications(t, reg1, x, y)
	<-checkNotifications(t, reg2, x, y)

	hatMock.Press()
	msg = <-se
	if !msg.Screen[y-s.Window.Y][x-s.Window.X] {
		t.Errorf("msg.Screen[%d][%d] should be set", y, x)
	}
	<-checkNotifications(t, reg1, x, y, x, y)
	<-checkNotifications(t, reg2, x, y, x, y)
}

func checkNotifications(t *testing.T, reg chan []byte, x uint8, y uint8, points ...uint8) chan bool {
	doneCheckingNotifier := make(chan bool)
	go func() {
		defer close(doneCheckingNotifier)

		msg := <-reg
		webMsg, err := getCanvasFromMsg(msg)
		if err != nil {
			t.Fatal("getCanvasFromMsg", err)
		}
		if webMsg.Cursor.X != x {
			t.Errorf("webMsg.Cursor.X should be %d but it's %d", x, webMsg.Cursor.X)
		}
		if webMsg.Cursor.Y != y {
			t.Errorf("webMsg.Cursor.y should be %d but it's %d", y+1, webMsg.Cursor.Y)
		}

		if len(points) > 0 && len(points)%2 == 0 {
			for i := 0; i < len(points); i += 2 {
				px := points[i]
				py := points[i+1]

				if !webMsg.Canvas[py][px] {
					t.Error("webMsg.Canvas[py][px] should be set")
				}
			}
		}
	}()
	return doneCheckingNotifier
}

func getCanvasFromMsg(msg []byte) (*state.State, error) {
	s := &state.State{}
	if err := json.Unmarshal(msg, s); err != nil {
		return nil, err
	}
	return s, nil
}
