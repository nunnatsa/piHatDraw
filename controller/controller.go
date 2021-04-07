package controller

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nunnatsa/piHatDraw/common"

	"github.com/nunnatsa/piHatDraw/notifier"
	"github.com/nunnatsa/piHatDraw/webapp"

	"github.com/nunnatsa/piHatDraw/hat"
	"github.com/nunnatsa/piHatDraw/state"
)

type Controller struct {
	hat            hat.Interface
	joystickEvents chan hat.Event
	screenEvents   chan hat.DisplayMessage
	done           chan bool
	state          *state.State
	notifier       *notifier.Notifier
	clientEvents   <-chan webapp.ClientEvent
}

func NewController(notifier *notifier.Notifier, clientEvents <-chan webapp.ClientEvent, canvasWidth uint8, canvasHeight uint8) *Controller {
	je := make(chan hat.Event)
	se := make(chan hat.DisplayMessage)

	return &Controller{
		hat:            hat.NewHat(je, se),
		joystickEvents: je,
		screenEvents:   se,
		done:           make(chan bool),
		state:          state.NewState(canvasWidth, canvasHeight),
		notifier:       notifier,
		clientEvents:   clientEvents,
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

	defer c.stop(signals)

	c.hat.Start()

	msg := c.state.CreateDisplayMessage()
	c.screenEvents <- msg

	for {
		changed := false

		select {
		case <-signals:
			return

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

		case e := <-c.clientEvents:
			switch data := e.(type) {
			case webapp.ClientEventRegistered:
				id := uint64(data)
				c.registered(id)

			case webapp.ClientEventReset:
				if data {
					c.state.Reset()
					changed = true
				}

			case webapp.ClientEventSetColor:
				color := common.Color(data)
				changed = c.state.SetColor(color)

			case webapp.ClientEventSetTool:
				switch string(data) {
				case "pen":
					changed = c.state.SetPen()
				case "eraser":
					changed = c.state.SetEraser()
				default:
					log.Printf(`unknown tool "%s"`, data)
				}

			case webapp.ClientEventDownload:
				ch := chan [][]common.Color(data)
				ch <- c.state.Canvas.Clone()
			}
		}

		if changed {
			c.Update()
		}
	}
}

func (c *Controller) stop(signals chan os.Signal) {
	c.hat.Stop()
	<-c.joystickEvents // wait for the hat graceful shutdown
	signal.Stop(signals)
	close(c.done)
}

func (c *Controller) Update() {
	msg := c.state.CreateDisplayMessage()
	go func() {
		c.screenEvents <- msg
	}()

	js, err := json.Marshal(c.state)
	if err != nil {
		log.Println(err)
	} else {
		c.notifier.NotifyAll(js)
	}
}

func (c *Controller) registered(id uint64) {
	js, err := json.Marshal(c.state)
	if err != nil {
		log.Println(err)
	} else {
		c.notifier.NotifyOne(id, js)
	}
}
