package main

import (
	"L0/database"
	nats "L0/nats"
	"encoding/json"
	"log"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Response struct {
	Order database.CacheOrder `json:"order"`
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{AllowOrigins: []string{"*"}}))

	r.LoadHTMLGlob("templates/*")
	rdb := database.GetRedisClient()

	nc, err := nats.GetNatsConn()
	if err != nil {
		log.Panic(err)
	}

	subject := "order"

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
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Nats",
		})
	})
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
	r.POST(("/orders/new"), func(c *gin.Context) {
		order := nats.CreateFakeOrder(false, false)
		b, err := json.Marshal(order)
		if err != nil {
			log.Println(err)
		}

		err = nc.Publish(subject, b)
		if err != nil {
			log.Println(err)
		}

		log.Println("Created new order from HTTP request:", order.OrderUID)
		c.JSON(http.StatusCreated, gin.H{"OrderUID:": order.OrderUID})

	})
	go nats.StartNatsPub()
	go nats.StartNatsSub()
	r.Run()
}
