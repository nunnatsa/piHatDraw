package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nunnatsa/piHatDraw/notifier"
	"github.com/nunnatsa/piHatDraw/webapp"

	"github.com/nunnatsa/piHatDraw/controller"
)

var (
	canvasWidth, canvasHeight uint8
	port                      uint16
)

func init() {
	var width, height, prt uint
	flag.UintVar(&width, "width", 24, "Canvas width in pixels")
	flag.UintVar(&height, "height", 24, "Canvas height in pixels")
	flag.UintVar(&prt, "port", 8080, "The application port")

	flag.Parse()

	if width < 8 {
		fmt.Println("The minimum width of the canvas is 8 pixels; setting it for you")
		width = 8
	}

	if width > 40 {
		log.Fatal("ERROR: The maximum width of the canvas is 40 pixels")
	}
	canvasWidth = uint8(width)

	if height < 8 {
		fmt.Println("The minimum height of the canvas is 8 pixels; setting it for you")
		height = 8
	}

	if height > 40 {
		log.Fatal("ERROR: The maximum height of the canvas is 40 pixels")
	}
	canvasHeight = uint8(height)

	port = uint16(prt)

	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("In your web browser, go to http://%s:%d\n", hostname, port)
}

func main() {
	n := notifier.NewNotifier()

	clientEvents := make(chan webapp.ClientEvent)
	webApplication := webapp.NewWebApplication(n, port, clientEvents)

	portStr := fmt.Sprintf(":%d", port)
	server := http.Server{Addr: portStr, Handler: webApplication.GetMux()}

	control := controller.NewController(n, clientEvents, canvasWidth, canvasHeight)
	done := control.Start()
	go func() {
		<-done
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown the server; %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Panic(err)
	} else {
		fmt.Println("\nGood Bye!")
	}

}
