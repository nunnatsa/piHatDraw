package state

import (
	"log"

	"github.com/nunnatsa/piHatDraw/hat"

	"github.com/nunnatsa/piHatDraw/common"
)

type canvas [][]common.Color

type cursor struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

type window struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

type tool interface {
	GetColor() common.Color
}

type pen struct {
	Color common.Color `json:"color"`
}

func (p pen) GetColor() common.Color {
	return p.Color
}

func (p *pen) SetColor(c common.Color) {
	p.Color = c
}

type Eraser struct{}

func (Eraser) GetColor() common.Color {
	return 0
}

type State struct {
	Canvas       canvas `json:"canvas,omitempty"`
	Cursor       cursor `json:"cursor,omitempty"`
	Window       window `json:"window,omitempty"`
	canvasWidth  uint8
	canvasHeight uint8
	ToolName     string `json:"tool"`
	Pen          *pen   `json:"pen"`
	tool         tool
}

func NewState(canvasWidth, canvasHeight uint8) *State {
	s := &State{
		canvasWidth:  canvasWidth,
		canvasHeight: canvasHeight,
	}

	s.Reset()

	return s
}

func (s *State) Reset() {
	c := make([][]common.Color, s.canvasHeight)
	for y := uint8(0); y < s.canvasHeight; y++ {
		c[y] = make([]common.Color, s.canvasWidth)
	}

	cr := cursor{X: s.canvasWidth / 2, Y: s.canvasHeight / 2}
	halfWindow := uint8(common.WindowSize / 2)
	win := window{X: cr.X - halfWindow, Y: cr.Y - halfWindow}

	s.Canvas = c
	s.Cursor = cr
	s.Window = win
	s.Pen = &pen{Color: 0xFFFFFF}
	s.SetPen()
}

func (s *State) GoUp() bool {
	if s.Cursor.Y > 0 {
		s.Cursor.Y--
		if s.Cursor.Y < s.Window.Y {
			s.Window.Y = s.Cursor.Y
		}
		return true
	}
	return false
}

func (s *State) GoLeft() bool {
	if s.Cursor.X > 0 {
		s.Cursor.X--
		if s.Cursor.X < s.Window.X {
			s.Window.X = s.Cursor.X
		}
		return true
	}
	return false
}

func (s *State) GoDown() bool {
	if s.Cursor.Y < s.canvasHeight-1 {
		s.Cursor.Y++
		if s.Cursor.Y > s.Window.Y+common.WindowSize-1 {
			s.Window.Y++
		}
		return true
	}

	return false
}

func (s *State) GoRight() bool {
	if s.Cursor.X < s.canvasWidth-1 {
		s.Cursor.X++
		if s.Cursor.X > s.Window.X+common.WindowSize-1 {
			s.Window.X++
		}
		return true
	}

	return false
}

func (s *State) PaintPixel() bool {
	if s.Cursor.Y >= s.canvasHeight || s.Cursor.X >= s.canvasWidth {
		log.Printf("Error: Cursor (%d, %d) is out of canvas\n", s.Cursor.X, s.Cursor.Y)
		return false
	}

	c := s.tool.GetColor()
	if s.Canvas[s.Cursor.Y][s.Cursor.X] != c {
		s.Canvas[s.Cursor.Y][s.Cursor.X] = c
		return true
	}

	return false
}

func (s State) CreateDisplayMessage() hat.DisplayMessage {
	c := make([][]common.Color, common.WindowSize)
	for y := uint8(0); y < common.WindowSize; y++ {
		c[y] = make([]common.Color, 0, common.WindowSize)
		c[y] = append(c[y], s.Canvas[s.Window.Y+y][s.Window.X:s.Window.X+common.WindowSize]...)
	}

	return hat.NewDisplayMessage(c, s.Cursor.X-s.Window.X, s.Cursor.Y-s.Window.Y)
}

func (s *State) SetColor(cl common.Color) bool {
	if s.Pen.GetColor() != cl {
		s.Pen.SetColor(cl)
		return true
	}
	return false
}

func (s *State) SetPen() bool {
	s.tool = s.Pen
	if s.ToolName != "pen" {
		s.ToolName = "pen"
		return true
	}
	return false
}

var eraser Eraser

func (s *State) SetEraser() bool {
	if s.ToolName != "eraser" {
		s.ToolName = "eraser"
		s.tool = eraser
		return true
	}
	return false
}
