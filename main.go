package main

import (
	"fmt"

	"github.com/nunnatsa/piHatDraw/controller"
)

func main() {
	control := controller.NewController()
	<-control.Start()

	fmt.Println("\nGood Bye!")
}
