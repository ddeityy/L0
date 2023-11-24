package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func SendFakeData() error {
	nc, err := nats.Connect("0.0.0.0:4222", nats.Name("Sender"))
	if err != nil {
		log.Println(err)
		return fmt.Errorf("could not connect to nats: %v", err)
	}
	subject := "order"
	badItems := randBool()
	badPayment := randBool()
	order := CreateFakeOrder(badPayment, badItems)

	b, err := json.Marshal(order)

	if err != nil {
		return err
	}

	err = nc.Publish(subject, b)
	if err != nil {
		return err
	}
	log.Println("Published:", order.OrderUID, "Bad items:", badItems, "Bad payment:", badPayment)
	return nil
}
