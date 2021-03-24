package webapp

import (
	"github.com/gorilla/websocket"
	"github.com/nunnatsa/piHatDraw/notifier"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebApplication_register(t *testing.T) {
	n := notifier.NewNotifier()
	ce := make(chan ClientEvent)

	wa := NewWebApplication(n, 8080, ce)
	server := httptest.NewServer(wa.GetMux())
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/canvas/register"

	numClients := ClientEventRegistered(10)
	sockets := make([]*websocket.Conn, 0, numClients)
	message := "another message"

	for i := ClientEventRegistered(1); i <= numClients; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		defer ws.Close()

		sockets = append(sockets, ws)

		clientEvent := <-ce
		subscriberID, ok := clientEvent.(ClientEventRegistered)
		if !ok {
			t.Fatal("wrong client event type")
		}
		if subscriberID != i {
			t.Fatalf("wrong subscriberId; shoud be %d but it's %d", i, subscriberID)
		}

		n.NotifyOne(uint64(subscriberID), []byte(message))

		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if string(p) != message {
			t.Fatalf("bad message")
		}
	}

	message = "hello there"
	n.NotifyAll([]byte(message))

	for _, ws := range sockets {
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if string(p) != message {
			t.Fatalf("bad message")
		}
	}
}
