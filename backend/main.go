package main

import (
	"L0/cache"
	"L0/database"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Order database.CacheOrder `json:"order"`
}

func main() {
	r := gin.Default()
	db, err := database.Connect()
	if err != nil {
		log.Panic(err)
	}
	err = db.AutoMigrate(database.DBOrder{}, database.Delivery{}, database.Payment{}, database.OrderItem{})
	if err != nil {
		log.Panic(err)
	}
	rdb := cache.GetRedisClient()

	//cache.RestoreCacheFromDB(db, rdb)

	r.GET("/orders/:uuid", func(c *gin.Context) {
		uid := c.Param("uuid")
		cachedOrder, err := cache.GetFromCache(uid, rdb)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		log.Println("Returning cached order:", uid)
		co := database.CacheOrder{}
		json.Unmarshal([]byte(cachedOrder), &co)
		c.JSON(http.StatusOK, co)
	})
	r.Run()
}
