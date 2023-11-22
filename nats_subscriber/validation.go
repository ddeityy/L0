package main

import (
	"L0/database"
	"errors"
)

func validateOrder(order database.Order) (database.Order, error) {
	if order.Payment.Transaction != order.OrderUID {
		return order, errors.New("payment does not match")
	}
	if order.Delivery.OrderUID != order.OrderUID {
		return order, errors.New("delivery does not match")
	}
	for _, item := range order.Items {
		if item.TrackNumber != order.TrackNumber {
			return order, errors.New("tracking does not match")
		}
	}
	return order, nil
}
