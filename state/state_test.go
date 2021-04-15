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

	change := s.GoUp()
	if s.Cursor.Y != y-1 {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y-1, s.Cursor.Y)
	}
	if s.Cursor.X != x {
		t.Errorf("s.Cursor.X should be %d but it's %d", x, s.Cursor.X)
	}
	if change.Cursor.Y != y-1 {
		t.Errorf("change.Cursor.Y should be %d but it's %d", y-1, change.Cursor.Y)
	}

	if change.Cursor.X != x {
		t.Errorf("change.Cursor.X should be %d but it's %d", x, change.Cursor.X)
	}

	if l := undoList.len(); l > 0 {
		t.Errorf("undo list should be empty, but it's with length of %d", l)
	}

	s.Cursor.Y = 0
	change = s.GoUp()
	if change != nil {
		t.Errorf("change should be nil")
	}
	if s.Cursor.Y != 0 {
		t.Errorf("s.Cursor.Y should be %d but it's %d", 0, s.Cursor.Y)
	}
	if s.Cursor.X != x {
		t.Errorf("s.Cursor.X should be %d but it's %d", x, s.Cursor.X)
	}
}

func TestStateGoDown(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.Cursor.X
	y := s.Cursor.Y

	change := s.GoDown()
	if s.Cursor.Y != y+1 {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y+1, s.Cursor.Y)
	}
	if s.Cursor.X != x {
		t.Errorf("s.Cursor.X should be %d but it's %d", x, s.Cursor.X)
	}
	if change.Cursor.Y != y+1 {
		t.Errorf("change.Cursor.Y should be %d but it's %d", y+1, change.Cursor.Y)
	}
	if change.Cursor.X != x {
		t.Errorf("change.Cursor.X should be %d but it's %d", x, change.Cursor.X)
	}

	if l := undoList.len(); l > 0 {
		t.Errorf("undo list should be empty, but it's with length of %d", l)
	}

	s.Cursor.Y = canvasHeight - 1
	change = s.GoDown()
	if change != nil {
		t.Errorf("change should be nil")
	}
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

	change := s.GoLeft()
	if s.Cursor.X != x-1 {
		t.Errorf("s.Cursor.X should be %d but it's %d", x-1, s.Cursor.X)
	}
	if s.Cursor.Y != y {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y, s.Cursor.Y)
	}
	if change.Cursor.X != x-1 {
		t.Errorf("change.Cursor.X should be %d but it's %d", x-1, change.Cursor.X)
	}

	if change.Cursor.Y != y {
		t.Errorf("change.Cursor.Y should be %d but it's %d", y, change.Cursor.Y)
	}

	if l := undoList.len(); l > 0 {
		t.Errorf("undo list should be empty, but it's with length of %d", l)
	}

	s.Cursor.X = 0
	change = s.GoLeft()
	if change != nil {
		t.Errorf("change should be nil")
	}
	if s.Cursor.X != 0 {
		t.Errorf("s.Cursor.X should be %d but it's %d", 0, s.Cursor.X)
	}
	if s.Cursor.Y != y {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y, s.Cursor.Y)
	}
}

func TestStateGoRight(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.Cursor.X
	y := s.Cursor.Y

	change := s.GoRight()
	if s.Cursor.X != x+1 {
		t.Errorf("s.Cursor.X should be %d but it's %d", x+1, s.Cursor.X)
	}
	if s.Cursor.Y != y {
		t.Errorf("s.Cursor.Y should be %d but it's %d", y, s.Cursor.Y)
	}
	if change.Cursor.X != x+1 {
		t.Errorf("change.Cursor.X should be %d but it's %d", x+1, change.Cursor.X)
	}

	if change.Cursor.Y != y {
		t.Errorf("change.Cursor.Y should be %d but it's %d", y, change.Cursor.Y)
	}

	if l := undoList.len(); l > 0 {
		t.Errorf("undo list should be empty, but it's with length of %d", l)
	}

	s.Cursor.X = canvasWidth - 1
	change = s.GoRight()
	if change != nil {
		t.Errorf("change should be nil")
	}
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
	if res == nil {
		t.Fatal("should return a change")
	}
	if s.Canvas[s.Cursor.Y][s.Cursor.X] != 0xFFFFFF {
		t.Fatalf("s.Canvas[%d][%d] should be set, but it's not", s.Cursor.Y, s.Cursor.X)
	}

	if len(res.Pixels) != 1 {
		t.Errorf("res.Pixels should be with len of 1, but it is %d.", len(res.Pixels))
	}
	pixel := res.Pixels[0]
	if pixel.X != s.Cursor.X || pixel.Y != s.Cursor.Y || pixel.Color != s.Pen.Color {
		t.Errorf("x should be %d, y should be %d and color should be #%06x; but pixel is %#v", s.Cursor.X, s.Cursor.Y, s.Pen.Color, pixel)
	}
	if l := undoList.len(); l != 1 {
		t.Errorf("undo list should be with len of 1, but it's with length of %d", l)
	}

	change := undoList.pop()
	if l := len(change.Pixels); l != 1 {
		t.Fatalf("pixel list should be with len of 1, but it's with length of %d", l)
	}
	if change.Pixels[0].X != res.Pixels[0].X || change.Pixels[0].Y != res.Pixels[0].Y || change.Pixels[0].Color != 0 {
		t.Errorf("x should be %d, y should be %d and color should be #%06x; but pixel is %#v", s.Cursor.X, s.Cursor.Y, s.Pen.Color, change)
	}

	res = s.PaintPixel()
	if res != nil {
		t.Fatal("should not return a change")
	}
	if s.Canvas[s.Cursor.Y][s.Cursor.X] != 0xFFFFFF {
		t.Errorf("s.Canvas[%d][%d] should be set, but it's not", s.Cursor.Y, s.Cursor.X)
	}

	s.Cursor.X = canvasWidth
	res = s.PaintPixel()
	if res != nil {
		t.Fatal("should not return a change")
	}
	if undoList.len() != 0 {
		t.Error("undo list should be empty")
	}

	s.Cursor.X = canvasWidth / 2

	s.Cursor.Y = canvasHeight
	res = s.PaintPixel()
	if res != nil {
		t.Fatal("should not return a change")
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

func TestState_Undo(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)
	if s.Canvas[3][4] != 0 {
		t.Error("s.Canvas[3][4] should be 0")
	}
	if s.Canvas[10][20] != 0 {
		t.Error("s.Canvas[3][4] should be 0")
	}

	c := &Change{
		Pixels: []Pixel{{
			X: 4, Y: 3, Color: 0x112233,
		}, {
			X: 20, Y: 10, Color: 0x112233,
		}},
	}

	undoList.push(c)

	s.Undo()
	if s.Canvas[3][4] != 0x112233 {
		t.Error("s.Canvas[3][4] should be 0x112233")
	}
	if s.Canvas[10][20] != 0x112233 {
		t.Error("s.Canvas[3][4] should be 0x112233")
	}

}
