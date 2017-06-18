package main

import (
	"fmt"
	"log"
	"os"
)

func printTitle() {
	log.Println("+ ------------------------------------ +")
	log.Println("|  Running Pterodactyl Daemon Canary   |")
	log.Println("|  Copyright 2015 - 2017 Dane Everitt  |")
	log.Println("+ ------------------------------------ +")
	log.Println("Loading modules, this could take a few seconds.")
}

func storageCheck() {
	cinfo, err := os.Stat("./config")
	if err != nil {
		fmt.Println("Missing config directories generating now")
		os.Mkdir("./config", os.ModeDir)
		return
	}
	fmt.Println("Server config directories in place:", cinfo.IsDir())

	linfo, err := os.Stat("./logs")
	if err != nil {
		fmt.Println("Missing log directory generating now")
		os.Mkdir("./logs", os.ModeDir)
		return
	}
	fmt.Println("Log directory in place:", linfo.IsDir())

	pinfo, err := os.Stat("./packs")
	if err != nil {
		fmt.Println("Missing pack directories generating now")
		os.Mkdir("./packs", os.ModeDir)
		return
	}
	fmt.Println("Server pack directory in place:", pinfo.IsDir())
}
