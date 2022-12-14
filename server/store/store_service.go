package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type StorageService struct {
	redisClient *redis.Client
}

var (
	storeService = &StorageService{}
	userService  = &StorageService{}
	ctx          = context.Background()
)

const CacheDuration = 6 * time.Hour

func InitializeStore() *StorageService {
	redisInfo := fmt.Sprintf("%s:%s", os.Getenv("REDIS_STORE_HOST"), os.Getenv("REDIS_STORE_PORT"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisInfo,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis store instance started successfully: pong message = {%s}\n", pong)

	storeService.redisClient = redisClient
	return storeService
}

func InitializeUser() *StorageService {
	redisInfo := fmt.Sprintf("%s:%s", os.Getenv("REDIS_USER_HOST"), os.Getenv("REDIS_USER_PORT"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisInfo,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis user instance started successfully: pong message = {%s}\n", pong)

	userService.redisClient = redisClient
	return userService
}

func SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key url | Error: %v - short url: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}
}

func RetrieveInitialUrl(shortUrl string) string {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | error: %v - shortUrl: %s\n", err, shortUrl))
	}

	return result
}

func GetRedis() *redis.Client {
	return userService.redisClient
}
