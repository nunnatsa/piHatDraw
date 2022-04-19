package controller

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/hat"
	"github.com/nunnatsa/piHatDraw/notifier"
	"github.com/nunnatsa/piHatDraw/state"
	"github.com/nunnatsa/piHatDraw/webapp"
)

const (
	eraserToolName = "eraser"
	bucketToolName = "bucket"
	penToolName    = "pen"

	canvasWidth  = 40
	canvasHeight = 24
	y            = uint8(canvasHeight / 2)
	x            = uint8(canvasWidth / 2)
)

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")

}

var _ = Describe("Controller Test", func() {
	s := state.NewState(canvasWidth, canvasHeight)
	je := make(chan hat.Event, 1)
	se := make(chan hat.DisplayMessage, 1)
	n := notifier.NewNotifier()
	ce := make(chan webapp.ClientEvent, 1)

	hatMock := &hatMock{
		je: je,
		se: se,
	}

	reg1 := make(chan []byte, 1)
	reg2 := make(chan []byte, 1)

	client1 := n.Subscribe(reg1)
	client2 := n.Subscribe(reg2)

	done := make(chan struct{})
	defer close(done)

	c := &Controller{
		hat:            hatMock,
		joystickEvents: je,
		screenEvents:   se,
		done:           done,
		state:          s,
		notifier:       n,
		clientEvents:   ce,
	}

	c.Start()

	ce <- webapp.ClientEventRegistered(client1)
	ce <- webapp.ClientEventRegistered(client2)

	It("should build messages with the initial values", func() {
		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(4))
			Expect(msg.CursorY).Should(BeEquivalentTo(4))
			Expect(msg.Screen).To(HaveLen(8))
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg1, x, y)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg2, x, y)
		}).Should(BeTrue())
	})

	It("should do nothing for the undo command, when there was no change yet", func() {
		ce <- webapp.ClientEventUndo(true)
		Consistently(se).ShouldNot(Receive())
		Consistently(reg1).ShouldNot(Receive())
		Consistently(reg2).ShouldNot(Receive())
	})

	It("should move down", func() {
		hatMock.MoveDown()

		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(4))
			Expect(msg.CursorY).Should(BeEquivalentTo(5))
			Expect(msg.Screen).To(HaveLen(8))
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg1, x, y+1)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg2, x, y+1)
		}).Should(BeTrue())

		By("should do nothing for the undo command. no change in the canvas")
		ce <- webapp.ClientEventUndo(true)
		Consistently(se).ShouldNot(Receive())
		Consistently(reg1).ShouldNot(Receive())
		Consistently(reg2).ShouldNot(Receive())
	})

	It("should move up", func() {
		hatMock.MoveUp()

		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(4))
			Expect(msg.CursorY).Should(BeEquivalentTo(4))
			Expect(msg.Screen).To(HaveLen(8))
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg1, x, y)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg2, x, y)
		}).Should(BeTrue())

		By("should do nothing for the undo command for move event")
		ce <- webapp.ClientEventUndo(true)
		Consistently(se).ShouldNot(Receive())
		Consistently(reg1).ShouldNot(Receive())
		Consistently(reg2).ShouldNot(Receive())
	})

	It("should move right", func() {
		hatMock.MoveRight()

		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(5))
			Expect(msg.CursorY).Should(BeEquivalentTo(4))
			Expect(msg.Screen).To(HaveLen(8))
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg1, x+1, y)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg2, x+1, y)
		}).Should(BeTrue())

		By("should do nothing for the undo command for move event")
		ce <- webapp.ClientEventUndo(true)
		Consistently(se).ShouldNot(Receive())
		Consistently(reg1).ShouldNot(Receive())
		Consistently(reg2).ShouldNot(Receive())
	})

	It("should move left", func() {
		hatMock.MoveLeft()

		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(4))
			Expect(msg.CursorY).Should(BeEquivalentTo(4))
			Expect(msg.Screen).To(HaveLen(8))
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg1, x, y)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkMoveNotifications(<-reg2, x, y)
		}).Should(BeTrue())

		By("should do nothing for the undo command for move event")
		ce <- webapp.ClientEventUndo(true)
		Consistently(se).ShouldNot(Receive())
		Consistently(reg1).ShouldNot(Receive())
		Consistently(reg2).ShouldNot(Receive())
	})

	It("should paint", func() {
		hatMock.Press()

		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(4))
			Expect(msg.CursorY).Should(BeEquivalentTo(4))
			Expect(msg.Screen).To(HaveLen(8))
			Expect(msg.Screen[4][4]).Should(BeEquivalentTo(0xFFFFFF))
			return true
		}).Should(BeTrue())

		expected := state.Pixel{X: x, Y: y, Color: 0xFFFFFF}
		Eventually(func() bool {
			return checkPaintNotifications(<-reg1, expected)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkPaintNotifications(<-reg2, expected)
		}).Should(BeTrue())

		By("should undo painting")
		ce <- webapp.ClientEventUndo(true)
		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(4))
			Expect(msg.CursorY).Should(BeEquivalentTo(4))
			Expect(msg.Screen).To(HaveLen(8))
			Expect(msg.Screen[4][4]).Should(BeEquivalentTo(0))
			return true
		}).Should(BeTrue())

		expected = state.Pixel{X: x, Y: y, Color: 0}
		Eventually(func() bool {
			return checkPaintNotifications(<-reg1, expected)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkPaintNotifications(<-reg2, expected)
		}).Should(BeTrue())
	})

	It("should set color", func() {
		clr := common.Color(0x123456)
		ce <- webapp.ClientEventSetColor(clr)

		Eventually(c.screenEvents).Should(Receive())

		Eventually(func() bool {
			webMsg, err := getChangeFromMsg(<-reg1)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, *webMsg.Color).To(Equal(clr))
			ExpectWithOffset(1, webMsg.Cursor).To(BeNil())
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			webMsg, err := getChangeFromMsg(<-reg2)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, *webMsg.Color).To(Equal(clr))
			ExpectWithOffset(1, webMsg.Cursor).To(BeNil())
			return true
		}).Should(BeTrue())

		By("do nothing if the color was not changed")
		ce <- webapp.ClientEventSetColor(clr)
		Consistently(c.screenEvents).ShouldNot(Receive())
		Consistently(reg1).ShouldNot(Receive())
		Consistently(reg2).ShouldNot(Receive())

	})

	It("should set tool to eraser", func() {
		ce <- webapp.ClientEventSetTool(eraserToolName)

		Eventually(c.screenEvents).Should(Receive())

		Eventually(func() bool {
			webMsg, err := getChangeFromMsg(<-reg1)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, webMsg.ToolName).To(Equal(eraserToolName))
			ExpectWithOffset(1, webMsg.Cursor).To(BeNil())
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			webMsg, err := getChangeFromMsg(<-reg2)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, webMsg.ToolName).To(Equal(eraserToolName))
			ExpectWithOffset(1, webMsg.Cursor).To(BeNil())
			return true
		}).Should(BeTrue())
	})

	It("should set tool to bucket", func() {
		ce <- webapp.ClientEventSetTool(bucketToolName)

		Eventually(c.screenEvents).Should(Receive())

		Eventually(func() bool {
			webMsg, err := getChangeFromMsg(<-reg1)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, webMsg.ToolName).To(Equal(bucketToolName))
			ExpectWithOffset(1, webMsg.Cursor).To(BeNil())
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			webMsg, err := getChangeFromMsg(<-reg2)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, webMsg.ToolName).To(Equal(bucketToolName))
			ExpectWithOffset(1, webMsg.Cursor).To(BeNil())
			return true
		}).Should(BeTrue())
	})

	It("should set tool to pen", func() {
		ce <- webapp.ClientEventSetTool(penToolName)

		Eventually(c.screenEvents).Should(Receive())

		Eventually(func() bool {
			webMsg, err := getChangeFromMsg(<-reg1)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, webMsg.ToolName).To(Equal(penToolName))
			ExpectWithOffset(1, webMsg.Cursor).To(BeNil())
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			webMsg, err := getChangeFromMsg(<-reg2)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, webMsg.ToolName).To(Equal(penToolName))
			ExpectWithOffset(1, webMsg.Cursor).To(BeNil())
			return true
		}).Should(BeTrue())

		By("should do nothing if the tool was not changed", func() {
			ce <- webapp.ClientEventSetTool(penToolName)

			Consistently(c.screenEvents).ShouldNot(Receive())
			Consistently(reg1).ShouldNot(Receive())
			Consistently(reg2).ShouldNot(Receive())
		})
	})

	It("should ignore unknown tools", func() {
		ce <- webapp.ClientEventSetTool("wrongToolName")

		Consistently(c.screenEvents).ShouldNot(Receive())
		Consistently(reg1).ShouldNot(Receive())
		Consistently(reg2).ShouldNot(Receive())
	})

	It("should use bucket", func() {
		clr := common.Color(0x00112233)
		change, err := c.state.SetTool(bucketToolName)
		Expect(err).ToNot(HaveOccurred())
		Expect(change.ToolName).Should(Equal(bucketToolName))

		change = c.state.SetColor(clr)
		Expect(*change.Color).Should(Equal(clr))

		_ = c.state.GoUp()
		_ = c.state.GoUp()
		change = c.state.GoUp()
		Expect(change.Cursor.Y).Should(Equal(y - 3))
		_ = c.state.GoRight()
		_ = c.state.GoRight()
		_ = c.state.GoRight()
		change = c.state.GoRight()
		Expect(change.Cursor.X).Should(Equal(x + 4))

		hatMock.Press()

		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(7))
			Expect(msg.CursorY).Should(BeEquivalentTo(1))
			Expect(msg.Screen).To(HaveLen(8))
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					Expect(msg.Screen[4][4]).Should(BeEquivalentTo(clr))
				}
			}
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkBucketNotifications(<-reg1, clr)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkBucketNotifications(<-reg2, clr)
		}).Should(BeTrue())

		By("check reset")
		ce <- webapp.ClientEventReset(true)
		Eventually(func() bool {
			msg := <-c.screenEvents
			Expect(msg.CursorX).Should(BeEquivalentTo(4))
			Expect(msg.CursorY).Should(BeEquivalentTo(4))
			Expect(msg.Screen).To(HaveLen(8))
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					Expect(msg.Screen[4][4]).Should(BeEquivalentTo(0))
				}
			}
			return true
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkResetNotifications(<-reg1)
		}).Should(BeTrue())

		Eventually(func() bool {
			return checkResetNotifications(<-reg2)
		}).Should(BeTrue())

	})
})

func checkMoveNotifications(msg []byte, x uint8, y uint8) bool {
	webMsg, err := getChangeFromMsg(msg)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	ExpectWithOffset(1, webMsg.Cursor.X).To(Equal(x))
	ExpectWithOffset(1, webMsg.Cursor.Y).To(Equal(y))

	return true
}

func checkPaintNotifications(msg []byte, pixels ...state.Pixel) bool {
	webMsg, err := getChangeFromMsg(msg)
	ExpectWithOffset(1, err).ShouldNot(HaveOccurred())

	ExpectWithOffset(1, webMsg.Pixels).To(HaveLen(len(pixels)))

	for i, p := range pixels {
		mp := webMsg.Pixels[i]
		Expect(mp).Should(Equal(p))
	}

	return true
}

func checkBucketNotifications(msg []byte, clr common.Color) bool {
	webMsg, err := getChangeFromMsg(msg)
	ExpectWithOffset(1, err).ShouldNot(HaveOccurred())

	ExpectWithOffset(1, webMsg.Pixels).To(HaveLen(canvasWidth * canvasHeight))

	checkBoard := [][]bool{
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	}

	for _, px := range webMsg.Pixels {
		Expect(px.Color).Should(Equal(clr))
		checkBoard[px.Y][px.X] = true
	}

	for y := 0; y < canvasHeight; y++ {
		for x := 0; x < canvasWidth; x++ {
			Expect(checkBoard[y][x]).Should(BeTrue())
		}
	}

	return true
}

func getChangeFromMsg(msg []byte) (*state.Change, error) {
	s := &state.Change{}
	if err := json.Unmarshal(msg, s); err != nil {
		return nil, err
	}
	return s, nil
}

func checkResetNotifications(msg []byte) bool {
	webMsg, err := getChangeFromMsg(msg)
	ExpectWithOffset(1, err).ShouldNot(HaveOccurred())

	ExpectWithOffset(1, webMsg.Canvas).To(HaveLen(canvasHeight))
	for _, line := range webMsg.Canvas {
		ExpectWithOffset(1, line).To(HaveLen(canvasWidth))
		for _, px := range line {
			ExpectWithOffset(1, px).Should(BeEquivalentTo(0))
		}
	}
	ExpectWithOffset(1, *webMsg.Color).Should(BeEquivalentTo(0xFFFFFF))
	ExpectWithOffset(1, webMsg.Cursor).ShouldNot(BeNil())
	ExpectWithOffset(1, webMsg.Cursor.X).Should(BeEquivalentTo(x))
	ExpectWithOffset(1, webMsg.Cursor.Y).Should(BeEquivalentTo(y))
	ExpectWithOffset(1, webMsg.ToolName).Should(BeEquivalentTo(penToolName))
	ExpectWithOffset(1, webMsg.Window).ShouldNot(BeNil())
	ExpectWithOffset(1, webMsg.Window.X).Should(BeEquivalentTo(x - 4))
	ExpectWithOffset(1, webMsg.Window.Y).Should(BeEquivalentTo(y - 4))

	return true
}
