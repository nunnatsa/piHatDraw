package state

import (
	"fmt"
	"log"

	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/hat"
)

const (
	penName    = "pen"
	eraserName = "eraser"
)

type Canvas [][]common.Color

func (c Canvas) Clone() Canvas {
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

type tool func() *Change

type State struct {
	canvas       Canvas
	cursor       cursor
	window       window
	canvasWidth  uint8
	canvasHeight uint8
	toolName     string
	color        common.Color
	tools        map[string]tool
}

func NewState(canvasWidth, canvasHeight uint8) *State {
	s := &State{
		canvasWidth:  canvasWidth,
		canvasHeight: canvasHeight,
	}

	s.tools = map[string]tool{
		penName:    s.pen,
		eraserName: s.eraser,
	}

	_ = s.Reset()

	return s
}

func (s State) GetCanvasClone() Canvas {
	return s.canvas.Clone()
}

func (s *State) Reset() *Change {
	if len(s.canvas) > 0 {
		chng := &Change{
			Canvas: s.canvas.Clone(),
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

	s.canvas = c
	s.cursor = cr
	s.window = win
	s.color = 0xFFFFFF
	s.toolName = penName

	return s.GetFullChange()
}

func (s *State) GoUp() *Change {
	if s.cursor.Y > 0 {
		s.cursor.Y--
		if s.cursor.Y < s.window.Y {
			s.window.Y = s.cursor.Y
		}
		return s.getPositionChange()
	}
	return nil
}

func (s *State) GoLeft() *Change {
	if s.cursor.X > 0 {
		s.cursor.X--
		if s.cursor.X < s.window.X {
			s.window.X = s.cursor.X
		}
		return s.getPositionChange()
	}
	return nil
}

func (s *State) GoDown() *Change {
	if s.cursor.Y < s.canvasHeight-1 {
		s.cursor.Y++
		if s.cursor.Y > s.window.Y+common.WindowSize-1 {
			s.window.Y++
		}
		return s.getPositionChange()
	}

	return nil
}

func (s *State) GoRight() *Change {
	if s.cursor.X < s.canvasWidth-1 {
		s.cursor.X++
		if s.cursor.X > s.window.X+common.WindowSize-1 {
			s.window.X++
		}
		return s.getPositionChange()
	}

	return nil
}

func (s *State) Paint() *Change {
	return s.tools[s.toolName]()
}

func (s *State) pen() *Change {
	px := s.paintPixel(s.color)
	if px == nil {
		return nil
	}

	return &Change{
		Pixels: []Pixel{*px},
	}
}

func (s *State) eraser() *Change {
	px := s.paintPixel(0)
	if px == nil {
		return nil
	}

	return &Change{
		Pixels: []Pixel{*px},
	}
}

func (s *State) paintPixel(color common.Color) *Pixel {
	if s.cursor.Y >= s.canvasHeight || s.cursor.X >= s.canvasWidth {
		log.Printf("Error: Cursor (%d, %d) is out of canvas\n", s.cursor.X, s.cursor.Y)
		return nil
	}

	if s.canvas[s.cursor.Y][s.cursor.X] != color {
		chng := &Change{
			Pixels: []Pixel{{X: s.cursor.X, Y: s.cursor.Y, Color: s.canvas[s.cursor.Y][s.cursor.X]}},
		}
		undoList.push(chng)
		s.canvas[s.cursor.Y][s.cursor.X] = color
		return &Pixel{
			X:     s.cursor.X,
			Y:     s.cursor.Y,
			Color: color,
		}
	}

	return nil
}

func (s State) CreateDisplayMessage() hat.DisplayMessage {
	c := make([][]common.Color, common.WindowSize)
	for y := uint8(0); y < common.WindowSize; y++ {
		c[y] = make([]common.Color, 0, common.WindowSize)
		c[y] = append(c[y], s.canvas[s.window.Y+y][s.window.X:s.window.X+common.WindowSize]...)
	}

	return hat.NewDisplayMessage(c, s.cursor.X-s.window.X, s.cursor.Y-s.window.Y)
}

func (s *State) SetColor(cl common.Color) *Change {
	if s.color != cl {
		s.color = cl
		return &Change{
			Color: &cl,
		}
	}
	return nil
}

func (s *State) SetTool(toolName string) (*Change, error) {
	if _, found := s.tools[toolName]; !found {
		return nil, fmt.Errorf(`unknown tool "%s"`, toolName)
	}

	if toolName == s.toolName {
		return nil, nil
	}

	s.toolName = toolName
	return &Change{
		ToolName: toolName,
	}, nil
}

func (s State) getPositionChange() *Change {
	return &Change{
		Cursor: &cursor{
			X: s.cursor.X,
			Y: s.cursor.Y,
		},
		Window: &window{
			X: s.window.X,
			Y: s.window.Y,
		},
	}
}

func (s State) GetFullChange() *Change {
	return &Change{
		Canvas:   s.canvas.Clone(),
		Cursor:   &s.cursor,
		Window:   &s.window,
		ToolName: s.toolName,
		Color:    &s.color,
	}
}

func (s *State) Undo() *Change {
	chng := undoList.pop()
	if chng != nil {
		if chng.Canvas != nil {
			s.canvas = chng.Canvas
		} else if len(chng.Pixels) > 0 {
			for _, pixel := range chng.Pixels {
				s.canvas[pixel.Y][pixel.X] = pixel.Color
			}
		}
	}

	return chng
}
