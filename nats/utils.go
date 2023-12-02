package nats

import (
	"L0/database"
	"errors"

	"github.com/nats-io/nats.go"
)

func ValidateOrder(order database.CacheOrder) error {
	if order.Payment.Transaction != order.OrderUID {
		return errors.New("payment does not match")
	}

	for _, item := range order.Items {
		if item.TrackNumber != order.TrackNumber {
			return errors.New("track numbers don't match")
		}
	}

	return nil
}

func GetNatsConn() (*nats.Conn, error) {
	nc, err := nats.Connect("nats-server:4222", nats.Name("HTTP Server"))
	if err != nil {
		return nil, err
	}
	return nc, nil
}
