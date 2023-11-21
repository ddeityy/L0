package main

import "time"

func main() {
	go SendFakeData(1 * time.Second)
	StartReader()
}
