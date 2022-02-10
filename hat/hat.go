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
	redColor color.Color  = 0b1111100000000000
	rmask    common.Color = 0b111110000000000000000000
	gmask    common.Color = 0b000000001111110000000000
	bmask    common.Color = 0b000000000000000011111000
)

// to convert 24-bit color to 16-bit color, we are taking only the 5 (for red and
// blue) or 6 (for green) MS bits
func toHatColor(c common.Color) color.Color {
	r := color.Color((c & rmask) >> 8)
	g := color.Color((c & gmask) >> 5)
	b := color.Color((c & bmask) >> 3)

	return r | g | b
}

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
	done   chan struct{}
	input  *stick.Device
}

func NewHat(joystickEvents chan<- Event, screenEvents <-chan DisplayMessage) *Hat {
	return &Hat{
		events: joystickEvents,
		screen: screenEvents,
		done:   make(chan struct{}),
	}
}

func (h *Hat) Start() {
	h.init()
	go h.do()
}

func (h *Hat) Stop() {
	close(h.done)
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
			fb.SetPixel(x, y, toHatColor(screenChange.Screen[y][x]))
		}
	}

	cursorOrigColor := toHatColor(screenChange.Screen[screenChange.CursorY][screenChange.CursorX])
	cursorColor := reversColor(cursorOrigColor)

	fb.SetPixel(int(screenChange.CursorX), int(screenChange.CursorY), cursorColor)
	err := screen.Draw(fb)
	if err != nil {
		log.Println("error while printing to HAT display:", err)
	}
}

func reversColor(c color.Color) color.Color {
	return c ^ 0b1111111111111111
}

func (h Hat) gracefulShutDown() {
	screen.Clear()
	// signal the controller we've done
	close(h.events)
}
