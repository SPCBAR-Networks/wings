package main

import (
	"fmt"
)

func init() {
	storageCheck()
}

func main() {

	fmt.Println("The Wings Daemon is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
