package main

import (
	"L0/database"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func createItems(n int, trackNumber uuid.UUID, incorrect bool) []database.OrderItem {
	var items []database.OrderItem
	for i := 0; i < n; i++ {
		item := database.OrderItem{}
		if incorrect {
			item.TrackNumber = uuid.New()
		} else {
			item.TrackNumber = trackNumber
		}
		gofakeit.Struct(item)
		items = append(items, item)
	}
	return items
}

func CreateFakeOrder(badPayment bool, badItems bool, badDelivery bool) database.Order {
	order := database.Order{}
	order.OrderUID = uuid.New()
	order.TrackNumber = uuid.New()
	delivery := database.Delivery{}
	payment := database.Payment{}

	order.Items = createItems(rand.Intn(10), order.TrackNumber, badItems)

	if badDelivery {
		delivery.OrderUID = uuid.New()
	} else {
		delivery.OrderUID = order.OrderUID
	}
	gofakeit.Struct(&delivery)
	order.Delivery = delivery

	if badPayment {
		payment.Transaction = uuid.New()
	} else {
		payment.Transaction = order.OrderUID
	}
	gofakeit.Struct(&payment)
	order.Payment = payment

	gofakeit.Struct(&order)

	return order
}

func randBool() bool {
	return rand.Intn(20) == 1
}
