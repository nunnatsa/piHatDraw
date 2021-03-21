package webapp

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

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
