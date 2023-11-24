package main

import (
	"L0/cache"
	"L0/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	cache.RestoreCacheFromDB(db, rdb)

	r.GET("/orders/:uuid", func(c *gin.Context) {
		uid := c.Param("uuid")
		cachedOrder, err := cache.GetFromCache(uid, rdb)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, cachedOrder)
	})
	r.Run()
}
