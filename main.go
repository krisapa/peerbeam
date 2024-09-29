package main

import (
	"github.com/6b70/peerbeam/cmd"
	"golang.design/x/clipboard"
	"log"
)

func main() {
	err := clipboard.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.App()
	if err != nil {
		log.Fatal(err)
	}

}
