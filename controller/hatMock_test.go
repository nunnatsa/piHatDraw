package controller

import "github.com/nunnatsa/piHatDraw/common"

type hatMock struct {
	je chan common.HatEvent
	se chan *common.DisplayMessage
}

func (h *hatMock) MoveUp() {
	h.je <- common.MoveUp
}

func (h *hatMock) MoveDown() {
	h.je <- common.MoveDown
}

func (h *hatMock) MoveRight() {
	h.je <- common.MoveRight
}

func (h *hatMock) MoveLeft() {
	h.je <- common.MoveLeft
}

func (h *hatMock) Start() {

}

func (h *hatMock) Press() {
	h.je <- common.Pressed
}
