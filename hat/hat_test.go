package hat

import (
	"os"
	"path"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
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

	Context("test findJoystickDeviceFile", func() {
		origFunc := getDevicesFilePath

		AfterEach(func() {
			getDevicesFilePath = origFunc
		})

		It("should return the the event file name from a valid file", func() {
			getDevicesFilePath = func() string {
				return path.Join(getTestFileLocation(), "validDeviceFile.txt")
			}

			eventFileName, err := findJoystickDeviceFile()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(eventFileName).To(Equal("/dev/input/event2"))
		})

		It("should return error if can't find the device in the file", func() {
			getDevicesFilePath = func() string {
				return path.Join(getTestFileLocation(), "notFoundDeviceFile.txt")
			}

			_, err := findJoystickDeviceFile()
			Expect(err).Should(HaveOccurred())
		})

		It("should return error if can't find the device name in the file", func() {
			getDevicesFilePath = func() string {
				return path.Join(getTestFileLocation(), "noDevNameDeviceFile.txt")
			}

			_, err := findJoystickDeviceFile()
			Expect(err).Should(HaveOccurred())
		})

		It("should return error if can't find the event name in the file", func() {
			getDevicesFilePath = func() string {
				return path.Join(getTestFileLocation(), "noEventNameDeviceFile.txt")
			}

			_, err := findJoystickDeviceFile()
			Expect(err).Should(HaveOccurred())
		})

		It("should return error if can't find the devices file", func() {
			getDevicesFilePath = func() string {
				return path.Join(getTestFileLocation(), "notExistsFile")
			}

			_, err := findJoystickDeviceFile()
			Expect(err).Should(HaveOccurred())
			Expect(os.IsNotExist(err)).Should(BeTrue())
		})
	})
})

func getTestFileLocation() string {
	wd, err := os.Getwd()
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	if strings.HasSuffix(wd, "hat") {
		return path.Join(wd, "testFiles")
	}
	return path.Join(wd, "hat", "testFiles")
}
