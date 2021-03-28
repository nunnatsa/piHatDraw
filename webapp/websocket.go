package webapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/notifier"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1500,
		WriteBufferSize: 1500,
	}
)

type ClientEvent interface{}

type ClientEventRegistered uint64

type ClientEventSetColor common.Color

type ClientEventSetTool string

type ClientEventReset bool

type WebApplication struct {
	mux          *http.ServeMux
	notifier     *notifier.Notifier
	clientEvents chan<- ClientEvent
}

func (ca WebApplication) GetMux() *http.ServeMux {
	return ca.mux
}

func NewWebApplication(mailbox *notifier.Notifier, port uint16, ch chan<- ClientEvent) *WebApplication {
	mux := http.NewServeMux()
	ca := &WebApplication{mux: mux, notifier: mailbox, clientEvents: ch}
	mux.Handle("/", newIndexPage(port))
	mux.HandleFunc("/api/canvas/register", ca.register)
	mux.HandleFunc("/api/canvas/color", ca.setColor)
	mux.HandleFunc("/api/canvas/tool", ca.setTool)
	mux.HandleFunc("/api/canvas/reset", ca.reset)

	return ca
}

func (ca WebApplication) register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		if err != nil {
			log.Println("Error:", err)
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))

			return
		}

		defer conn.Close()

		subscription := make(chan []byte)

		id := ca.notifier.Subscribe(subscription)
		defer ca.notifier.Unsubscribe(id)
		ca.clientEvents <- ClientEventRegistered(id)

		for js := range subscription {
			log.Printf("got event; updating client %d\n", id)
			if err := conn.WriteMessage(websocket.TextMessage, js); err != nil {
				log.Printf("failed to send message to the client %d: %v\n", id, err)
				return
			}
		}
		log.Println("Connection is closed")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type setColorRq struct {
	Color common.Color `json:"color"`
}

func (ca WebApplication) setColor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		enc := json.NewDecoder(r.Body)
		msg := &setColorRq{}
		err := enc.Decode(msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "can't parse json'"}`)
			return
		}

		log.Printf("Got set color request. Color = #%06x", msg.Color)

		clientEvent := ClientEventSetColor(msg.Color)
		ca.clientEvents <- clientEvent
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type setToolRq struct {
	ToolName string `json:"toolName"`
}

func (ca WebApplication) setTool(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		enc := json.NewDecoder(r.Body)
		msg := &setToolRq{}
		err := enc.Decode(msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "can't parse json'"}`)
			return
		}

		log.Printf("Got set tool request. tool name = %v", msg.ToolName)

		clientEvent := ClientEventSetTool(msg.ToolName)
		ca.clientEvents <- clientEvent
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type resetRq struct {
	Reset bool `json:"reset"`
}

func (ca WebApplication) reset(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		enc := json.NewDecoder(r.Body)
		msg := &resetRq{}
		err := enc.Decode(msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error": "can't parse json'"}`)
			return
		}

		if msg.Reset {
			log.Printf("Got reset request")
		}

		clientEvent := ClientEventReset(true)
		ca.clientEvents <- clientEvent
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
