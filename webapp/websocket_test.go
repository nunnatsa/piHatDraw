package webapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/nunnatsa/piHatDraw/notifier"
)

func TestWebApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "webapp Suite")
}

var _ = Describe("Test the web application", func() {
	It("test register", func() {
		n := notifier.NewNotifier()
		defer n.Close()

		ce := make(chan ClientEvent)

		wa := NewWebApplication(n, 8080, ce)
		server := httptest.NewServer(wa.GetMux())
		defer server.Close()

		url := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/canvas/register"

		numClients := ClientEventRegistered(10)
		sockets := make([]*websocket.Conn, 0, numClients)
		message := "another message"
		defer func() {
			for _, ws := range sockets {
				ws.Close()
			}
		}()

		for i := ClientEventRegistered(1); i <= numClients; i++ {
			ws, _, err := websocket.DefaultDialer.Dial(url, nil)
			Expect(err).ToNot(HaveOccurred())

			sockets = append(sockets, ws)

			clientEvent := <-ce
			subscriberID, ok := clientEvent.(ClientEventRegistered)
			Expect(ok).To(BeTrue())
			Expect(subscriberID).Should(BeEquivalentTo(i))

			n.NotifyOne(uint64(subscriberID), []byte(message))

			_, p, err := ws.ReadMessage()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(p)).Should(Equal(message))
		}

		message = "hello there"
		n.NotifyAll([]byte(message))

		for _, ws := range sockets {
			_, p, err := ws.ReadMessage()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(p)).Should(Equal(message))
		}
	})

	Context("test set color request", func() {
		var (
			n      *notifier.Notifier
			ce     chan ClientEvent
			wa     *WebApplication
			server *httptest.Server
		)

		BeforeEach(func() {
			n = notifier.NewNotifier()
			ce = make(chan ClientEvent, 1)
			wa = NewWebApplication(n, 8080, ce)
			server = httptest.NewServer(wa.GetMux())
		})

		AfterEach(func() {
			n.Close()
			close(ce)
			server.Close()
		})

		DescribeTable("send valid POST request", func(url, reqBody string, expected interface{}) {
			url = server.URL + url

			res, err := server.Client().Post(url, "application/json", strings.NewReader(reqBody))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))

			Eventually(func() bool {
				clientEvent := <-ce
				Expect(clientEvent).Should(BeEquivalentTo(expected))
				return true
			}, time.Millisecond*10, time.Millisecond).Should(BeTrue())
		},
			Entry("test set color request", "/api/canvas/color", `{"color": "#123456"}`, 0x123456),
			Entry("test set tool request", "/api/canvas/tool", `{"toolName": "pen"}`, "pen"),
			Entry("test reset request", "/api/canvas/reset", `{"reset": true}`, true),
			Entry("test undo request", "/api/canvas/undo", `{"undo": true}`, true),
		)

		DescribeTable("should reject if not a POST request", func(url string) {
			url = server.URL + url

			res, err := server.Client().Get(url)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).Should(Equal(http.StatusMethodNotAllowed))
		},
			Entry("wrong method in set color request", "/api/canvas/color"),
			Entry("wrong method in set tool request", "/api/canvas/tool"),
			Entry("wrong method in reset request", "/api/canvas/reset"),
			Entry("wrong method in undo request", "/api/canvas/undo"),
		)

		DescribeTable("should reject if not the body is in wrong json format", func(url string) {
			url = server.URL + url

			res, err := server.Client().Post(url, "application/json", strings.NewReader(`bad json`))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).Should(Equal(http.StatusBadRequest))

		},
			Entry("wrong json in set color request", "/api/canvas/color"),
			Entry("wrong json in set tool request", "/api/canvas/tool"),
			Entry("wrong json in reset request", "/api/canvas/reset"),
			Entry("wrong json in undo request", "/api/canvas/undo"),
		)
	})

})
