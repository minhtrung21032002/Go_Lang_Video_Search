package config

import (
	"context"
	"log"
	"os"
	"strconv" // Corrected import statement

	"github.com/go-redis/redis/v8"
)

func RedisClient() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable not set")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		// Handle error converting REDIS_DB to an integer
		log.Fatalf("Error converting REDIS_DB to integer: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test the connection
	err = rdb.Ping(context.Background()).Err() // Reusing 'err' variable
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return rdb
}
