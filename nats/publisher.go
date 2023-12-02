package nats

import (
	"L0/database"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

func StartNatsPub(delay time.Duration) {
	nc, err := GetNatsConn()
	if err != nil {
		log.Panic("could not connect to nats:", err)

	}
	defer nc.Close()

	sc, err := stan.Connect("cluster", "pub", stan.NatsConn(nc))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, "nats-server:4222")
	}
	defer sc.Close()

	subject := "order"

	for {
		order := CreateFakeOrder(randBool(), randBool())

		b, err := json.Marshal(order)
		if err != nil {
			log.Println(err)
		}

		err = sc.Publish(subject, b)
		if err != nil {
			log.Println(err)
		}

		log.Println("Published:", order.OrderUID)
		time.Sleep(delay)
	}
}

// Returns a slice of (potentially incorrect) OrderItems
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

// Returns a fake (potentially incorrect) CacheOrder
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
	payment.PaymentDt = int(gofakeit.Date().Unix())
	order.Payment = payment

	gofakeit.Struct(&order)

	return &order
}

func randBool() bool {
	return rand.Intn(20) == 1
}
