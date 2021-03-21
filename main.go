package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nunnatsa/piHatDraw/notifier"
	"github.com/nunnatsa/piHatDraw/webapp"

	"github.com/nunnatsa/piHatDraw/controller"
)

func main() {

	n := notifier.NewNotifier()

	clientEvents := make(chan webapp.ClientEvent)
	webApplication := webapp.NewWebApplication(n, 8080, clientEvents)

	server := http.Server{Addr: ":8080", Handler: webApplication.GetMux()}

	control := controller.NewController(n, clientEvents)
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
