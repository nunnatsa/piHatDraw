package hat

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "hat Suite")
}

var _ = Describe("test the hat package", func() {
	Context("test toHatColor", func() {
		It("should convert 0xFFFFFF", func() {
			Expect(toHatColor(0xFFFFFF)).Should(BeEquivalentTo(0xFFFF))
		})

		It("should ignore LSB", func() {
			Expect(toHatColor(0b000001110000000000000000)).Should(BeEquivalentTo(0))
			Expect(toHatColor(0b000000000000001100000000)).Should(BeEquivalentTo(0))
			Expect(toHatColor(0b000000000000000000000111)).Should(BeEquivalentTo(0))
			Expect(toHatColor(0b000001110000001100000111)).Should(BeEquivalentTo(0))
		})

		It("should convert red", func() {
			Expect(toHatColor(0b111110000000000000000000)).Should(BeEquivalentTo(0xF800))
			Expect(toHatColor(0b111111110000000000000000)).Should(BeEquivalentTo(0xF800))
		})

		It("should convert green", func() {
			Expect(toHatColor(0b000000001111110000000000)).Should(BeEquivalentTo(0x07E0))
			Expect(toHatColor(0b000000001111111100000000)).Should(BeEquivalentTo(0x07E0))
		})

		It("should convert blue", func() {
			Expect(toHatColor(0b000000000000000011111000)).Should(BeEquivalentTo(0x01F))
			Expect(toHatColor(0b000000000000000011111111)).Should(BeEquivalentTo(0x01F))
		})
	})
})
