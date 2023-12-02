package nats

import (
	"L0/database"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/redis/go-redis/v9"
	"github.com/xeipuuv/gojsonschema"
	"gorm.io/gorm"
)

func StartNatsSub() {
	rdb := database.GetRedisClient()

	db, err := database.Connect()
	if err != nil {
		log.Panic("Could not connect to db:", err)
	}

	nc, err := GetNatsConn()
	if err != nil {
		log.Panic("could not connect to nats:", err)
	}
	defer nc.Close()

	err = db.AutoMigrate(database.DBOrder{}, database.Delivery{}, database.Payment{}, database.OrderItem{})
	if err != nil {
		log.Panic(err)
	}

	StartReader(db, rdb, nc)
}

func StartReader(db *gorm.DB, rdb *redis.Client, nc *nats.Conn) error {

	var err error
	bufferSize := 64
	msgCh := make(chan *stan.Msg, bufferSize)
	defer close(msgCh)

	sc, err := stan.Connect("cluster", "sub", stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, "nats-server:4222")
	}
	defer sc.Close()

	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", "nats-server:4222", "cluster", "sub")

	mcb := func(msg *stan.Msg) {
		msgCh <- msg
	}

	sub, err := sc.QueueSubscribe("order", "orders", mcb)
	if err != nil {
		return fmt.Errorf("could not subscribe to nats: %v", err)
	}

	schemaPath, err := filepath.Abs("./schema.json")
	if err != nil {
		log.Panicln("Could not locate validation schema:", err)
	}
	schema := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%v", schemaPath))

	wg := sync.WaitGroup{}
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(ch chan *stan.Msg, wg *sync.WaitGroup) {
			for msg := range ch {
				jsonData := gojsonschema.NewStringLoader(string(msg.Data))

				result, err := gojsonschema.Validate(schema, jsonData)
				if err != nil {
					log.Println(err)
				}
				if !result.Valid() {
					log.Println("Invalid schema:")
					for _, desc := range result.Errors() {
						log.Printf("- %s\n", desc)
					}
				}

				order := database.CacheOrder{}

				err = json.Unmarshal(msg.Data, &order)
				if err != nil {
					log.Println(err)
				}

				err = ValidateOrder(order)
				if err != nil {
					log.Printf("Invalid order: %v", err)
				} else {
					err = database.SaveToDB(order, db)
					if err != nil {
						log.Println(err)
					} else {
						log.Println("Saved to db:", order.OrderUID, "items:", len(order.Items))
						err := database.SaveToCache(order, rdb)
						if err != nil {
							log.Println(err)
						}
						log.Println("Saved to cache:", order.OrderUID)

					}
				}
			}
		}(msgCh, &wg)
	}
	wg.Wait()
	_ = sub.Unsubscribe()
	return nil
}
