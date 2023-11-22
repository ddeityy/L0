package cache

import "github.com/redis/go-redis/v9"

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6969",
		Password: "password",
		DB:       0,
	})
}

func RestoreCacheFromDB(rdb *redis.Client) error {
	return nil
}
