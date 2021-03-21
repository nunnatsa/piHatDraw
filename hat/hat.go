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
	Stop()
}

// The format of the HAT color is 16-bit: 5 MS bits are the red color, the middle 6 bits are
// green and the 5 LB bits are blue
// rrrrrggggggbbbbb
const (
	redColor   color.Color = 0b1111100000000000
	whiteColor color.Color = 0b1111111111111111
)

const defaultJoystickFile = "/dev/input/event0"

// Joystick events
type Event uint8

const (
	Pressed Event = iota
	MoveUp
	MoveLeft
	MoveDown
	MoveRight
)

// HAT display events
type DisplayMessage struct {
	Screen  [][]common.Color
	CursorX uint8
	CursorY uint8
}

func NewDisplayMessage(mat [][]common.Color, x, y uint8) DisplayMessage {
	return DisplayMessage{
		Screen:  mat,
		CursorX: x,
		CursorY: y,
	}
}

type Hat struct {
	events chan<- Event
	screen <-chan DisplayMessage
	done   chan bool
	input  *stick.Device
}

func NewHat(joystickEvents chan<- Event, screenEvents <-chan DisplayMessage) *Hat {
	return &Hat{
		events: joystickEvents,
		screen: screenEvents,
		done:   make(chan bool),
	}
}

func (h *Hat) Start() {
	h.init()
	go h.do()
}

func (h *Hat) Stop() {
	h.done <- true
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
	defer h.gracefulShutDown()

	for {
		select {
		case event := <-h.input.Events:
			switch event.Code {
			case stick.Enter:
				h.events <- Pressed
				log.Println("Joystick Event: Pressed")

			case stick.Up:
				h.events <- MoveUp
				log.Println("Joystick Event: MoveUp")

			case stick.Down:
				h.events <- MoveDown
				log.Println("Joystick Event: MoveDown")

			case stick.Left:
				h.events <- MoveLeft
				log.Println("Joystick Event: MoveLeft")

			case stick.Right:
				h.events <- MoveRight
				log.Println("Joystick Event: MoveRight")
			}

		case screenChange := <-h.screen:
			h.drawScreen(screenChange)

		case <-h.done:
			return
		}
	}
}

func (h *Hat) drawScreen(screenChange DisplayMessage) {
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

func (h Hat) gracefulShutDown() {
	screen.Clear()
	// signal the controller we've done
	close(h.events)
}
