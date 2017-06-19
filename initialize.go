package main

import (
	"net"
	"os"
	"strconv"
)

func printTitle() {
	log.Info("+ ------------------------------------ +")
	log.Info("|  Running Pterodactyl Daemon " + version + "   |")
	log.Info("|  Copyright 2015 - 2017 Dane Everitt  |")
	log.Info("+ ------------------------------------ +")
	log.Info("Loading modules, this could take a few seconds.")
}

func checkStorage() {
	cinfo, err := os.Stat("./config")
	if err != nil {
		log.Errorln("Missing config directory generating now")
		os.Mkdir("./config", os.ModeDir)
		return
	}
	log.Info("Server config directories in place:", cinfo.IsDir())

	linfo, err := os.Stat("./logs")
	if err != nil {
		log.Error("Missing log directory generating now")
		os.Mkdir("./logs", os.ModeDir)
		return
	}
	log.Info("Log directory in place:", linfo.IsDir())

	pinfo, err := os.Stat("./packs")
	if err != nil {
		log.Error("Missing pack directories generating now")
		os.Mkdir("./packs", os.ModeDir)
		return
	}
	log.Info("Server pack directory in place:", pinfo.IsDir())
}

// Check if a port is available
func checkPort() (status bool, err error) {

	// Concatenate a colon and the port
	host := ":" + strconv.Itoa(5)

	// Try to create a server with the port
	server, err := net.Listen("tcp", host)

	// if it fails then the port is likely taken
	if err != nil {
		return false, err
	}

	// close the server
	server.Close()

	// we successfully used and closed the port
	// so it's now available to be used again
	return true, nil

}
