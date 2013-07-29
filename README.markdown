# Notify My Android on the Go

This is a [go][go] client for [Notify my Android][NMA].

With this, you can send simple notifications directly to your phone
and other android devices quickly and easily.

## Installation

`go get github.com/dustin/go-nma`

## Usage

```go
package main

import "github.com/dustin/go-nma"

func main() {
    n := nma.New("yourapikey")
    e := nma.Notification{
      Event: "It worked!",
      Description: "I was able to send a message!",
      Priority: 1,
    }

	if err := n.Notify(&e); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message:  %v\n", err)
		os.Exit(1)
	}
}
```

[go]: http://golang.org/
[NMA]: http://www.notifymyandroid.com/
