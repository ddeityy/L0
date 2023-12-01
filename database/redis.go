package database

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"runtime"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	if pong := client.Ping(context.Background()); pong.String() != "ping: PONG" {
		log.Println("-------------Error connection redis ----------:", pong)
	}
	return client
}

func chunkify(orders []DBOrder) [][]DBOrder {
	max := runtime.GOMAXPROCS(0)
	var divided [][]DBOrder

	chunkSize := (len(orders) + max - 1) / max
	for i := 0; i < len(orders); i += chunkSize {
		end := i + chunkSize

		if end > len(orders) {
			end = len(orders)
		}

		divided = append(divided, orders[i:end])
	}
	return divided
}

func RestoreCacheFromDB(db *gorm.DB, rdb *redis.Client) error {
	orders := []DBOrder{}
	db.Find(&orders)
	log.Println("Orders to restore from cache:", len(orders))
	if len(orders) == 0 {
		return errors.New("no orders to restore")
	}
	batches := chunkify(orders)
	var wg sync.WaitGroup
	for i := 0; i < len(batches); i++ {
		wg.Add(1)
		go func(batch []DBOrder) {
			defer wg.Done()
			for _, order := range batch {
				delivery := Delivery{}
				payment := Payment{}
				items := []OrderItem{}
				db.First(&delivery, "order_uid = ?", order.OrderUID)
				db.First(&payment, "transaction = ?", order.OrderUID)
				db.Find(&items, "track_number = ?", order.TrackNumber)

				delivery.OrderUID = ""

				cacheOrder := CacheOrder{
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
				err := SaveToCache(cacheOrder, rdb)
				if err != nil {
					log.Println(err)
				}
			}
		}(batches[i])
	}
	wg.Wait()
	log.Println("Restored from cache:", len(orders), "orders")
	return nil
}

func SaveToCache(order CacheOrder, rdb *redis.Client) error {
	ctx := context.Background()
	jsonOrder, _ := json.Marshal(order)
	err := rdb.Set(ctx, order.OrderUID, jsonOrder, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetFromCache(key string, rdb *redis.Client) (*CacheOrder, error) {
	val, err := rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, redis.Nil
	} else if err != nil {
		return nil, err
	} else {
		c := CacheOrder{}
		err := json.Unmarshal([]byte(val), &c)
		if err != nil {
			log.Println(err)
		}
		return &c, nil
	}
}
