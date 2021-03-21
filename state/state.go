package state

import (
	"log"

	"github.com/nunnatsa/piHatDraw/hat"

	"github.com/nunnatsa/piHatDraw/common"
)

const (
	canvasHeight = common.WindowSize
	canvasWidth  = common.WindowSize
)

type canvas [][]common.Color

type cursor struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

type State struct {
	Canvas canvas `json:"canvas,omitempty"`
	Cursor cursor `json:"cursor,omitempty"`
}

func NewState() *State {
	c := make([][]common.Color, canvasHeight)
	for y := 0; y < canvasHeight; y++ {
		c[y] = make([]common.Color, canvasWidth)
	}

	return &State{
		Canvas: c,
		Cursor: cursor{X: canvasWidth / 2, Y: canvasHeight / 2},
	}
}

func (s *State) GoUp() bool {
	if s.Cursor.Y > 0 {
		s.Cursor.Y--
		return true
	}
	return false
}

func (s *State) GoLeft() bool {
	if s.Cursor.X > 0 {
		s.Cursor.X--
		return true
	}
	return false
}

func (s *State) GoDown() bool {
	if s.Cursor.Y < canvasHeight-1 {
		s.Cursor.Y++
		return true
	}

	return false
}

func (s *State) GoRight() bool {
	if s.Cursor.X < canvasWidth-1 {
		s.Cursor.X++
		return true
	}

	return false
}

func (s *State) PaintPixel() bool {
	if s.Cursor.Y >= canvasHeight || s.Cursor.X >= canvasWidth {
		log.Printf("Error: Cursor (%d, %d) is out of canvas\n", s.Cursor.X, s.Cursor.Y)
		return false
	}
	if !s.Canvas[s.Cursor.Y][s.Cursor.X] {
		s.Canvas[s.Cursor.Y][s.Cursor.X] = true
		return true
	}

	return false
}

func (s State) CreateDisplayMessage() hat.DisplayMessage {
	c := make([][]common.Color, common.WindowSize)
	for y := 0; y < common.WindowSize; y++ {
		c[y] = make([]common.Color, 0, common.WindowSize)
		c[y] = append(c[y], s.Canvas[y]...)
	}

	return hat.NewDisplayMessage(c, s.Cursor.X, s.Cursor.Y)
}
