package main

import "github.com/nunnatsa/piHatDraw/controller"

func main() {
	control := controller.NewController()
	<-control.Start()
}
