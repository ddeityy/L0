package main

import (
	"L0/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Response struct {
	Order database.CacheOrder `json:"order"`
}

func main() {
	r := gin.Default()
	rdb := database.GetRedisClient()

	db, err := database.Connect()
	if err != nil {
		log.Panic(err)
	}

	err = db.AutoMigrate(database.DBOrder{}, database.Delivery{}, database.Payment{}, database.OrderItem{})
	if err != nil {
		log.Panic(err)
	}

	err = database.RestoreCacheFromDB(db, rdb)
	if err != nil {
		log.Println("Failed to restore cache:", err)
	}

	r.GET("/orders/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")

		cachedOrder, err := database.GetFromCache(uuid, rdb)
		if err == redis.Nil {
			co, err := database.GetFromDB(uuid, db)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			}
			log.Println("Returning order from database:", uuid)
			c.JSON(http.StatusOK, co)
		}

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		log.Println("Returning cached order:", uuid)
		c.JSON(http.StatusOK, cachedOrder)
	})
	r.Run()
}
