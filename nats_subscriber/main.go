package main

import (
	"L0/database"
	"log"
)

func main() {
	rdb := database.GetRedisClient()

	db, err := database.Connect()
	if err != nil {
		log.Panic("Could not connect to db:", err)
	}

	err = db.AutoMigrate(database.DBOrder{}, database.Delivery{}, database.Payment{}, database.OrderItem{})
	if err != nil {
		log.Panic(err)
	}

	StartReader(db, rdb)
}
