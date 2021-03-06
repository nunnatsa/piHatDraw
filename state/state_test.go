package state

import "testing"

func TestNewState(t *testing.T) {
	s := NewState()

	if len(s.Canvas) != canvasHeight {
		t.Errorf("len(s.Canvas should be %d but it's %d", canvasHeight, len(s.Canvas))
	}

	for i, line := range s.Canvas {
		if len(line) != canvasWidth {
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
	s := NewState()

	s.Cursor.X = 7
	s.Cursor.Y = 3

	s.Canvas[4][3] = true

	msg := s.CreateDisplayMessage()
	if !msg.Screen[4][3] {
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
	s := NewState()

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
	s := NewState()

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
	s := NewState()

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
	s := NewState()

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
	s := NewState()

	if s.Canvas[s.Cursor.Y][s.Cursor.X] {
		t.Errorf("s.Canvas[%d][%d] should not be set, but it is", s.Cursor.Y, s.Cursor.X)
	}

	res := s.PaintPixel()
	if !res {
		t.Error("should return true")
	}
	if !s.Canvas[s.Cursor.Y][s.Cursor.X] {
		t.Errorf("s.Canvas[%d][%d] should be set, but it's not", s.Cursor.Y, s.Cursor.X)
	}

	res = s.PaintPixel()
	if res {
		t.Error("should return false")
	}
	if !s.Canvas[s.Cursor.Y][s.Cursor.X] {
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
