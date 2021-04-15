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

func NewWebApplication(mailbox *notifier.Notifier, port uint16, ch chan<- ClientEvent) *WebApplication {
    mux := http.NewServeMux()
    ca := &WebApplication{mux: mux, notifier: mailbox, clientEvents: ch}
    mux.Handle("/", newIndexPage(port))
    mux.HandleFunc("/api/canvas/register", ca.register)
    mux.HandleFunc("/api/canvas/color", ca.setColor)
    mux.HandleFunc("/api/canvas/tool", ca.setTool)
    mux.HandleFunc("/api/canvas/reset", ca.reset)
    mux.HandleFunc("/api/canvas/download", ca.downloadImage)
    mux.HandleFunc("/api/canvas/undo", ca.undo)

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

func (ca WebApplication) downloadImage(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        r.ParseForm()

        pixelSizeStr := r.Form.Get("pixelSize")
        pixelSize, err := strconv.Atoi(pixelSizeStr)
        if err != nil || pixelSize < 1 || pixelSize > 20 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, `{"error": "wrong pixel size"}`)
            return
        }

        canvasChannel := make(chan [][]common.Color)
        ca.clientEvents <- ClientEventDownload(canvasChannel)
        imageData := <-canvasChannel

        imageCanvas, err := getImageCanvas(imageData, pixelSize)
        if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, `{"error": "%v"}`, err)
            return
        }

        w.Header().Add("Content-Disposition", `attachment; filename="untitled.png"`)
        w.Header().Set("Content-Type", "image/png")

        log.Println("downloading a file; pixel size =", pixelSize)
        png.Encode(w, imageCanvas)

    } else {
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

type undoRq struct {
    Undo bool `json:"undo"`
}

func (ca WebApplication) undo(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
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
    } else {
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
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
