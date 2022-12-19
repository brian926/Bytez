package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/brian926/Bytez/server/db"
	"github.com/brian926/Bytez/server/forms"
	"github.com/go-redis/redis/v8"
)

type StorageService struct {
	redisClient *redis.Client
}

type Url struct {
	ID        int64  `db:"id, primarykey, autoincrement" json:"id"`
	ShortUrl  string `db:"shortUrl" json:"shortUrl"`
	LongUrl   string `db:"longUrl" json:"longUrl"`
	CreatedAt int64  `db:"created_at" json:"-"`
}

type UrlModel struct{}

var (
	storeService = &StorageService{}
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

func (u UrlModel) SaveUrlMapping(urlCreation forms.UrlCreationRequest) (url Url, err error) {
	checkUser, err := db.GetDB().SelectInt("SELECT count(id) FROM public.urls WHERE longUrl=LOWER($1) LIMIT 1", urlCreation.LongUrl)
	if err != nil {
		return url, errors.New("something went wrong checking for long url, please try again later")
	}

	if checkUser > 0 {
		return url, errors.New("mapping already exists")
	}

	err = db.GetDB().QueryRow("INSERT INTO public.urls(shortUrl, longUrl) VALUES($1, $2) RETURNING id", urlCreation.ShortUrl, urlCreation.LongUrl).Scan(&url.ID)
	if err != nil {
		return url, errors.New("something went wrong creating url, please try again later")
	}

	url.ShortUrl = urlCreation.ShortUrl
	url.LongUrl = urlCreation.LongUrl

	fmt.Println("redis?")
	errAccess := storeService.redisClient.Set(ctx, url.ShortUrl, url.LongUrl, CacheDuration).Err()
	if errAccess != nil {
		return url, errors.New("failed to set reddis")
	}

	return url, nil
}

func RetrieveInitialUrl(shortUrl string) string {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed RetrieveInitialUrl url | error: %v - shortUrl: %s\n", err, shortUrl))
	}

	return result
}

func GetRedis() *redis.Client {
	return storeService.redisClient
}
