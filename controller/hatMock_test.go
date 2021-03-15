package controller

import (
	"github.com/nunnatsa/piHatDraw/hat"
)

type hatMock struct {
	je chan hat.Event
	se chan hat.DisplayMessage
}

func (h *hatMock) MoveUp() {
	h.je <- hat.MoveUp
}

func (h *hatMock) MoveDown() {
	h.je <- hat.MoveDown
}

func (h *hatMock) MoveRight() {
	h.je <- hat.MoveRight
}

func (h *hatMock) MoveLeft() {
	h.je <- hat.MoveLeft
}

func (h *hatMock) Start() {

}

func (h *hatMock) Press() {
	h.je <- hat.Pressed
}
