package main

import (
	client "L0/nats-client"
	publish "L0/nats-publish"
)

func main() {
	go publish.NatsWriter()
	client.StartReader()
}
