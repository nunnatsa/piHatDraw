package common

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCommon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "common Suite")
}

var _ = Describe("Test the common package", func() {
	Context("test color.MarshalJSON", func() {
		It("should decode 0", func() {
			c := Color(0)
			var buf bytes.Buffer
			Expect(json.NewEncoder(&buf).Encode(c)).ToNot(HaveOccurred())

			Expect(buf.String()).Should(Equal("\"#000000\"\n"))
		})

		It("should decode 0xFFFFFF", func() {
			c := Color(0xFFFFFF)
			var buf bytes.Buffer
			Expect(json.NewEncoder(&buf).Encode(c)).ToNot(HaveOccurred())

			Expect(buf.String()).Should(Equal("\"#ffffff\"\n"))
		})

		It("should ignore 8 MSB", func() {
			c := Color(0xFFFFFFFF)
			var buf bytes.Buffer
			Expect(json.NewEncoder(&buf).Encode(c)).ToNot(HaveOccurred())

			Expect(buf.String()).Should(Equal("\"#ffffff\"\n"))
		})

		It("should encode an array of colors", func() {
			var buf bytes.Buffer
			cs := []Color{0xFF0000, 0x00FF00, 0x0000FF}

			Expect(json.NewEncoder(&buf).Encode(cs)).ToNot(HaveOccurred())
			Expect(buf.String()).Should(Equal("[\"#ff0000\",\"#00ff00\",\"#0000ff\"]\n"))
		})
	})

	Context("test color.UnmarshalJSON", func() {
		It("should decode #000000", func() {
			var c Color
			Expect(json.Unmarshal([]byte(`"#000000"`), &c)).ToNot(HaveOccurred())
			Expect(c).Should(BeEquivalentTo(0))
		})

		It("should decode #ffffff", func() {
			var c Color
			Expect(json.Unmarshal([]byte(`"#ffffff"`), &c)).ToNot(HaveOccurred())
			Expect(c).Should(BeEquivalentTo(0xFFFFFF))
		})

		It("should decode an array of colors", func() {
			var c []Color
			Expect(json.Unmarshal([]byte(`["#ff0000", "#00ff00", "#0000ff"]`), &c)).ToNot(HaveOccurred())
			Expect(c).Should(Equal([]Color{0xFF0000, 0x00FF00, 0x0000FF}))
		})

		It("should return error if the json format is wrong", func() {
			var c *Color
			Expect(json.Unmarshal([]byte(`"#ff0000`), c)).To(HaveOccurred())
			Expect(c).Should(BeNil())
		})

		It("should return error if the HTML format is wrong", func() {
			var cs []Color
			Expect(json.Unmarshal([]byte(`["ff0000"]`), &cs)).To(HaveOccurred())

			var c *Color
			err := json.Unmarshal([]byte(`"#ff00"`), c)
			Expect(err).To(HaveOccurred())
			Expect(c).Should(BeNil())
		})
	})
})
