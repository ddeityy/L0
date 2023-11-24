package main

import (
	"L0/database"
	"errors"
)

func validateOrder(order database.CacheOrder) (database.CacheOrder, error) {
	if order.Payment.Transaction != order.OrderUID {
		return order, errors.New("payment does not match")
	}
	for _, item := range order.Items {
		if item.TrackNumber != order.TrackNumber {
			return order, errors.New("tracking does not match")
		}
	}
	return order, nil
}
