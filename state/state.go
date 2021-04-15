package state

import (
	"log"

	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/hat"
)

type canvas [][]common.Color

func (c canvas) Clone() canvas {
	if len(c) == 0 || len(c[0]) == 0 {
		return nil
	}

	newCanvas := make([][]common.Color, len(c))
	for y, line := range c {
		newCanvas[y] = make([]common.Color, len(line))
		copy(newCanvas[y], line)
	}

	return newCanvas
}

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
	ToolName     string `json:"toolName"`
	Pen          *pen   `json:"pen"`
	tool         tool
}

func NewState(canvasWidth, canvasHeight uint8) *State {
	s := &State{
		canvasWidth:  canvasWidth,
		canvasHeight: canvasHeight,
	}

	_ = s.Reset()

	return s
}

func (s *State) Reset() *Change {
	if len(s.Canvas) > 0 {
		cv := s.Canvas.Clone()
		chng := &Change{
			Canvas: &cv,
		}

		undoList.push(chng)
	}

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

	return s.GetFullChange()
}

func (s *State) GoUp() *Change {
	if s.Cursor.Y > 0 {
		s.Cursor.Y--
		if s.Cursor.Y < s.Window.Y {
			s.Window.Y = s.Cursor.Y
		}
		return s.getPositionChange()
	}
	return nil
}

func (s *State) GoLeft() *Change {
	if s.Cursor.X > 0 {
		s.Cursor.X--
		if s.Cursor.X < s.Window.X {
			s.Window.X = s.Cursor.X
		}
		return s.getPositionChange()
	}
	return nil
}

func (s *State) GoDown() *Change {
	if s.Cursor.Y < s.canvasHeight-1 {
		s.Cursor.Y++
		if s.Cursor.Y > s.Window.Y+common.WindowSize-1 {
			s.Window.Y++
		}
		return s.getPositionChange()
	}

	return nil
}

func (s *State) GoRight() *Change {
	if s.Cursor.X < s.canvasWidth-1 {
		s.Cursor.X++
		if s.Cursor.X > s.Window.X+common.WindowSize-1 {
			s.Window.X++
		}
		return s.getPositionChange()
	}

	return nil
}

func (s *State) PaintPixel() *Change {
	if s.Cursor.Y >= s.canvasHeight || s.Cursor.X >= s.canvasWidth {
		log.Printf("Error: Cursor (%d, %d) is out of canvas\n", s.Cursor.X, s.Cursor.Y)
		return nil
	}

	c := s.tool.GetColor()
	if s.Canvas[s.Cursor.Y][s.Cursor.X] != c {
		chng := &Change{
			Pixels: []Pixel{{X: s.Cursor.X, Y: s.Cursor.Y, Color: s.Canvas[s.Cursor.Y][s.Cursor.X]}},
		}
		undoList.push(chng)
		s.Canvas[s.Cursor.Y][s.Cursor.X] = c
		return &Change{
			Pixels: []Pixel{{
				X:     s.Cursor.X,
				Y:     s.Cursor.Y,
				Color: c,
			}},
		}
	}

	return nil
}

func (s State) CreateDisplayMessage() hat.DisplayMessage {
	c := make([][]common.Color, common.WindowSize)
	for y := uint8(0); y < common.WindowSize; y++ {
		c[y] = make([]common.Color, 0, common.WindowSize)
		c[y] = append(c[y], s.Canvas[s.Window.Y+y][s.Window.X:s.Window.X+common.WindowSize]...)
	}

	return hat.NewDisplayMessage(c, s.Cursor.X-s.Window.X, s.Cursor.Y-s.Window.Y)
}

func (s *State) SetColor(cl common.Color) *Change {
	if s.Pen.GetColor() != cl {
		s.Pen.SetColor(cl)
		return &Change{
			Pen: &pen{
				Color: cl,
			},
		}
	}
	return nil
}

const (
	penName    = "pen"
	eraserName = "eraser"
)

func (s *State) SetPen() *Change {
	s.tool = s.Pen
	if s.ToolName != penName {
		s.ToolName = penName
		return &Change{
			ToolName: penName,
		}
	}
	return nil
}

var eraser Eraser

func (s *State) SetEraser() *Change {
	if s.ToolName != eraserName {
		s.ToolName = eraserName
		s.tool = eraser
		return &Change{
			ToolName: eraserName,
		}
	}
	return nil
}

func (s State) getPositionChange() *Change {
	return &Change{
		Cursor: &cursor{
			X: s.Cursor.X,
			Y: s.Cursor.Y,
		},
		Window: &window{
			X: s.Window.X,
			Y: s.Window.Y,
		},
	}
}

func (s State) GetFullChange() *Change {
	cv := s.Canvas.Clone()
	return &Change{
		Canvas:   &cv,
		Cursor:   &s.Cursor,
		Window:   &s.Window,
		ToolName: s.ToolName,
		Pen:      s.Pen,
	}
}

func (s *State) Undo() *Change {
	chng := undoList.pop()
	if chng != nil {
		if chng.Canvas != nil {
			s.Canvas = *chng.Canvas
		} else if len(chng.Pixels) > 0 {
			for _, pixel := range chng.Pixels {
				s.Canvas[pixel.Y][pixel.X] = pixel.Color
			}
		}
	}

	return chng
}
