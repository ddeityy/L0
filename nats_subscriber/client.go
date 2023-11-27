package main

import (
	"L0/database"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/xeipuuv/gojsonschema"
	"gorm.io/gorm"
)

func StartReader(db *gorm.DB, rdb *redis.Client) error {

	var err error
	errorChan := make(chan error, 1)
	bufferSize := 64
	natsChan := make(chan *nats.Msg, bufferSize)

	nc, err := nats.Connect("nats://0.0.0.0:4222", nats.Name("Reader"))
	if err != nil {
		return fmt.Errorf("could not connect to nats: %v", err)
	}
	defer nc.Close()

	sub, err := nc.ChanQueueSubscribe("order", "orders", natsChan)
	if err != nil {
		return fmt.Errorf("could not subscribe to nats: %v", err)
	}

	nc.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		log.Panicln("Read error:", err.Error())
		errorChan <- err
	})
	nc.SetDisconnectErrHandler(func(_ *nats.Conn, err error) {
		log.Panicln("Reader disconnected:", err.Error())
	})
	nc.SetClosedHandler(func(_ *nats.Conn) {
		log.Panicln("Connection closed")
	})

	schemaPath, err := filepath.Abs("./schema.json")
	if err != nil {
		log.Panicln("Could not locate validation schema:", err)
	}
	schema := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%v", schemaPath))

	wg := sync.WaitGroup{}

	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(ch chan *nats.Msg, wg *sync.WaitGroup) {
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
					errorChan <- err
				}

				err = validateOrder(order)
				if err != nil {
					log.Printf("Invalid order: %v", err)
				} else {
					err = database.SaveToDB(order, db)
					if err != nil {
						log.Println(err)
					} else {
						log.Println("Saved to db:", order.OrderUID, "items:", len(order.Items))
						err := database.SaveToCache(order)
						if err != nil {
							log.Println(err)
						}
						log.Println("Saved to cache:", order.OrderUID)

					}
				}
			}
		}(natsChan, &wg)
	}
	wg.Wait()
	_ = sub.Unsubscribe()
	close(natsChan)
	return nil
}
