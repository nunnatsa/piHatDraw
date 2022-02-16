package webapp

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/nunnatsa/piHatDraw/common"
	"github.com/nunnatsa/piHatDraw/notifier"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1500,
		WriteBufferSize: 1500,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type ClientEvent interface{}

type ClientEventRegistered uint64

type ClientEventSetColor common.Color

type ClientEventSetTool string

type ClientEventReset bool

type ClientEventDownload chan [][]common.Color

type ClientEventUndo bool

type WebApplication struct {
	mux          *http.ServeMux
	notifier     *notifier.Notifier
	clientEvents chan<- ClientEvent
}

func (ca WebApplication) GetMux() *http.ServeMux {
	return ca.mux
}

func NewWebApplication(mailbox *notifier.Notifier, ch chan<- ClientEvent) *WebApplication {
	mux := http.NewServeMux()
	ca := &WebApplication{mux: mux, notifier: mailbox, clientEvents: ch}

	mux.Handle("/", ui)
	mux.Handle("/api/canvas/register", GetOnlyRequest(ca.register))
	mux.Handle("/api/canvas/color", PostOnlyRequest(ca.setColor))
	mux.Handle("/api/canvas/tool", PostOnlyRequest(ca.setTool))
	mux.Handle("/api/canvas/reset", PostOnlyRequest(ca.reset))
	mux.Handle("/api/canvas/download", GetOnlyRequest(ca.downloadImage))
	mux.Handle("/api/canvas/undo", PostOnlyRequest(ca.undo))

	return ca
}

func (ca WebApplication) register(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		log.Println("Error:", err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	defer conn.Close()

	subscription := make(chan []byte, 1)

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

	log.Printf("Connection %d is closed\n", id)
}

type setColorRq struct {
	Color common.Color `json:"color"`
}

func (ca WebApplication) setColor(w http.ResponseWriter, r *http.Request) {
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
}

type setToolRq struct {
	ToolName string `json:"toolName"`
}

func (ca WebApplication) setTool(w http.ResponseWriter, r *http.Request) {
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
}

type resetRq struct {
	Reset bool `json:"reset"`
}

func (ca WebApplication) reset(w http.ResponseWriter, r *http.Request) {
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
}

func (ca WebApplication) downloadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error": "wrong request"}`)
		return
	}

	pixelSizeStr := r.Form.Get("pixelSize")
	pixelSize, err := strconv.Atoi(pixelSizeStr)

	if err != nil || pixelSize < 1 || pixelSize > 20 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error": "wrong pixel size"}`)
		return
	}

	fileName := r.Form.Get("fileName")
	if fileName == "" {
		fileName = "untitled.png"
	}

	canvasChannel := make(chan [][]common.Color, 1)
	defer close(canvasChannel)
	ca.clientEvents <- ClientEventDownload(canvasChannel)
	imageData := <-canvasChannel

	imageCanvas, err := getImageCanvas(imageData, pixelSize)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, err)
		return
	}

	w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	w.Header().Set("Content-Type", "image/png")

	log.Printf("downloading a file %s; pixel size = %d\n", fileName, pixelSize)
	_ = png.Encode(w, imageCanvas)
}

type undoRq struct {
	Undo bool `json:"undo"`
}

func (ca WebApplication) undo(w http.ResponseWriter, r *http.Request) {
	enc := json.NewDecoder(r.Body)
	msg := &undoRq{}
	err := enc.Decode(msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "can't undo'"}`)
		return
	}

	if msg.Undo {
		log.Printf("Got undo request")
	}

	clientEvent := ClientEventUndo(true)
	ca.clientEvents <- clientEvent
}

func getImageCanvas(imageData [][]common.Color, pixelSize int) (*image.RGBA, error) {
	height := len(imageData) * pixelSize
	if height == 0 {
		return nil, fmt.Errorf("can't get the data")
	}

	width := len(imageData[0]) * pixelSize
	if width == 0 {
		return nil, fmt.Errorf("can't get the data")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y, line := range imageData {
		for x, pixel := range line {
			setPixel(img, x, y, toColor(pixel), pixelSize)
		}
	}

	return img, nil
}

func setPixel(img *image.RGBA, x int, y int, pixel color.Color, pixelSize int) {
	x = x * pixelSize
	y = y * pixelSize
	for x1 := x; x1 < x+pixelSize; x1++ {
		for y1 := y; y1 < y+pixelSize; y1++ {
			img.Set(x1, y1, pixel)
		}
	}
}

func toColor(pixel common.Color) color.Color {
	r := uint8((pixel >> 16) & 0xFF)
	g := uint8((pixel >> 8) & 0xFF)
	b := uint8(pixel & 0xFF)

	return color.RGBA{A: 0xFF, R: r, G: g, B: b}
}

func GetOnlyRequest(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		next(w, r)
	})
}

func PostOnlyRequest(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		next(w, r)
	})
}
