package webapp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/nunnatsa/piHatDraw/notifier"
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

func TestWebApplication_color(t *testing.T) {
	n := notifier.NewNotifier()
	ce := make(chan ClientEvent)

	wa := NewWebApplication(n, 8080, ce)
	server := httptest.NewServer(wa.GetMux())
	defer server.Close()

	url := server.URL + "/api/canvas/color"

	timeout := time.After(time.Millisecond * 10)

	go func() {
		select {
		case clientEvent := <-ce:
			if color, ok := clientEvent.(ClientEventSetColor); !ok {
				t.Errorf("should be ClientEventSetColor")
			} else {
				if uint32(color) != 0x123456 {
					t.Errorf("color should be 0x123456, but it's #%6x", color)
				}
			}
		case <-timeout:
			t.Fatal("Timeout")
		}
	}()

	res, err := server.Client().Post(url, "application/json", strings.NewReader(`{"color": "#123456"}`))
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Got non 200-OK status: %v", res.Status)
	}
	<-timeout

	res, err = server.Client().Get(url)
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("should be 405 status, but got %v instead", res.Status)
	}

	res, err = server.Client().Post(url, "application/json", strings.NewReader(`bad json`))
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("should be 400 status, but got %v instead", res.Status)
	}
}

func TestWebApplication_tools(t *testing.T) {
	n := notifier.NewNotifier()
	ce := make(chan ClientEvent)

	wa := NewWebApplication(n, 8080, ce)
	server := httptest.NewServer(wa.GetMux())
	defer server.Close()

	url := server.URL + "/api/canvas/tool"

	timeout := time.After(time.Millisecond * 10)

	go func() {
		select {
		case clientEvent := <-ce:
			if tool, ok := clientEvent.(ClientEventSetTool); !ok {
				t.Errorf("should be ClientEventSetTool")
			} else {
				if string(tool) != "pen" {
					t.Errorf("tool should be pen, but it's %v", tool)
				}
			}
		case <-timeout:
			t.Fatal("Timeout")
		}
	}()

	res, err := server.Client().Post(url, "application/json", strings.NewReader(`{"toolName": "pen"}`))
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Got non 200-OK status: %v", res.Status)
	}
	<-timeout

	res, err = server.Client().Get(url)
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("should be 405 status, but got %v instead", res.Status)
	}

	res, err = server.Client().Post(url, "application/json", strings.NewReader(`bad json`))
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("should be 400 status, but got %v instead", res.Status)
	}
}

func TestWebApplication_reset(t *testing.T) {
	n := notifier.NewNotifier()
	ce := make(chan ClientEvent)

	wa := NewWebApplication(n, 8080, ce)
	server := httptest.NewServer(wa.GetMux())
	defer server.Close()

	url := server.URL + "/api/canvas/reset"

	timeout := time.After(time.Millisecond * 10)

	go func() {
		select {
		case clientEvent := <-ce:
			if reset, ok := clientEvent.(ClientEventReset); !ok {
				t.Errorf("should be ClientEventReset")
			} else {
				if !reset {
					t.Errorf("reset should be true, but it's %v", reset)
				}
			}
		case <-timeout:
			t.Fatal("Timeout")
		}
	}()

	res, err := server.Client().Post(url, "application/json", strings.NewReader(`{"reset": true}`))
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Got non 200-OK status: %v", res.Status)
	}
	<-timeout

	res, err = server.Client().Get(url)
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("should be 405 status, but got %v instead", res.Status)
	}

	res, err = server.Client().Post(url, "application/json", strings.NewReader(`bad json`))
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("should be 400 status, but got %v instead", res.Status)
	}
}

func TestWebApplication_undo(t *testing.T) {
	n := notifier.NewNotifier()
	ce := make(chan ClientEvent)

	wa := NewWebApplication(n, 8080, ce)
	server := httptest.NewServer(wa.GetMux())
	defer server.Close()

	url := server.URL + "/api/canvas/undo"

	timeout := time.After(time.Millisecond * 10)

	go func() {
		select {
		case clientEvent := <-ce:
			if undo, ok := clientEvent.(ClientEventUndo); !ok {
				t.Errorf("should be ClientEventUndo")
			} else {
				if !undo {
					t.Errorf("undo should be true, but it's %v", undo)
				}
			}
		case <-timeout:
			t.Fatal("Timeout")
		}
	}()

	res, err := server.Client().Post(url, "application/json", strings.NewReader(`{"undo": true}`))
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Got non 200-OK status: %v", res.Status)
	}
	<-timeout

	res, err = server.Client().Get(url)
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("should be 405 status, but got %v instead", res.Status)
	}

	res, err = server.Client().Post(url, "application/json", strings.NewReader(`bad json`))
	if err != nil {
		t.Fatalf("Got error from the server: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("should be 400 status, but got %v instead", res.Status)
	}
}
