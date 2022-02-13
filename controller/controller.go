package controller

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/hat"
	"github.com/nunnatsa/piHatDraw/notifier"
	"github.com/nunnatsa/piHatDraw/state"
	"github.com/nunnatsa/piHatDraw/webapp"
)

type Controller struct {
	hat            hat.Interface
	joystickEvents chan hat.Event
	screenEvents   chan hat.DisplayMessage
	done           chan struct{}
	state          *state.State
	notifier       *notifier.Notifier
	clientEvents   <-chan webapp.ClientEvent
}

func NewController(notifier *notifier.Notifier, clientEvents <-chan webapp.ClientEvent, canvasWidth uint8, canvasHeight uint8) *Controller {
	je := make(chan hat.Event, 1)
	se := make(chan hat.DisplayMessage, 1)

	return &Controller{
		hat:            hat.NewHat(je, se),
		joystickEvents: je,
		screenEvents:   se,
		done:           make(chan struct{}),
		state:          state.NewState(canvasWidth, canvasHeight),
		notifier:       notifier,
		clientEvents:   clientEvents,
	}
}

func (c Controller) Start() <-chan struct{} {
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
		var change *state.Change = nil

		select {
		case <-signals:
			return

		case je := <-c.joystickEvents:
			change = c.handleJoystickEvent(je)

		case e := <-c.clientEvents:
			change = c.handleWebClientEvent(e)
		}

		if change != nil {
			c.Update(change)
		}
	}
}

func (c *Controller) handleWebClientEvent(e webapp.ClientEvent) *state.Change {
	switch data := e.(type) {
	case webapp.ClientEventRegistered:
		id := uint64(data)
		c.registered(id)

	case webapp.ClientEventReset:
		if data {
			return c.state.Reset()
		}

	case webapp.ClientEventSetColor:
		color := common.Color(data)
		return c.state.SetColor(color)

	case webapp.ClientEventSetTool:
		var err error
		change, err := c.state.SetTool(string(data))
		if err != nil {
			log.Printf(err.Error())
			return nil
		}
		return change

	case webapp.ClientEventDownload:
		ch := chan [][]common.Color(data)
		ch <- c.state.GetCanvasClone()

	case webapp.ClientEventUndo:
		return c.state.Undo()
	}

	return nil
}

func (c *Controller) handleJoystickEvent(je hat.Event) *state.Change {
	switch je {
	case hat.MoveUp:
		return c.state.GoUp()

	case hat.MoveLeft:
		return c.state.GoLeft()

	case hat.MoveDown:
		return c.state.GoDown()

	case hat.MoveRight:
		return c.state.GoRight()

	case hat.Pressed:
		return c.state.Paint()
	}
	return nil
}

func (c *Controller) stop(signals chan os.Signal) {
	c.hat.Stop()
	<-c.joystickEvents // wait for the hat graceful shutdown
	signal.Stop(signals)
	close(c.done)
}

func (c *Controller) Update(change *state.Change) {
	msg := c.state.CreateDisplayMessage()
	go func() {
		c.screenEvents <- msg
	}()

	js, err := json.Marshal(change)
	if err != nil {
		log.Println(err)
	} else {
		c.notifier.NotifyAll(js)
	}
}

func (c *Controller) registered(id uint64) {
	change := c.state.GetFullChange()
	js, err := json.Marshal(change)
	if err != nil {
		log.Println(err)
	} else {
		c.notifier.NotifyOne(id, js)
	}
}
