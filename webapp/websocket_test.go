package webapp

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    neturl "net/url"
    "strings"
    "testing"
    "time"

    "github.com/gorilla/websocket"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/ginkgo/extensions/table"
    . "github.com/onsi/gomega"

    "github.com/nunnatsa/piHatDraw/common"
    "github.com/nunnatsa/piHatDraw/notifier"
)

func TestWebApp(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "webapp Suite")
}

var _ = Describe("Test the web application", func() {
    Context("test register request", func() {

        var (
            n      *notifier.Notifier
            server *httptest.Server
            wa     *WebApplication
            ce     chan ClientEvent
        )

        BeforeEach(func() {
            n = notifier.NewNotifier()
            ce = make(chan ClientEvent)
            wa = NewWebApplication(n, ce)
            server = httptest.NewServer(wa.GetMux())
        })

        AfterEach(func() {
            n.Close()
            server.Close()
            close(ce)
        })

        It("test register", func() {
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

        It("should reject if the method is wrong", func() {
            url := server.URL + "/api/canvas/register"

            res, err := server.Client().Post(url, "text/plain", strings.NewReader("test request body"))
            Expect(err).ShouldNot(HaveOccurred())
            Expect(res.StatusCode).Should(Equal(http.StatusMethodNotAllowed))
        })
    })

    Context("test POST requests", func() {
        var (
            n      *notifier.Notifier
            ce     chan ClientEvent
            wa     *WebApplication
            server *httptest.Server
        )

        BeforeEach(func() {
            n = notifier.NewNotifier()
            ce = make(chan ClientEvent, 1)
            wa = NewWebApplication(n, ce)
            server = httptest.NewServer(wa.GetMux())
        })

        AfterEach(func() {
            n.Close()
            close(ce)
            server.Close()
        })

        DescribeTable("send valid POST request", func(url, reqBody string, expected any) {
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

    Context("test download request", func() {
        var (
            n      *notifier.Notifier
            ce     chan ClientEvent
            wa     *WebApplication
            server *httptest.Server
        )

        BeforeEach(func() {
            n = notifier.NewNotifier()
            ce = make(chan ClientEvent, 1)
            wa = NewWebApplication(n, ce)
            server = httptest.NewServer(wa.GetMux())
        })

        AfterEach(func() {
            n.Close()
            close(ce)
            server.Close()
        })

        DescribeTable("should download an image file", func(fileName string) {
            done := make(chan bool)

            go sendImageData(ce, done)

            url := server.URL + "/api/canvas/download?pixelSize=3"
            if len(fileName) > 0 {
                url = fmt.Sprintf("%s&fileName=%s.png", url, fileName)
            } else {
                fileName = "untitled"
            }

            res, err := server.Client().Get(url)
            Expect(err).ToNot(HaveOccurred())
            Expect(res.StatusCode).Should(Equal(http.StatusOK))

            Eventually(func() bool { return <-done }, time.Millisecond*10, time.Millisecond).Should(BeTrue())

            Eventually(func() bool {
                defer res.Body.Close()

                body, err := ioutil.ReadAll(res.Body)
                Expect(err).ToNot(HaveOccurred())
                Expect(body).ToNot(BeEmpty())
                fileType := string(body[:4])
                Expect(fileType).Should(Equal("\x89PNG"))

                Expect(res.Header.Get("Content-Type")).Should(Equal("image/png"))
                expected := fmt.Sprintf(`attachment; filename="%s.png"`, fileName)
                Expect(res.Header.Get("Content-Disposition")).Should(Equal(expected))

                return true
            }).Should(BeTrue())
        },
            Entry("with file name", "test"),
            Entry("without file name", ""),
        )

        DescribeTable("should return error if the data is empty", func(data [][]common.Color) {
            done := make(chan bool)

            go func() {
                defer close(done)
                clientEvent := <-ce

                cb, ok := clientEvent.(ClientEventDownload)
                if !ok {
                    done <- false
                    return
                }
                Expect(ok).Should(BeTrue())

                cb <- data
                done <- true
            }()

            url := server.URL + "/api/canvas/download?pixelSize=3"
            res, err := server.Client().Get(url)

            Eventually(func() bool { return <-done }, time.Millisecond*10, time.Millisecond).Should(BeTrue())

            Expect(err).ToNot(HaveOccurred())
            Expect(res.StatusCode).Should(Equal(http.StatusInternalServerError))

            Eventually(func() bool {
                defer res.Body.Close()

                Expect(res.Header.Get("Content-Type")).Should(Equal("application/json"))

                dec := json.NewDecoder(res.Body)
                errMsg := &errorResponse{}
                err = dec.Decode(errMsg)
                Expect(err).ToNot(HaveOccurred())
                Expect(errMsg.Error).Should(Equal("can't get the data"))

                return true
            }).Should(BeTrue())
        },
            Entry("no data received", nil),
            Entry("empty data received", [][]common.Color{}),
            Entry("empty lines received", [][]common.Color{{}, {}, {}, {}}),
        )

        It("should return error if there is no pixelSize query parameter", func() {
            url := server.URL + "/api/canvas/download"
            res, err := server.Client().Get(url)

            Consistently(ce).ShouldNot(Receive())

            Eventually(func() bool {
                defer res.Body.Close()
                Expect(err).ToNot(HaveOccurred())
                Expect(res.StatusCode).Should(Equal(http.StatusBadRequest))

                Expect(res.Header.Get("Content-Type")).Should(Equal("application/json"))

                dec := json.NewDecoder(res.Body)
                errMsg := &errorResponse{}
                err = dec.Decode(errMsg)
                Expect(err).ToNot(HaveOccurred())
                Expect(errMsg.Error).Should(Equal("wrong pixel size"))

                return true
            }).Should(BeTrue())
        })

        It("should return error if there is the pixelSize is too big", func() {
            url := server.URL + "/api/canvas/download?pixelSize=50"
            res, err := server.Client().Get(url)

            Consistently(ce).ShouldNot(Receive())

            Eventually(func() bool {
                defer res.Body.Close()
                Expect(err).ToNot(HaveOccurred())
                Expect(res.StatusCode).Should(Equal(http.StatusBadRequest))

                Expect(res.Header.Get("Content-Type")).Should(Equal("application/json"))

                dec := json.NewDecoder(res.Body)
                errMsg := &errorResponse{}
                err = dec.Decode(errMsg)
                Expect(err).ToNot(HaveOccurred())
                Expect(errMsg.Error).Should(Equal("wrong pixel size"))

                return true
            }).Should(BeTrue())
        })

        It("should return error if the method is wrong", func() {
            url := server.URL + "/api/canvas/download"
            form := neturl.Values{
                "pixelSize": []string{"3"},
                "fileName":  []string{"test.png"},
            }
            res, err := server.Client().PostForm(url, form)
            Expect(err).ToNot(HaveOccurred())
            Expect(res.StatusCode).Should(Equal(http.StatusMethodNotAllowed))

            Consistently(ce).ShouldNot(Receive())
        })
    })

})

type errorResponse struct {
    Error string `json:"error,omitempty"`
}

func sendImageData(ce <-chan ClientEvent, done chan<- bool) {
    defer close(done)
    clientEvent := <-ce

    cb, ok := clientEvent.(ClientEventDownload)
    if !ok {
        done <- false
        return
    }
    Expect(ok).Should(BeTrue())

    cb <- [][]common.Color{
        {0, 0, 0, 0, 0, 0, 0, 0},
        {0, 0xFFFFFF, 0xFFFFFF, 0xFFFFFF, 0xFFFFFF, 0xFFFFFF, 0xFFFFFF, 0},
        {0, 0xFFFFFF, 0, 0, 0, 0, 0xFFFFFF, 0},
        {0, 0xFFFFFF, 0, 0, 0, 0, 0xFFFFFF, 0},
        {0, 0xFFFFFF, 0, 0, 0, 0, 0xFFFFFF, 0},
        {0, 0xFFFFFF, 0, 0, 0, 0, 0xFFFFFF, 0},
        {0, 0xFFFFFF, 0xFFFFFF, 0xFFFFFF, 0xFFFFFF, 0xFFFFFF, 0xFFFFFF, 0},
        {0, 0, 0, 0, 0, 0, 0, 0},
    }
    done <- true
}
