package hat

import (
	"log"

	"github.com/nathany/bobblehat/sense/screen"
	"github.com/nathany/bobblehat/sense/screen/color"
	"github.com/nathany/bobblehat/sense/stick"
	"github.com/nunnatsa/piHatDraw/common"
)

type Interface interface {
	Start()
}

// The format of the HAT color is 16-bit: 5 MS bits are the red color, the middle 6 bits are
// green and the 5 LB bits are blue
// rrrrrggggggbbbbb
const (
	redColor   color.Color = 0b1111100000000000
	whiteColor color.Color = 0b1111111111111111
)

const defaultJoystickFile = "/dev/input/event0"

type Hat struct {
	events chan<- common.HatEvent
	screen <-chan *common.DisplayMessage
	input  *stick.Device
}

func NewHat(joystickEvents chan<- common.HatEvent, screenEvents <-chan *common.DisplayMessage) *Hat {
	return &Hat{events: joystickEvents, screen: screenEvents}
}

func (h *Hat) Start() {
	h.init()
	go h.do()
}

func (h *Hat) init() {
	var err error
	h.input, err = stick.Open(defaultJoystickFile)
	if err != nil {
		log.Panic("can't open "+defaultJoystickFile, err)
	}

	if err = screen.Clear(); err != nil {
		log.Panic("can't clear the HAT display", err)
	}
}

func (h *Hat) do() {
	for {
		select {
		case event := <-h.input.Events:
			switch event.Code {
			case stick.Enter:
				h.events <- common.Pressed
				log.Println("Joystick Event: Pressed")

			case stick.Up:
				h.events <- common.MoveUp
				log.Println("Joystick Event: MoveUp")

			case stick.Down:
				h.events <- common.MoveDown
				log.Println("Joystick Event: MoveDown")

			case stick.Left:
				h.events <- common.MoveLeft
				log.Println("Joystick Event: MoveLeft")

			case stick.Right:
				h.events <- common.MoveRight
				log.Println("Joystick Event: MoveRight")
			}

		case screenChange := <-h.screen:
			h.drawScreen(screenChange)
		}
	}
}

func (h *Hat) drawScreen(screenChange *common.DisplayMessage) {
	fb := screen.NewFrameBuffer()
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if screenChange.Screen[y][x] {
				fb.SetPixel(x, y, whiteColor)
			}
		}
	}
	fb.SetPixel(int(screenChange.CursorX), int(screenChange.CursorY), redColor)
	err := screen.Draw(fb)
	if err != nil {
		log.Println("error while printing to HAT display:", err)
	}
}
