package main

import (
	"time"
)

func main() {
	for {
		SendFakeData()
		time.Sleep(1 * time.Second)
	}
}
