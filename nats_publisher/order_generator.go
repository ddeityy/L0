package main

import (
	"L0/database"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func createItems(n int, trackNumber string, incorrect bool) []database.OrderItem {
	var items []database.OrderItem

	for i := 0; i < n+1; i++ {
		item := database.OrderItem{}

		if incorrect {
			item.TrackNumber = uuid.New().String()
		} else {
			item.TrackNumber = trackNumber
		}

		gofakeit.Struct(&item)
		items = append(items, item)
	}
	return items
}

func CreateFakeOrder(badPayment bool, badItems bool) *database.CacheOrder {
	order := database.CacheOrder{}
	order.OrderUID = uuid.New().String()
	order.TrackNumber = uuid.New().String()
	delivery := database.Delivery{}
	delivery.OrderUID = order.OrderUID
	payment := database.Payment{}

	order.Items = createItems(rand.Intn(10), order.TrackNumber, badItems)

	gofakeit.Struct(&delivery)
	order.Delivery = delivery

	if badPayment {
		payment.Transaction = uuid.New().String()
	} else {
		payment.Transaction = order.OrderUID
	}

	gofakeit.Struct(&payment)
	order.Payment = payment

	gofakeit.Struct(&order)

	return &order
}

func randBool() bool {
	return rand.Intn(20) == 1
}
