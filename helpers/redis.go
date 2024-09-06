package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func RedisHealth() bool {
	flag := true

	dbHost := os.Getenv("DB_REDIS_HOST") + ":" + os.Getenv("DB_REDIS_PORT")
	dbPass := os.Getenv("DB_REDIS_PASSWORD")
	dbName, _ := strconv.Atoi(os.Getenv("DB_REDIS_NAME"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     dbHost,
		Password: dbPass,
		DB:       dbName,
	})

	_, err := rdb.Ping(ctx).Result()

	rdb.Close()
	if err != nil {
		flag = false
	}

	return flag

}

func SetRedis(key string, data interface{}, expire time.Duration) {
	dbHost := os.Getenv("DB_REDIS_HOST") + ":" + os.Getenv("DB_REDIS_PORT")
	dbPass := os.Getenv("DB_REDIS_PASSWORD")
	dbName, _ := strconv.Atoi(os.Getenv("DB_REDIS_NAME"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     dbHost,
		Password: dbPass,
		DB:       dbName,
	})
	json, _ := json.Marshal(data)
	err := rdb.Set(ctx, key, json, expire).Err()
	rdb.Close()
	ErrorCheck(err)
}

func GetRedis(key string) (string, bool) {
	flag := true
	dbHost := os.Getenv("DB_REDIS_HOST") + ":" + os.Getenv("DB_REDIS_PORT")
	dbPass := os.Getenv("DB_REDIS_PASSWORD")
	dbName, _ := strconv.Atoi(os.Getenv("DB_REDIS_NAME"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     dbHost,
		Password: dbPass,
		DB:       dbName,
	})

	resultredis, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		flag = false
	}
	rdb.Close()
	return resultredis, flag
}
func DeleteRedis(key string) int {
	dbHost := os.Getenv("DB_REDIS_HOST") + ":" + os.Getenv("DB_REDIS_PORT")
	dbPass := os.Getenv("DB_REDIS_PASSWORD")
	dbName, _ := strconv.Atoi(os.Getenv("DB_REDIS_NAME"))
	value := 0
	rdb := redis.NewClient(&redis.Options{
		Addr:     dbHost,
		Password: dbPass,
		DB:       dbName,
	})
	value = int(rdb.Del(ctx, key).Val())
	rdb.Close()
	return value
}
func IncrPipeRedis(key, db string, expire time.Duration) string {
	dbHost := os.Getenv("DB_REDIS_HOST") + ":" + os.Getenv("DB_REDIS_PORT")
	dbPass := os.Getenv("DB_REDIS_PASSWORD")
	dbName, _ := strconv.Atoi(os.Getenv("DB_REDIS_NAME"))

	if db != "" {
		dbName, _ = strconv.Atoi(db)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     dbHost,
		Password: dbPass,
		DB:       dbName,
	})
	defer rdb.Close()
	pipe := rdb.Pipeline()

	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expire)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Println("increment redis error: ", err)
	}
	s := fmt.Sprintf("%d", incr.Val())

	// The value is available only after Exec is called.
	return s
}
