package redis

import (
	"fmt"
	logger "go-todolist/utils/log"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

func InitRedis() *redis.Client {
	// Load .env file
	errEnv := godotenv.Load()
	if errEnv != nil {
		logger.Panic("Failed to load env file")
	}

	dbHost := os.Getenv("REDIS_HOST")
	dbPassword := os.Getenv("REDIS_PASSWORD")
	dbPort := os.Getenv("REDIS_PORT")

	log.Println("Testing Golang Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Password: dbPassword, // no password set
		DB:       0,          // use default DB
	})

	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		fmt.Println("Failed to connect redis : " + err.Error())
		logger.Panic("Failed to connect redis : " + err.Error())
	}

	return rdb
}

func Close(rdb *redis.Client) {
	rdb.Close()
}
