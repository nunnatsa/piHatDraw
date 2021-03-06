package controller

import (
	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/hat"
	"github.com/nunnatsa/piHatDraw/state"
	"os"
	"os/signal"
	"syscall"
)

type Controller struct {
	hat            hat.Interface
	joystickEvents chan common.HatEvent
	screenEvents   chan *common.DisplayMessage
	done           chan bool
	state          *state.State
}

func NewController() *Controller {
	je := make(chan common.HatEvent)
	se := make(chan *common.DisplayMessage)

	return &Controller{
		hat:            hat.NewHat(je, se),
		joystickEvents: je,
		screenEvents:   se,
		done:           make(chan bool),
		state:          state.NewState(),
	}
}

func (c Controller) Start() <-chan bool {
	go c.do()
	return c.done
}

func (c *Controller) do() {
	// Set up a signals channel (stop the loop using Ctrl-C)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	c.hat.Start()

	for {
		changed := false

		select {
		case <-signals:
			close(c.done)

		case je := <-c.joystickEvents:
			switch je {
			case common.MoveUp:
				changed = c.state.GoUp()

			case common.MoveLeft:
				changed = c.state.GoLeft()

			case common.MoveDown:
				changed = c.state.GoDown()

			case common.MoveRight:
				changed = c.state.GoRight()

			case common.Pressed:
				changed = c.state.PaintPixel()
			}
		}

		if changed {
			msg := c.state.CreateDisplayMessage()
			c.screenEvents <- msg
		}
	}
}
