package state

import (
	"testing"

	"github.com/nunnatsa/piHatDraw/common"
)

var (
	canvasWidth  = uint8(40)
	canvasHeight = uint8(24)
)

func TestNewState(t *testing.T) {

	s := NewState(canvasWidth, canvasHeight)

	if len(s.Canvas) != int(canvasHeight) {
		t.Errorf("len(s.Canvas should be %d but it's %d", canvasHeight, len(s.Canvas))
	}

	for i, line := range s.Canvas {
		if len(line) != int(canvasWidth) {
			t.Errorf("line %d length should be %d but it's %d", i, canvasWidth, len(line))
		}
	}

	if s.Cursor.X != canvasWidth/2 {
		t.Errorf("Cursor.X should be %d but it's %d", canvasWidth/2, s.Cursor.X)
	}

	if s.Cursor.Y != canvasHeight/2 {
		t.Errorf("Cursor.Y should be %d but it's %d", canvasHeight/2, s.Cursor.Y)
	}
}

func TestCreateDisplayMessage(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	s.Cursor.X = 7
	s.Cursor.Y = 3

	s.Window.X = 0
	s.Window.Y = 0

	s.Canvas[4][3] = common.Color(0xAABBCC)

	msg := s.CreateDisplayMessage()
	if msg.Screen[4][3] != common.Color(0xAABBCC) {
		t.Error("msg.Screen[4][3] should be set")
	}

	if msg.CursorX != 7 {
		t.Errorf("msg.CursorX should be 7 but it's %d", msg.CursorX)
	}

	if msg.CursorY != 3 {
		t.Errorf("msg.CursorX should be 3 but it's %d", msg.CursorY)
	}
}

func TestStateGoUp(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.Cursor.X
	y := s.Cursor.Y

	s.GoUp()
	if s.Cursor.Y != y-1 {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y-1, s.Cursor.Y)
	}
	if s.Cursor.X != x {
		t.Errorf("s.Cursor.X should be %d but it's %d", x, s.Cursor.X)
	}

	s.Cursor.Y = 0
	s.GoUp()
	if s.Cursor.Y != 0 {
		t.Errorf("s.Cursor.Y should be 0 but it's %d", s.Cursor.Y)
	}
	if s.Cursor.X != x {
		t.Errorf("s.Cursor.X should be %d but it's %d", x, s.Cursor.X)
	}
}

func TestStateGoDown(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.Cursor.X
	y := s.Cursor.Y

	s.GoDown()
	if s.Cursor.Y != y+1 {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y+1, s.Cursor.Y)
	}
	if s.Cursor.X != x {
		t.Errorf("s.Cursor.X should be %d but it's %d", x, s.Cursor.X)
	}

	s.Cursor.Y = canvasHeight - 1
	s.GoDown()
	if s.Cursor.Y != canvasHeight-1 {
		t.Errorf("s.Cursor.Y should be %d but it's %d", canvasHeight-1, s.Cursor.Y)
	}
	if s.Cursor.X != x {
		t.Errorf("s.Cursor.X should be %d but it's %d", x, s.Cursor.X)
	}
}

func TestStateGoLeft(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.Cursor.X
	y := s.Cursor.Y

	s.GoLeft()
	if s.Cursor.X != x-1 {
		t.Errorf("s.Cursor.Y should be %d but it's %d", x-1, s.Cursor.X)
	}
	if s.Cursor.Y != y {
		t.Errorf("s.Cursor.X should be %d but it's %d", y, s.Cursor.Y)
	}

	s.Cursor.X = 0
	s.GoLeft()
	if s.Cursor.X != 0 {
		t.Errorf("s.Cursor.X should be 0 but it's %d", s.Cursor.X)
	}
	if s.Cursor.Y != y {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y, s.Cursor.Y)
	}
}

func TestStateGoRight(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.Cursor.X
	y := s.Cursor.Y

	s.GoRight()
	if s.Cursor.X != x+1 {
		t.Errorf("s.Cursor.X should be %d but it's %d", x+1, s.Cursor.X)
	}
	if s.Cursor.Y != y {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y, s.Cursor.Y)
	}

	s.Cursor.X = canvasWidth - 1
	s.GoRight()
	if s.Cursor.X != canvasWidth-1 {
		t.Errorf("s.Cursor.X should be %d but it's %d", canvasWidth-1, s.Cursor.X)
	}
	if s.Cursor.Y != y {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y, s.Cursor.Y)
	}
}

func TestStatePaintPixel(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	if s.Canvas[s.Cursor.Y][s.Cursor.X] != 0 {
		t.Errorf("s.Canvas[%d][%d] should not be set, but it is", s.Cursor.Y, s.Cursor.X)
	}

	res := s.PaintPixel()
	if !res {
		t.Error("should return true")
	}
	if s.Canvas[s.Cursor.Y][s.Cursor.X] != 0xFFFFFF {
		t.Errorf("s.Canvas[%d][%d] should be set, but it's not", s.Cursor.Y, s.Cursor.X)
	}

	res = s.PaintPixel()
	if res {
		t.Error("should return false")
	}
	if s.Canvas[s.Cursor.Y][s.Cursor.X] != 0xFFFFFF {
		t.Errorf("s.Canvas[%d][%d] should be set, but it's not", s.Cursor.Y, s.Cursor.X)
	}

	s.Cursor.X = canvasWidth
	res = s.PaintPixel()
	if res {
		t.Error("should return false")
	}

	s.Cursor.X = canvasWidth / 2

	s.Cursor.Y = canvasHeight
	res = s.PaintPixel()
	if res {
		t.Error("should return false")
	}
}

func TestState_getColor(t *testing.T) {
	cr := cursor{
		X: 4,
		Y: 5,
	}

	s := State{
		Cursor: cr,
		Pen:    &pen{Color: 0x12345},
	}
	s.SetPen()

	if c := s.tool.GetColor(); c != 0x12345 {
		t.Errorf("color should be 0x12345 but it's 0x%x", c)
	}

	s.Pen.Color = 0x12346

	if c := s.tool.GetColor(); c != 0x12346 {
		t.Errorf("color should be 0x12346 but it's 0x%x", c)
	}

	s.SetEraser()

	if c := s.Pen.GetColor(); c != 0x12346 {
		t.Errorf("color should be 0x12346 but it's 0x%x", c)
	}

	if c := s.tool.GetColor(); c != 0 {
		t.Errorf("color should be 0 but it's 0x%x", c)
	}

	s.SetPen()

	if c := s.Pen.GetColor(); c != 0x12346 {
		t.Errorf("color should be 0x12346 but it's 0x%x", c)
	}

	if c := s.tool.GetColor(); c != 0x12346 {
		t.Errorf("color should be 0x12346 but it's 0x%x", c)
	}
}
