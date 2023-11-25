package cache

import (
	"L0/database"
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if pong := client.Ping(context.Background()); pong.String() != "ping: PONG" {
		log.Println("-------------Error connection redis ----------:", pong)
	}
	return client
}

func RestoreCacheFromDB(db *gorm.DB, rdb *redis.Client) error {
	orders := []database.DBOrder{}
	db.Find(&orders)
	for _, order := range orders {
		delivery := database.Delivery{}
		payment := database.Payment{}
		items := []database.OrderItem{}
		db.First(&delivery, "order_uid = ?", order.OrderUID)
		db.First(&payment, "transaction = ?", order.OrderUID)
		db.Find(&items, "track_number = ?", order.TrackNumber)
		cacheOrder := database.CacheOrder{
			OrderUID:          order.OrderUID,
			TrackNumber:       order.TrackNumber,
			Entry:             order.Entry,
			Locale:            order.Locale,
			InternalSignature: order.InternalSignature,
			CustomerID:        order.CustomerID,
			DeliveryService:   order.DeliveryService,
			Shardkey:          order.Shardkey,
			SmID:              order.SmID,
			DateCreated:       order.DateCreated,
			OofShard:          order.OofShard,
			Delivery:          delivery,
			Payment:           payment,
			Items:             items,
		}
		err := SaveToCache(cacheOrder)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveToCache(order database.CacheOrder) error {
	rdb := GetRedisClient()
	ctx := context.Background()
	jsonOrder, _ := json.Marshal(order)
	err := rdb.Set(ctx, order.OrderUID.String(), jsonOrder, 0).Err()
	if err != nil {
		return err
	}
	log.Println("Saved to cache:", order.TrackNumber, order.Items[0].TrackNumber)
	return nil
}

type KeyNotFoundError struct{}

func (e *KeyNotFoundError) Error() string {
	return "Key does not exist"
}

func GetFromCache(key string, rdb *redis.Client) (string, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", redis.Nil
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}
