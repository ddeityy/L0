package main

import (
	"L0/cache"
	"L0/database"
	"log"
)

func main() {
	db, err := database.Connect()
	rdb := cache.GetRedisClient()
	if err != nil {
		log.Panic("Could not connect to db:", err)
	}
	err = db.AutoMigrate(database.DBOrder{}, database.Delivery{}, database.Payment{}, database.OrderItem{})
	if err != nil {
		log.Panic(err)
	}
	StartReader(db, rdb)
}
