package notifier

import (
	"fmt"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNotifier(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "notifier Suite")
}

var _ = Describe("test notifier", func() {

	var notifier *Notifier
	BeforeEach(func() {
		notifier = NewNotifier()
	})

	AfterEach(func() {
		for id := range notifier.clientMap {
			notifier.Unsubscribe(id)
		}

	})

	It("should register to the notifier", func() {
		const numSubscribers = 10
		wg := &sync.WaitGroup{}

		done := make(chan bool)
		wg.Add(numSubscribers)
		go func() {
			wg.Wait()
			close(done)
		}()

		for i := 0; i < numSubscribers; i++ {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				ch := make(chan []byte)
				notifier.Subscribe(ch)
			}(wg)
		}

		<-done

		Expect(notifier.clientMap).Should(HaveLen(numSubscribers), fmt.Sprintf("number of subscribers should be %d", numSubscribers))

	})

	It("should unsubscribe", func() {
		ch := make(chan []byte)

		id := notifier.Subscribe(ch)

		notifier.Unsubscribe(id)

		Eventually(ch).Should(BeClosed())

		Expect(notifier.clientMap).To(HaveLen(0))
	})

	It("should notify all", func() {
		const numSubscribers = 10

		channels := make([]chan []byte, numSubscribers)

		for i := 0; i < numSubscribers; i++ {
			ch := make(chan []byte, 1)
			channels[i] = ch

			notifier.Subscribe(ch)
		}

		Expect(notifier.clientMap).Should(HaveLen(numSubscribers), fmt.Sprintf("Number of subscribers should be %d", numSubscribers))

		notifier.NotifyAll([]byte("message"))

		for i := 0; i < numSubscribers; i++ {
			Eventually(<-channels[i]).Should(BeEquivalentTo("message"))
		}
	})

	It("should notify one", func() {
		ch1 := make(chan []byte, 1)
		ch2 := make(chan []byte, 1)
		id1 := notifier.Subscribe(ch1)
		_ = notifier.Subscribe(ch2)

		notifier.NotifyOne(id1, []byte("message"))
		Eventually(<-ch1).Should(BeEquivalentTo("message"))
		Consistently(ch2).ShouldNot(Receive())
	})
})
