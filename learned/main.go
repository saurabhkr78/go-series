package main

import (
	"github.com/saurabhkr78/learned/config"
	"log"
)

func main() {
	if err := config.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	defer config.Close()

}
