package main

import (
	"L0/database"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/nats-io/nats.go"
	"github.com/xeipuuv/gojsonschema"
)

func StartReader() error {

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

		order := database.Order{}

		err = json.Unmarshal(msg.Data, &order)

		if err != nil {
			errorChan <- err
		}

		order, err = validateOrder(order)
		if err != nil {
			log.Printf("Invalid order: %v", err)
		} else {
			log.Printf("Valid order: %v", order.OrderUID)
		}
	}

	_ = sub.Unsubscribe()
	close(natsChan)
	return nil
}
