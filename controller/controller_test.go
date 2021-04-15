package controller

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/hat"
	"github.com/nunnatsa/piHatDraw/notifier"
	"github.com/nunnatsa/piHatDraw/state"
	"github.com/nunnatsa/piHatDraw/webapp"
)

const (
	eraserToolName = "eraser"
	penToolName    = "pen"
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
	err := <-checkMoveNotifications(reg1, x, y)
	if err != nil {
		t.Fatal(err)
	}
	ce <- webapp.ClientEventRegistered(client2)
	err = <-checkMoveNotifications(reg2, x, y)
	if err != nil {
		t.Fatal(err)
	}

	if msg.CursorY != y-s.Window.Y {
		t.Errorf("msg.CursorY should be %d but it's %d", y-s.Window.Y, msg.CursorY)
	}

	ce <- webapp.ClientEventUndo(true)
	if len(se) != 0 || len(reg1) != 0 || len(reg2) != 0 {
		t.Error("should not initiate a chenage")
	}

	hatMock.MoveDown()
	msg = <-se
	if msg.CursorY != y+1-s.Window.Y {
		t.Errorf("msg.CursorY should be %d but it's %d", y+1-s.Window.Y, msg.CursorY)
	}
	err = <-checkMoveNotifications(reg1, x, y+1)
	if err != nil {
		t.Fatal(err)
	}
	err = <-checkMoveNotifications(reg2, x, y+1)
	if err != nil {
		t.Fatal(err)
	}

	ce <- webapp.ClientEventUndo(true)
	if len(se) != 0 || len(reg1) != 0 || len(reg2) != 0 {
		t.Error("should not initiate a chenage")
	}

	hatMock.MoveUp()
	msg = <-se
	if msg.CursorY != y-s.Window.Y {
		t.Errorf("msg.CursorY should be %d but it's %d", y-s.Window.Y, msg.CursorY)
	}
	err = <-checkMoveNotifications(reg1, x, y)
	if err != nil {
		t.Fatal(err)
	}
	err = <-checkMoveNotifications(reg2, x, y)
	if err != nil {
		t.Fatal(err)
	}

	ce <- webapp.ClientEventUndo(true)
	if len(se) != 0 || len(reg1) != 0 || len(reg2) != 0 {
		t.Error("should not initiate a chenage")
	}

	hatMock.MoveRight()
	msg = <-se
	if msg.CursorX != x+1-s.Window.X {
		t.Errorf("msg.CursorX should be %d but it's %d", x+1-s.Window.X, msg.CursorY)
	}
	err = <-checkMoveNotifications(reg1, x+1, y)
	if err != nil {
		t.Fatal(err)
	}

	err = <-checkMoveNotifications(reg2, x+1, y)
	if err != nil {
		t.Fatal(err)
	}

	ce <- webapp.ClientEventUndo(true)
	if len(se) != 0 || len(reg1) != 0 || len(reg2) != 0 {
		t.Error("should not initiate a change")
	}

	hatMock.MoveLeft()
	msg = <-se
	if msg.CursorX != x-s.Window.X {
		t.Errorf("msg.CursorX should be %d but it's %d", x-s.Window.X, msg.CursorX)
	}
	err = <-checkMoveNotifications(reg1, x, y)
	if err != nil {
		t.Fatal(err)
	}

	err = <-checkMoveNotifications(reg2, x, y)
	if err != nil {
		t.Fatal(err)
	}

	ce <- webapp.ClientEventUndo(true)
	if len(se) != 0 || len(reg1) != 0 || len(reg2) != 0 {
		t.Error("should not initiate a change")
	}

	hatMock.Press()
	msg = <-se
	if msg.Screen[y-s.Window.Y][x-s.Window.X] != 0xFFFFFF {
		t.Errorf("msg.Screen[%d][%d] should be set", y, x)
	}
	<-checkPaintNotifications(t, reg1, state.Pixel{X: x, Y: y, Color: 0xFFFFFF})
	<-checkPaintNotifications(t, reg2, state.Pixel{X: x, Y: y, Color: 0xFFFFFF})

	clr := common.Color(0x123456)
	ce <- webapp.ClientEventSetColor(clr)
	<-checkNotificationsColor(t, reg1, &clr, "")
	<-checkNotificationsColor(t, reg2, &clr, "")

	ce <- webapp.ClientEventSetTool(eraserToolName)
	<-checkNotificationsColor(t, reg1, nil, eraserToolName)
	<-checkNotificationsColor(t, reg2, nil, eraserToolName)

	clr = common.Color(0x654321)
	ce <- webapp.ClientEventSetColor(clr)
	<-checkNotificationsColor(t, reg1, &clr, "")
	<-checkNotificationsColor(t, reg2, &clr, "")

	ce <- webapp.ClientEventSetTool(penToolName)
	<-checkNotificationsColor(t, reg1, nil, penToolName)
	<-checkNotificationsColor(t, reg2, nil, penToolName)

	ce <- webapp.ClientEventUndo(true)
	<-checkPaintNotifications(t, reg1, state.Pixel{X: x, Y: y, Color: 0})
	<-checkPaintNotifications(t, reg2, state.Pixel{X: x, Y: y, Color: 0})

	hatMock.Press()
	msg = <-se
	if msg.Screen[y-s.Window.Y][x-s.Window.X] != 0xFFFFFF {
		t.Errorf("msg.Screen[%d][%d] should be set", y, x)
	}
	<-checkPaintNotifications(t, reg1, state.Pixel{X: x, Y: y, Color: 0x654321})
	<-checkPaintNotifications(t, reg2, state.Pixel{X: x, Y: y, Color: 0x654321})

	ce <- webapp.ClientEventReset(true)
	initColor := common.Color(0xFFFFFF)
	<-checkNotificationsColor(t, reg1, &initColor, penToolName)
	err = <-checkMoveNotifications(reg2, x, y)
	if err != nil {
		t.Fatal(err)
	}

	ce <- webapp.ClientEventUndo(true)
	ns := state.NewState(40, 24)
	ns.Canvas[12][20] = 0x654321
	err = <-checkResetNotifications(reg1, ns.Canvas)
	if err != nil {
		t.Fatal(err)
	}
	err = <-checkResetNotifications(reg2, ns.Canvas)
	if err != nil {
		t.Fatal(err)
	}
}

func checkMoveNotifications(reg chan []byte, x uint8, y uint8) chan error {
	doneCheckingNotifier := make(chan error)
	go func() {
		defer close(doneCheckingNotifier)

		msg := <-reg
		webMsg, err := getChangeFromMsg(msg)
		if err != nil {
			doneCheckingNotifier <- fmt.Errorf("getChangeFromMsg %v", err)
			return
		}
		if webMsg.Cursor.X != x {
			doneCheckingNotifier <- fmt.Errorf("webMsg.Cursor.X should be %d but it's %d", x, webMsg.Cursor.X)
			return
		}
		if webMsg.Cursor.Y != y {
			doneCheckingNotifier <- fmt.Errorf("webMsg.Cursor.y should be %d but it's %d", y+1, webMsg.Cursor.Y)
			return
		}
	}()
	return doneCheckingNotifier
}

func checkPaintNotifications(t *testing.T, reg chan []byte, pixels ...state.Pixel) chan bool {
	doneCheckingNotifier := make(chan bool)
	go func() {
		defer close(doneCheckingNotifier)

		msg := <-reg
		webMsg, err := getChangeFromMsg(msg)
		if err != nil {
			t.Fatal("getChangeFromMsg", err)
		}

		if len(pixels) != len(webMsg.Pixels) {
			t.Fatalf("wrong length of webMsg.Pixels; should be %d but it's %d", len(pixels), len(webMsg.Pixels))
		}

		for i, p := range pixels {
			mp := webMsg.Pixels[i]
			if mp.X != p.X || mp.Y != p.Y || mp.Color != p.Color {
				t.Errorf("wrong pixel. Expected: %#v; Actual: %#v", p, mp)
			}
		}

	}()
	return doneCheckingNotifier
}

func checkNotificationsColor(t *testing.T, reg chan []byte, color *common.Color, tool string) chan bool {
	doneCheckingNotifier := make(chan bool)
	go func() {
		defer close(doneCheckingNotifier)

		msg := <-reg
		webMsg, err := getChangeFromMsg(msg)
		if err != nil {
			t.Fatal("getChangeFromMsg", err)
		}
		if color != nil {
			if webMsg.Pen == nil {
				t.Fatalf("webMsg.Cursor.Pen.Color should be #%06x but pen is nil", *color)
			} else if webMsg.Pen.Color != *color {
				t.Fatalf("webMsg.Cursor.Pen.Color should be #%06x but it's %x", *color, webMsg.Pen.Color)
			}
		}
		if tool != webMsg.ToolName {
			t.Errorf(`webMsg.Cursor.ToolName should be "%s" but it's "%s"`, tool, webMsg.ToolName)
		}
	}()
	return doneCheckingNotifier
}

func getChangeFromMsg(msg []byte) (*state.Change, error) {
	s := &state.Change{}
	if err := json.Unmarshal(msg, s); err != nil {
		return nil, err
	}
	return s, nil
}

func checkResetNotifications(reg chan []byte, canvas [][]common.Color) chan error {
	doneCheckingNotifier := make(chan error)
	go func() {
		defer close(doneCheckingNotifier)

		msg := <-reg
		webMsg, err := getChangeFromMsg(msg)
		if err != nil {
			doneCheckingNotifier <- fmt.Errorf("getChangeFromMsg %v", err)
			return
		}

		webMsgCanvas := [][]common.Color(*webMsg.Canvas)
		if !reflect.DeepEqual(webMsgCanvas, canvas) {
			doneCheckingNotifier <- fmt.Errorf("canvas should contain the point before the reset")
			return
		}
	}()
	return doneCheckingNotifier
}
