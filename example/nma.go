package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dustin/nma.go"
)

var apikey, event, url string
var priority int

func init() {
	flag.StringVar(&apikey, "apikey", "", "Your API key")
	flag.StringVar(&event, "event", "", "NMA event")
	flag.IntVar(&priority, "pri", 0, "NMA priority (-2 to 2)")
}

func main() {
	flag.Parse()

	n := nma.New(apikey)

	e := nma.Notification{
		Description: strings.Join(flag.Args(), " "),
		Event:       event,
		Priority:    priority,
	}

	if err := n.Notify(&e); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message:  %v\n", err)
		os.Exit(1)
	}
}
