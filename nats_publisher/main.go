package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	for {
		nc, err := nats.Connect("0.0.0.0:4222", nats.Name("Sender"))
		if err != nil {
			log.Panic("could not connect to nats:", err)
		}

		subject := "order"
		order := CreateFakeOrder(randBool(), randBool())

		b, err := json.Marshal(order)
		if err != nil {
			log.Println(err)
		}

		err = nc.Publish(subject, b)
		if err != nil {
			log.Println(err)
		}

		log.Println("Published:", order.OrderUID)
		time.Sleep(2 * time.Millisecond)
	}
}
