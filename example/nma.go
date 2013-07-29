package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dustin/go-nma"
)

var apikey, event, application string
var priority int

func init() {
	flag.StringVar(&apikey, "apikey", "", "Your API key")
	flag.StringVar(&event, "event", "", "NMA event")
	flag.StringVar(&application, "application", "", "Notifying application")
	flag.IntVar(&priority, "pri", 0, "NMA priority (-2 to 2)")
}

func main() {
	flag.Parse()

	n := nma.New(apikey)

	e := nma.Notification{
		Description: strings.Join(flag.Args(), " "),
		Application: application,
		Event:       event,
		Priority:    nma.PriorityLevel(priority),
	}

	if err := n.Notify(&e); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message:  %v\n", err)
		os.Exit(1)
	}
}
