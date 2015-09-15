//+build ignore

package main

import (
	"log"

	"github.com/Luit/ami"
)

type QueueSummary struct {
	Queue string
}

type Originate struct {
	Channel string `ami:"empty"`

	Exten    string
	Context  string
	Priority string

	Application string
	Data        string

	Timeout    uint
	CallerID   string
	Variable   string
	Account    string
	EarlyMedia string
	Async      string
	Codecs     string
}

func main() {
	conn, err := ami.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Do(QueueSummary{}, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Do(Originate{}, nil)
	if err != nil {
		log.Fatal(err)
	}
}
