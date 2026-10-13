package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/gamechat"
)

func main() {
	client, err := gamechat.NewInfraiClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	address := os.Getenv("ADDR")
	if address == "" {
		address = ":8080"
	}
	fmt.Println("game chat backend listening on " + address)
	log.Fatal(http.ListenAndServe(address, gamechat.NewServer(client)))
}
