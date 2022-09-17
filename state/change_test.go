package state

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test change", func() {
	It("should be with length of 0 if it nil", func() {
		var s *changeStack
		Expect(s.len()).Should(BeZero())

		Expect(s.pop()).To(BeNil())
	})

	It("should be with length of 0 if it empty", func() {
		var s changeStack
		Expect(s.len()).Should(BeZero())
		Expect(s.pop()).To(BeNil())

		s = changeStack{}
		Expect(s.len()).Should(BeZero())
		Expect(s.pop()).To(BeNil())
	})

	It("should push to stack", func() {
		s := changeStack{}
		s.push(&Change{ToolName: "first"})
		Expect(s.len()).Should(Equal(1))
	})

	It("should be a LIFO collection", func() {
		s := changeStack{}

		s.push(&Change{ToolName: "first"})
		Expect(s.len()).Should(Equal(1))

		s.push(&Change{ToolName: "second"})
		Expect(s.len()).Should(Equal(2))

		s.push(&Change{ToolName: "third"})
		Expect(s.len()).Should(Equal(3))

		chng := s.pop()
		Expect(chng).ShouldNot(BeNil())
		Expect(chng.ToolName).Should(Equal("third"))
		Expect(s.len()).Should(Equal(2))

		chng = s.pop()
		Expect(chng).ShouldNot(BeNil())
		Expect(chng.ToolName).Should(Equal("second"))
		Expect(s.len()).Should(Equal(1))

		chng = s.pop()
		Expect(chng).ShouldNot(BeNil())
		Expect(chng.ToolName).Should(Equal("first"))
		Expect(s.len()).Should(BeZero())

		Expect(s.pop()).To(BeNil())
	})
})

func (s *changeStack) len() int {
	if s == nil {
		return 0
	}

	res := 0
	for p := s.head; p != nil; p = p.next {
		res++
	}

	return res
}
