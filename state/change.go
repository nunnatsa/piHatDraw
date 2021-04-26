package state

import (
	"github.com/nunnatsa/piHatDraw/common"
)

type Pixel struct {
	X     uint8        `json:"x"`
	Y     uint8        `json:"y"`
	Color common.Color `json:"color"`
}

type Change struct {
	Canvas   Canvas        `json:"canvas,omitempty"`
	Cursor   *cursor       `json:"cursor,omitempty"`
	Window   *window       `json:"window,omitempty"`
	ToolName string        `json:"toolName,omitempty"`
	Color    *common.Color `json:"color,omitempty"`

	Pixels []Pixel `json:"pixels,omitempty"`
}

type changeNode struct {
	data *Change
	next *changeNode
}

type changeStack struct {
	head *changeNode
}

func (s *changeStack) push(change *Change) {
	s.head = &changeNode{
		data: change,
		next: s.head,
	}
}

func (s *changeStack) pop() *Change {
	if s == nil || s.head == nil {
		return nil
	}
	res := s.head
	s.head = s.head.next

	return res.data
}

var undoList = &changeStack{
	head: nil,
}
