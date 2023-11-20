package nats

import (
	"L0/database"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/nats-io/nats.go"
)

func NatsWriter() error {

	nc, err := nats.Connect("0.0.0.0:4222", nats.Name("Sender"))

	if err != nil {
		return fmt.Errorf("could not connect to nats: %v", err)
	}

	order := database.Order{}
	gofakeit.Struct(&order)

	subject := "test"

	b, err := json.Marshal(order)

	if err != nil {

		return err

	}
	for {
		err := nc.Publish(subject, b)
		if err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
	}

}
