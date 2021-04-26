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

	if len(s.canvas) != int(canvasHeight) {
		t.Errorf("len(s.Canvas should be %d but it's %d", canvasHeight, len(s.canvas))
	}

	for i, line := range s.canvas {
		if len(line) != int(canvasWidth) {
			t.Errorf("line %d length should be %d but it's %d", i, canvasWidth, len(line))
		}
	}

	if s.cursor.X != canvasWidth/2 {
		t.Errorf("Cursor.X should be %d but it's %d", canvasWidth/2, s.cursor.X)
	}

	if s.cursor.Y != canvasHeight/2 {
		t.Errorf("Cursor.Y should be %d but it's %d", canvasHeight/2, s.cursor.Y)
	}
}

func TestCreateDisplayMessage(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	s.cursor.X = 7
	s.cursor.Y = 3

	s.window.X = 0
	s.window.Y = 0

	s.canvas[4][3] = common.Color(0xAABBCC)

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

	x := s.cursor.X
	y := s.cursor.Y

	change := s.GoUp()
	if s.cursor.Y != y-1 {
		t.Errorf("s.cursor.Y should be %d but it's %d", y-1, s.cursor.Y)
	}
	if s.cursor.X != x {
		t.Errorf("s.cursor.X should be %d but it's %d", x, s.cursor.X)
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

	s.cursor.Y = 0
	change = s.GoUp()
	if change != nil {
		t.Errorf("change should be nil")
	}
	if s.cursor.Y != 0 {
		t.Errorf("s.cursor.Y should be %d but it's %d", 0, s.cursor.Y)
	}
	if s.cursor.X != x {
		t.Errorf("s.cursor.X should be %d but it's %d", x, s.cursor.X)
	}
}

func TestStateGoDown(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.cursor.X
	y := s.cursor.Y

	change := s.GoDown()
	if s.cursor.Y != y+1 {
		t.Errorf("s.cursor.Y should be %d but it's %d", y+1, s.cursor.Y)
	}
	if s.cursor.X != x {
		t.Errorf("s.cursor.X should be %d but it's %d", x, s.cursor.X)
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

	s.cursor.Y = canvasHeight - 1
	change = s.GoDown()
	if change != nil {
		t.Errorf("change should be nil")
	}
	if s.cursor.Y != canvasHeight-1 {
		t.Errorf("s.cursor.Y should be %d but it's %d", canvasHeight-1, s.cursor.Y)
	}
	if s.cursor.X != x {
		t.Errorf("s.cursor.X should be %d but it's %d", x, s.cursor.X)
	}
}

func TestStateGoLeft(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.cursor.X
	y := s.cursor.Y

	change := s.GoLeft()
	if s.cursor.X != x-1 {
		t.Errorf("s.cursor.X should be %d but it's %d", x-1, s.cursor.X)
	}
	if s.cursor.Y != y {
		t.Errorf("s.cursor.Y should be %d but it's %d", y, s.cursor.Y)
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

	s.cursor.X = 0
	change = s.GoLeft()
	if change != nil {
		t.Errorf("change should be nil")
	}
	if s.cursor.X != 0 {
		t.Errorf("s.cursor.X should be %d but it's %d", 0, s.cursor.X)
	}
	if s.cursor.Y != y {
		t.Errorf("s.cursor.Y should be %d but it's %d", y, s.cursor.Y)
	}
}

func TestStateGoRight(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	x := s.cursor.X
	y := s.cursor.Y

	change := s.GoRight()
	if s.cursor.X != x+1 {
		t.Errorf("s.cursor.X should be %d but it's %d", x+1, s.cursor.X)
	}
	if s.cursor.Y != y {
		t.Errorf("s.cursor.Y should be %d but it's %d", y, s.cursor.Y)
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

	s.cursor.X = canvasWidth - 1
	change = s.GoRight()
	if change != nil {
		t.Errorf("change should be nil")
	}
	if s.cursor.X != canvasWidth-1 {
		t.Errorf("s.cursor.X should be %d but it's %d", canvasWidth-1, s.cursor.X)
	}
	if s.cursor.Y != y {
		t.Errorf("s.cursor.Y should be %d but it's %d", y, s.cursor.Y)
	}
}

func TestStatePaintPixel(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)

	if s.canvas[s.cursor.Y][s.cursor.X] != 0 {
		t.Errorf("s.Canvas[%d][%d] should not be set, but it is", s.cursor.Y, s.cursor.X)
	}

	res := s.Paint()
	if res == nil {
		t.Fatal("should return a change")
	}
	if s.canvas[s.cursor.Y][s.cursor.X] != 0xFFFFFF {
		t.Fatalf("s.Canvas[%d][%d] should be set, but it's not", s.cursor.Y, s.cursor.X)
	}

	if len(res.Pixels) != 1 {
		t.Errorf("res.Pixels should be with len of 1, but it is %d.", len(res.Pixels))
	}
	pixel := res.Pixels[0]
	if pixel.X != s.cursor.X || pixel.Y != s.cursor.Y || pixel.Color != s.color {
		t.Errorf("x should be %d, y should be %d and color should be #%06x; but pixel is %#v", s.cursor.X, s.cursor.Y, s.color, pixel)
	}
	if l := undoList.len(); l != 1 {
		t.Errorf("undo list should be with len of 1, but it's with length of %d", l)
	}

	change := undoList.pop()
	if l := len(change.Pixels); l != 1 {
		t.Fatalf("pixel list should be with len of 1, but it's with length of %d", l)
	}
	if change.Pixels[0].X != res.Pixels[0].X || change.Pixels[0].Y != res.Pixels[0].Y || change.Pixels[0].Color != 0 {
		t.Errorf("x should be %d, y should be %d and color should be #%06x; but pixel is %#v", s.cursor.X, s.cursor.Y, s.color, change)
	}

	res = s.Paint()
	if res != nil {
		t.Fatal("should not return a change")
	}
	if s.canvas[s.cursor.Y][s.cursor.X] != 0xFFFFFF {
		t.Errorf("s.Canvas[%d][%d] should be set, but it's not", s.cursor.Y, s.cursor.X)
	}

	s.cursor.X = canvasWidth
	res = s.Paint()
	if res != nil {
		t.Fatal("should not return a change")
	}
	if undoList.len() != 0 {
		t.Error("undo list should be empty")
	}

	s.cursor.X = canvasWidth / 2

	s.cursor.Y = canvasHeight
	res = s.Paint()
	if res != nil {
		t.Fatal("should not return a change")
	}
}

func TestState_setColor(t *testing.T) {
	s := NewState(8, 8)
	s.color = 0x123456

	change := s.SetColor(0x123456)
	if change != nil {
		t.Errorf("change should be nil, but it's %#v", change)
	}

	change = s.SetColor(0x654321)
	if c := *change.Color; c != 0x654321 {
		t.Errorf("color should be 0x12346 but it's 0x%x", c)
	}
}

func TestState_setTool(t *testing.T) {
	s := NewState(8, 8)

	change, err := s.SetTool(penName)
	if err != nil {
		t.Fatalf("error should be nil but it's %#v", err)
	}

	if change != nil {
		t.Fatalf("change should be nil, but it's %#v", change)
	}

	if s.toolName != penName {
		t.Errorf("color should be %s but it's 0x%x", penName, s.toolName)
	}

	change, err = s.SetTool(eraserName)
	if err != nil {
		t.Fatalf("error should be nil but it's %#v", err)
	}

	if change == nil {
		t.Fatal("change should be not nil")
	}

	if tn := change.ToolName; tn != eraserName {
		t.Errorf("color should be %s but it's 0x%x", eraserName, tn)
	}

	if s.toolName != eraserName {
		t.Errorf("color should be %s but it's 0x%x", eraserName, s.toolName)
	}

	change, err = s.SetTool(penName)
	if err != nil {
		t.Fatalf("error should be nil but it's %#v", err)
	}

	if change == nil {
		t.Fatal("change should be not nil")
	}

	if tn := change.ToolName; tn != penName {
		t.Errorf("color should be %s but it's 0x%x", penName, tn)
	}

	if s.toolName != penName {
		t.Errorf("color should be %s but it's 0x%x", penName, s.toolName)
	}
}

func TestState_Undo(t *testing.T) {
	s := NewState(canvasWidth, canvasHeight)
	if s.canvas[3][4] != 0 {
		t.Error("s.Canvas[3][4] should be 0")
	}
	if s.canvas[10][20] != 0 {
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
	if s.canvas[3][4] != 0x112233 {
		t.Error("s.Canvas[3][4] should be 0x112233")
	}
	if s.canvas[10][20] != 0x112233 {
		t.Error("s.Canvas[3][4] should be 0x112233")
	}

}
