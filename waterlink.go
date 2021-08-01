package main

import (
	"context"
	"fmt"
	"github.com/lukasl-dev/waterlink"
	"log"
	"net/url"
	"os"
)

var (
	userID     = os.Getenv("USERID")
	passphrase = os.Getenv("PASSPHRASE")
	host       = os.Getenv("HOST")

	httpHost, _ = url.Parse(fmt.Sprintf("http://%s", host))
	wsHost, _   = url.Parse(fmt.Sprintf("ws://%s", host))

	connOpts = waterlink.NewConnectOptions().WithUserID(userID).WithPassphrase(passphrase)
	reqOpts  = waterlink.NewRequesterOptions().WithPassphrase(passphrase)

	conn waterlink.Connection
	req  waterlink.Requester
)

// initWaterlink initializes waterlink.
func initWaterlink() {
	openConnection()
	createRequester()
}

// openConnection opens a connection to the server.
func openConnection() {
	var err error
	conn, err = waterlink.Connect(context.TODO(), *wsHost, connOpts)
	if err != nil {
		log.Fatalln("Opening connection failed:", err)
	}
	log.Println("Connection established.")
}

// createRequester creates a new requester.
func createRequester() {
	req = waterlink.NewRequester(*httpHost, reqOpts)
	log.Println("Requester created.")
}
