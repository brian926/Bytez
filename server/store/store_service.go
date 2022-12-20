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

type UrlModel struct{}
type StorageService struct {
	redisClient *redis.Client
}

type Url struct {
	ID        int64  `db:"id, primarykey, autoincrement" json:"id"`
	ShortUrl  string `db:"shortUrl" json:"shortUrl"`
	LongUrl   string `db:"longUrl" json:"longUrl"`
	CreatedAt int64  `db:"created_at" json:"-"`
}

var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

const CacheDuration = 6 * time.Hour

func (u UrlModel) InitializeStore() *StorageService {
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
	// Check redis first if already saved, else save it
	result, err := storeService.redisClient.Get(ctx, urlCreation.ShortUrl).Result()
	if err != nil {
		checkUser, err := db.GetDB().SelectInt("SELECT count(id) FROM public.urls WHERE longUrl=LOWER($1) LIMIT 1", urlCreation.LongUrl)
		if err != nil {
			return url, errors.New("something went wrong checking for long url, please try again later")
		}

		// even though the URL already exist, return with no errors
		if checkUser > 0 {
			return url, nil
		}

		err = db.GetDB().QueryRow("INSERT INTO public.urls(shortUrl, longUrl) VALUES($1, $2) RETURNING id", urlCreation.ShortUrl, urlCreation.LongUrl).Scan(&url.ID)
		if err != nil {
			return url, errors.New("something went wrong creating url, please try again later")
		}

		url.ShortUrl = urlCreation.ShortUrl
		url.LongUrl = urlCreation.LongUrl

		errAccess := storeService.redisClient.Set(ctx, url.ShortUrl, url.LongUrl, CacheDuration).Err()
		if errAccess != nil {
			return url, errors.New("failed to set reddis")
		}

		return url, nil
	}

	url.LongUrl = result
	url.ShortUrl = urlCreation.ShortUrl
	return url, nil
}

func (u UrlModel) RetrieveInitialUrl(urlCreation forms.UrlCreationRequest) (url Url, err error) {
	result, err := storeService.redisClient.Get(ctx, urlCreation.ShortUrl).Result()
	if err != nil {
		err = db.GetDB().SelectOne(&url, "SELECT shortUrl, longUrl FROM public.urls WHERE shortUrl=LOWER($1) LIMIT 1", urlCreation.ShortUrl)

		if err != nil {
			return url, errors.New("failed to set database urls")
		}

		errAccess := storeService.redisClient.Set(ctx, url.ShortUrl, url.LongUrl, CacheDuration).Err()
		if errAccess != nil {
			return url, errors.New("failed to set reddis")
		}
	}

	url.LongUrl = result

	return url, nil
}

func GetRedis() *redis.Client {
	return storeService.redisClient
}
