package nats

import (
	"L0/database"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

/*
periodically send post-requests to NATS
subscribe to NATS queue and write request to DB & redis cache

simple frontend to get data by order ID from

if server dies -> get cache from DB

*/

func StartReader() error {

	var err error

	nc, err := nats.Connect("0.0.0.0:4222", nats.Name("Reader"))

	if err != nil {
		return fmt.Errorf("could not connect to nats: %v", err)
	}

	bufferSize := 64
	natsChan := make(chan *nats.Msg, bufferSize)

	defer nc.Close()

	sub, err := nc.ChanSubscribe("test", natsChan)

	if err != nil {
		return fmt.Errorf("could not subscribe to nats: %v", err)
	}

	errorChan := make(chan error, 1)

	//Добавляем хэндлеры, которые будут отрабатывать в случае ошибки вычитки, дисконнектов, или закрытого коннекта

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

	for msg := range natsChan {
		order := database.Order{}

		err := json.Unmarshal(msg.Data, &order)

		log.Printf("%+v", order)

		if err != nil {
			errorChan <- err
		}
	}

	_ = sub.Unsubscribe()
	close(natsChan)
	return nil
}
