package state

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nunnatsa/piHatDraw/common"
)

const (
	canvasWidth  = uint8(40)
	canvasHeight = uint8(24)
)

var _ = Describe("test state", func() {
	Context("should create a new state", func() {
		s := NewState(canvasWidth, canvasHeight)

		It("should create canvas", func() {
			Expect(s.canvasHeight).Should(Equal(canvasHeight))
			Expect(s.canvasWidth).Should(Equal(canvasWidth))

			Expect(s.canvas).Should(HaveLen(int(canvasHeight)))
			for _, line := range s.canvas {
				Expect(line).Should(HaveLen(int(canvasWidth)))
			}
		})

		It("should place the cursor in the middle of the canvas", func() {
			Expect(s.cursor.X).Should(BeEquivalentTo(canvasWidth / 2))
			Expect(s.cursor.Y).Should(BeEquivalentTo(canvasHeight / 2))
		})

		It("should set the default color to white", func() {
			Expect(s.color).Should(Equal(common.Color(0xffffff)))
		})

		It("should set the default tool to pen", func() {
			Expect(s.toolName).Should(Equal("pen"))
		})

		It("should place the window in the middle of the canvas", func() {
			Expect(s.window.X).Should(BeEquivalentTo(canvasWidth/2 - 4))
			Expect(s.window.Y).Should(BeEquivalentTo(canvasHeight/2 - 4))
		})
	})

	Context("test CreateDisplayMessage", func() {
		s := NewState(canvasWidth, canvasHeight)

		s.cursor.X = 7
		s.cursor.Y = 3

		s.window.X = 0
		s.window.Y = 0

		s.canvas[4][3] = common.Color(0xAABBCC)

		msg := s.CreateDisplayMessage()

		It("Should set the pixels according to the window", func() {
			Expect(msg.Screen[4][3]).Should(Equal(common.Color(0xAABBCC)))
		})

		It("Should set the cursor, relatively to the window", func() {
			Expect(msg.CursorX).Should(BeEquivalentTo(7))
			Expect(msg.CursorY).Should(BeEquivalentTo(3))
		})
	})

	Context("test GoUp", func() {
		s := NewState(canvasWidth, canvasHeight)

		x := s.cursor.X
		y := s.cursor.Y

		It("should update the cursor", func() {
			change := s.GoUp()
			Expect(change).ShouldNot(BeNil())
			Expect(s.cursor.Y).Should(BeEquivalentTo(y - 1))
			Expect(s.cursor.X).Should(BeEquivalentTo(x))
			Expect(change.Cursor).ShouldNot(BeNil())
			Expect(change.Cursor.Y).Should(BeEquivalentTo(y - 1))
			Expect(change.Cursor.X).Should(BeEquivalentTo(x))

			Expect(undoList.len()).Should(BeZero())
		})

		It("should ignore if the cursor is at the top of the canvas", func() {
			s.cursor.Y = 0
			change := s.GoUp()
			Expect(change).Should(BeNil())
			Expect(s.cursor.Y).Should(BeEquivalentTo(0))
			Expect(s.cursor.X).Should(BeEquivalentTo(x))

			Expect(undoList.len()).Should(BeZero())
		})
	})

	Context("test GoDown", func() {
		s := NewState(canvasWidth, canvasHeight)

		x := s.cursor.X
		y := s.cursor.Y

		It("should update the cursor", func() {
			change := s.GoDown()
			Expect(change).ShouldNot(BeNil())
			Expect(s.cursor.Y).Should(BeEquivalentTo(y + 1))
			Expect(s.cursor.X).Should(BeEquivalentTo(x))
			Expect(change.Cursor).ShouldNot(BeNil())
			Expect(change.Cursor.Y).Should(BeEquivalentTo(y + 1))
			Expect(change.Cursor.X).Should(BeEquivalentTo(x))

			Expect(undoList.len()).Should(BeZero())
		})

		It("should ignore if the cursor is at the bottom of the canvas", func() {
			s.cursor.Y = canvasHeight - 1
			change := s.GoDown()
			Expect(change).Should(BeNil())
			Expect(s.cursor.Y).Should(BeEquivalentTo(canvasHeight - 1))
			Expect(s.cursor.X).Should(BeEquivalentTo(x))

			Expect(undoList.len()).Should(BeZero())
		})
	})

	Context("test GoLeft", func() {
		s := NewState(canvasWidth, canvasHeight)

		x := s.cursor.X
		y := s.cursor.Y

		It("should update the cursor", func() {
			change := s.GoLeft()
			Expect(change).ShouldNot(BeNil())
			Expect(s.cursor.Y).Should(BeEquivalentTo(y))
			Expect(s.cursor.X).Should(BeEquivalentTo(x - 1))
			Expect(change.Cursor).ShouldNot(BeNil())
			Expect(change.Cursor.Y).Should(BeEquivalentTo(y))
			Expect(change.Cursor.X).Should(BeEquivalentTo(x - 1))

			Expect(undoList.len()).Should(BeZero())
		})

		It("should ignore if the cursor is at the most left column of the canvas", func() {
			s.cursor.X = 0
			change := s.GoLeft()
			Expect(change).Should(BeNil())
			Expect(s.cursor.Y).Should(BeEquivalentTo(y))
			Expect(s.cursor.X).Should(BeEquivalentTo(0))

			Expect(undoList.len()).Should(BeZero())
		})
	})

	Context("test GoRight", func() {
		s := NewState(canvasWidth, canvasHeight)

		x := s.cursor.X
		y := s.cursor.Y

		It("should update the cursor", func() {
			change := s.GoRight()
			Expect(change).ShouldNot(BeNil())
			Expect(s.cursor.Y).Should(BeEquivalentTo(y))
			Expect(s.cursor.X).Should(BeEquivalentTo(x + 1))
			Expect(change.Cursor).ShouldNot(BeNil())
			Expect(change.Cursor.Y).Should(BeEquivalentTo(y))
			Expect(change.Cursor.X).Should(BeEquivalentTo(x + 1))

			Expect(undoList.len()).Should(BeZero())
		})

		It("should ignore if the cursor is at the most right column of the canvas", func() {
			s.cursor.X = canvasWidth - 1
			change := s.GoRight()
			Expect(change).Should(BeNil())
			Expect(s.cursor.Y).Should(BeEquivalentTo(y))
			Expect(s.cursor.X).Should(BeEquivalentTo(canvasWidth - 1))

			Expect(undoList.len()).Should(BeZero())
		})
	})

	Context("test StatePaintPixel", func() {
		var s *State

		BeforeEach(func() {
			s = NewState(canvasWidth, canvasHeight)
		})

		AfterEach(func() {
			emptyUndoList()
		})

		It("should be black when creating the state", func() {
			Expect(s.canvas[s.cursor.Y][s.cursor.X]).Should(Equal(common.Color(0)))
		})

		It("should set when painting", func() {
			By("painting")
			change := s.Paint()
			Expect(change).ToNot(BeNil())
			Expect(s.canvas[s.cursor.Y][s.cursor.X]).Should(Equal(common.Color(0xFFFFFF)))
			Expect(change.Pixels).Should(HaveLen(1))
			Expect(change.Pixels[0]).Should(Equal(Pixel{X: s.cursor.X, Y: s.cursor.Y, Color: common.Color(0xFFFFFF)}))
			Expect(undoList.len()).Should(Equal(1))

			By("undo")
			change = undoList.pop()
			Expect(change.Pixels).Should(HaveLen(1))
			Expect(change.Pixels[0]).Should(Equal(Pixel{X: s.cursor.X, Y: s.cursor.Y, Color: common.Color(0)}))
			Expect(undoList.len()).Should(BeZero())
		})

		It("should ignore if painting with the same color at the same place", func() {
			change := s.Paint()
			Expect(change).ToNot(BeNil())
			_ = undoList.pop()

			change = s.Paint()
			Expect(change).To(BeNil())
			Expect(s.canvas[s.cursor.Y][s.cursor.X]).Should(Equal(common.Color(0xFFFFFF)))
		})

		It("should ignore if painting out of canvas", func() {
			s.cursor.X = canvasWidth
			change := s.Paint()
			Expect(change).To(BeNil())
			Expect(undoList.len()).Should(BeZero())

			s.cursor.X = canvasWidth / 2
			s.cursor.Y = canvasHeight
			change = s.Paint()
			Expect(change).To(BeNil())
		})
	})

	Context("Test SetColor", func() {
		s := NewState(8, 8)
		s.color = 0x123456

		It("should ignore if setting the same color", func() {
			change := s.SetColor(0x123456)
			Expect(change).Should(BeNil())
		})

		It("should set color", func() {
			change := s.SetColor(0x654321)
			Expect(change).ShouldNot(BeNil())
			Expect(change.Color).ShouldNot(BeNil())
			Expect(*change.Color).Should(Equal(common.Color(0x654321)))
		})
	})

	Context("Test SetTool", func() {
		var s *State
		BeforeEach(func() {
			s = NewState(8, 8)
		})

		It("should ignore if setting the same tool", func() {
			change, err := s.SetTool(penName)
			Expect(err).ToNot(HaveOccurred())
			Expect(change).Should(BeNil())
			Expect(s.toolName).Should(Equal(penName))
		})

		It("should set tool to eraser", func() {
			change, err := s.SetTool(eraserName)
			Expect(err).ToNot(HaveOccurred())
			Expect(*change).Should(Equal(Change{ToolName: eraserName}))
			Expect(s.toolName).Should(Equal(eraserName))
		})

		It("should set tool to bucket", func() {
			change, err := s.SetTool(bucketName)
			Expect(err).ToNot(HaveOccurred())
			Expect(*change).Should(Equal(Change{ToolName: bucketName}))
			Expect(s.toolName).Should(Equal(bucketName))
		})

		It("should set the tool to pen", func() {
			s.toolName = bucketName
			s.tool = s.bucket

			change, err := s.SetTool(penName)
			Expect(err).ToNot(HaveOccurred())
			Expect(*change).Should(Equal(Change{ToolName: penName}))
			Expect(s.toolName).Should(Equal(penName))

		})

		It("should reject unknown tools", func() {
			change, err := s.SetTool("wrongToolName")
			Expect(err).To(HaveOccurred())
			Expect(change).Should(BeNil())
		})
	})

	Context("test undo", func() {
		s := NewState(canvasWidth, canvasHeight)

		It("should perform unod", func() {
			By("check the initial values")
			Expect(s.canvas[3][4]).Should(BeEquivalentTo(0))
			Expect(s.canvas[10][20]).Should(BeEquivalentTo(0))

			c := &Change{
				Pixels: []Pixel{{
					X: 4, Y: 3, Color: 0x112233,
				}, {
					X: 20, Y: 10, Color: 0x112233,
				}},
			}

			undoList.push(c)
			By("perform the undo")
			s.Undo()

			Expect(s.canvas[3][4]).Should(BeEquivalentTo(0x112233))
			Expect(s.canvas[10][20]).Should(BeEquivalentTo(0x112233))
		})
	})

	Context("test bucket", func() {
		var s *State

		BeforeEach(func() {
			s = NewState(8, 8)
		})

		It("should paint inner-closed shape", func() {
			s.cursor.Y = 4
			s.cursor.X = 4

			s.canvas = Canvas{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 1, 1, 1, 1, 1, 0},
				{0, 0, 1, 0, 0, 0, 1, 0},
				{0, 0, 1, 0, 0, 0, 1, 0},
				{0, 0, 1, 0, 0, 0, 1, 0},
				{0, 0, 1, 1, 1, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			}

			_, _ = s.SetTool(bucketName)
			s.color = 2
			s.Paint()

			expected := Canvas{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 1, 1, 1, 1, 1, 0},
				{0, 0, 1, 2, 2, 2, 1, 0},
				{0, 0, 1, 2, 2, 2, 1, 0},
				{0, 0, 1, 2, 2, 2, 1, 0},
				{0, 0, 1, 1, 1, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			}

			Expect(s.canvas).Should(Equal(expected))
		})

		It("should paint border of a shape", func() {
			s.cursor.Y = 2
			s.cursor.X = 2

			_, _ = s.SetTool(bucketName)
			s.color = 3

			s.canvas = Canvas{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 1, 1, 1, 1, 1, 0},
				{0, 0, 1, 2, 2, 2, 1, 0},
				{0, 0, 1, 2, 2, 2, 1, 0},
				{0, 0, 1, 2, 2, 2, 1, 0},
				{0, 0, 1, 1, 1, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			}

			s.Paint()
			expected := Canvas{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 3, 3, 3, 3, 3, 0},
				{0, 0, 3, 2, 2, 2, 3, 0},
				{0, 0, 3, 2, 2, 2, 3, 0},
				{0, 0, 3, 2, 2, 2, 3, 0},
				{0, 0, 3, 3, 3, 3, 3, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			}

			Expect(s.canvas).Should(Equal(expected))
		})

		It("should paint out of a shape", func() {
			s.cursor.X = 0
			s.cursor.Y = 0

			_, _ = s.SetTool(bucketName)
			s.color = 4

			s.canvas = Canvas{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 3, 3, 3, 3, 3, 0},
				{0, 0, 3, 2, 2, 2, 3, 0},
				{0, 0, 3, 2, 2, 2, 3, 0},
				{0, 0, 3, 2, 2, 2, 3, 0},
				{0, 0, 3, 3, 3, 3, 3, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			}

			s.Paint()
			expected := Canvas{
				{4, 4, 4, 4, 4, 4, 4, 4},
				{4, 4, 4, 4, 4, 4, 4, 4},
				{4, 4, 3, 3, 3, 3, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 3, 3, 3, 3, 4},
				{4, 4, 4, 4, 4, 4, 4, 4},
			}

			Expect(s.canvas).Should(Equal(expected))
		})

		It("should paint out of a shape", func() {
			s.cursor.X = 7
			s.cursor.Y = 7

			_, _ = s.SetTool(bucketName)
			s.color = 5

			s.canvas = Canvas{
				{4, 4, 4, 4, 4, 4, 4, 4},
				{4, 4, 4, 4, 4, 4, 4, 4},
				{4, 4, 3, 3, 3, 3, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 3, 3, 3, 3, 4},
				{4, 4, 4, 4, 4, 4, 4, 4},
			}

			s.Paint()
			expected := Canvas{
				{5, 5, 5, 5, 5, 5, 5, 5},
				{5, 5, 5, 5, 5, 5, 5, 5},
				{5, 5, 3, 3, 3, 3, 3, 5},
				{5, 5, 3, 2, 2, 2, 3, 5},
				{5, 5, 3, 2, 2, 2, 3, 5},
				{5, 5, 3, 2, 2, 2, 3, 5},
				{5, 5, 3, 3, 3, 3, 3, 5},
				{5, 5, 5, 5, 5, 5, 5, 5},
			}

			Expect(s.canvas).Should(Equal(expected))
		})

		It("should undo", func() {
			s.cursor.X = 7
			s.cursor.Y = 7

			_, _ = s.SetTool(bucketName)

			s.canvas = Canvas{
				{5, 5, 5, 5, 5, 5, 5, 5},
				{5, 5, 5, 5, 5, 5, 5, 5},
				{5, 5, 3, 3, 3, 3, 3, 5},
				{5, 5, 3, 2, 2, 2, 3, 5},
				{5, 5, 3, 2, 2, 2, 3, 5},
				{5, 5, 3, 2, 2, 2, 3, 5},
				{5, 5, 3, 3, 3, 3, 3, 5},
				{5, 5, 5, 5, 5, 5, 5, 5},
			}

			s.Undo()

			expected := Canvas{
				{4, 4, 4, 4, 4, 4, 4, 4},
				{4, 4, 4, 4, 4, 4, 4, 4},
				{4, 4, 3, 3, 3, 3, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 2, 2, 2, 3, 4},
				{4, 4, 3, 3, 3, 3, 3, 4},
				{4, 4, 4, 4, 4, 4, 4, 4},
			}

			Expect(s.canvas).Should(Equal(expected))
		})

		It("fill the entire canvas", func() {
			s.cursor.X = 7
			s.cursor.Y = 7

			_, _ = s.SetTool(bucketName)

			s.canvas = Canvas{
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0},
			}

			s.SetColor(6)
			s.Paint()

			expected := Canvas{
				{6, 6, 6, 6, 6, 6, 6, 6},
				{6, 6, 6, 6, 6, 6, 6, 6},
				{6, 6, 6, 6, 6, 6, 6, 6},
				{6, 6, 6, 6, 6, 6, 6, 6},
				{6, 6, 6, 6, 6, 6, 6, 6},
				{6, 6, 6, 6, 6, 6, 6, 6},
				{6, 6, 6, 6, 6, 6, 6, 6},
				{6, 6, 6, 6, 6, 6, 6, 6},
			}

			Expect(s.canvas).Should(Equal(expected))
		})
	})
})

func emptyUndoList() {
	for undoList.len() > 0 {
		_ = undoList.pop()
	}
}
