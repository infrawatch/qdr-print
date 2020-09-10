package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"qpid.apache.org/amqp"
	"qpid.apache.org/electron"
)

func main() {
	connstr := flag.String("connection", "amqp://127.0.0.1:5666/anycast/ceilometer/metering.sample", "QDR URL")
	flag.Parse()

	container := electron.NewContainer("qdr-print")
	url, err := amqp.ParseURL(*connstr)
	if err != nil {
		fmt.Printf("[ERROR] Failed parsing URL: %s\n", err)
		os.Exit(1)
	}

	conn, err := container.Dial("tcp", url.Host)
	if err != nil {
		fmt.Printf("[ERROR] Failed to connect: %s\n", err)
		os.Exit(1)
	}
	defer conn.Close(err)

	addr := strings.TrimPrefix(url.Path, "/")
	opts := []electron.LinkOption{electron.Source(addr)}

	receiver, err := conn.Receiver(opts...)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create receiver: %s\n", err)
		os.Exit(1)
	}

	for {
		if msg, err := receiver.Receive(); err == nil {
			msg.Accept()
			prettyJSON, err := json.MarshalIndent()
			fmt.Printf("%v\n", msg.Message)
		} else {
			fmt.Printf("[ERROR] Failed to receive message: %s\n", err)
		}
	}
}
