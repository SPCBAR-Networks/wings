package main

import (
	"github.com/sirupsen/logrus"
)

const (
	// Version of pterodactyld
	version = "0.0.1-alpha"
)

var (
	//logrus as a global var
	log = logrus.New()
)

func init() {
	//print start imagery
	printTitle()
	//initialize storage check and generation
	checkStorage()
	//Check port usage on startup and fail gracefully.
	checkPort()
}

func main() {
	log.Info("")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
