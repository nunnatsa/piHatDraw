package controller

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nunnatsa/piHatDraw/hat"
	"github.com/nunnatsa/piHatDraw/state"
)

type Controller struct {
	hat            hat.Interface
	joystickEvents chan hat.Event
	screenEvents   chan hat.DisplayMessage
	done           chan bool
	state          *state.State
}

func NewController() *Controller {
	je := make(chan hat.Event)
	se := make(chan hat.DisplayMessage)

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

	msg := c.state.CreateDisplayMessage()
	c.screenEvents <- msg

	for {
		changed := false

		select {
		case <-signals:
			close(c.done)

		case je := <-c.joystickEvents:
			switch je {
			case hat.MoveUp:
				changed = c.state.GoUp()

			case hat.MoveLeft:
				changed = c.state.GoLeft()

			case hat.MoveDown:
				changed = c.state.GoDown()

			case hat.MoveRight:
				changed = c.state.GoRight()

			case hat.Pressed:
				changed = c.state.PaintPixel()
			}
		}

		if changed {
			msg := c.state.CreateDisplayMessage()
			c.screenEvents <- msg
		}
	}
}
