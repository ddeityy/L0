package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func SendFakeData(interval time.Duration) error {

	nc, err := nats.Connect("0.0.0.0:4222", nats.Name("Sender"))

	if err != nil {
		return fmt.Errorf("could not connect to nats: %v", err)
	}
	subject := "test"

	for {
		order := CreateFakeOrder(randBool(), randBool(), randBool())

		b, err := json.Marshal(order)

		if err != nil {
			return err
		}

		err = nc.Publish(subject, b)
		if err != nil {
			return err
		}
		time.Sleep(interval)
	}

}
