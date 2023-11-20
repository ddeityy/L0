package main

import "L0/nats"

func main() {
	go nats.NatsWriter()
	nats.StartReader()
}
