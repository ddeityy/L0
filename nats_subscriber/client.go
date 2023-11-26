package main

import (
	"L0/cache"
	"L0/database"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/xeipuuv/gojsonschema"
	"gorm.io/gorm"
)

func StartReader(db *gorm.DB, rdb *redis.Client) error {

	var err error

	nc, err := nats.Connect("nats://0.0.0.0:4222", nats.Name("Reader"))

	if err != nil {
		return fmt.Errorf("could not connect to nats: %v", err)
	}

	bufferSize := 64
	natsChan := make(chan *nats.Msg, bufferSize)

	defer nc.Close()

	sub, err := nc.ChanSubscribe("order", natsChan)

	if err != nil {
		return fmt.Errorf("could not subscribe to nats: %v", err)
	}

	errorChan := make(chan error, 1)

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

	for msg := range natsChan {

		jsonData := gojsonschema.NewStringLoader(string(msg.Data))

		result, err := gojsonschema.Validate(schema, jsonData)
		if err != nil {
			log.Println("Invalid schema:", err)
		}
		if result.Valid() {
			log.Printf("The document is valid\n")
		} else {
			log.Printf("The document is not valid. see errors :\n")
			for _, desc := range result.Errors() {
				log.Printf("- %s\n", desc)
			}
		}

		order := database.CacheOrder{}

		err = json.Unmarshal(msg.Data, &order)

		if err != nil {
			errorChan <- err
		}

		order, err = validateOrder(order)
		if err != nil {
			log.Printf("Invalid order: %v", err)
		} else {
			_, err := cache.GetFromCache(order.OrderUID, rdb)
			log.Println(err)
			if err == redis.Nil {
				err = database.SaveToDB(order, db)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("Saved to db:", order.TrackNumber, order.Items[0].TrackNumber)
					err := cache.SaveToCache(order)
					if err != nil {
						log.Println(err)
					}
					log.Println("Saved to cache:", order.OrderUID)

				}
			}
		}

	}
	_ = sub.Unsubscribe()
	close(natsChan)
	return nil
}
