package main

import (
	"github.com/6b70/peerbeam/cmd"
	"log"
)

func main() {
	err := cmd.App()
	if err != nil {
		log.Fatal(err)
	}

}
