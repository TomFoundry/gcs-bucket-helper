package main

import (
	"log"

	"github.com/athera-io/gcs-bucket-helper/internal/server"
)

var (
	// appID is set at build time
	appID string
	// appPseudoSecret is set at build time
	appPseudoSecret string
)

func main() {

	serv := server.New("8001", appID, appPseudoSecret)

	if err := serv.Serve(); err != nil {
		log.Fatal("HTTP Server has failed: ", err)
	}
}
