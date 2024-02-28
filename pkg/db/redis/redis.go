package dbredis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arfan21/vocagame/config"
	"github.com/arfan21/vocagame/pkg/logger"
	"github.com/redis/go-redis/v9"
)

func New() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GetConfig().Redis.URL, config.GetConfig().Redis.Port),
		Username: config.GetConfig().Redis.Username,
		Password: config.GetConfig().Redis.Password,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		logger.Log(context.Background()).Error().Err(err).Msg("failed to ping redis")
		return nil, err
	}

	logger.Log(context.Background()).Info().Msg("dbredis: connection established")

	return client, nil
}

func Set(ctx context.Context, client *redis.Client, key string, value any, expiration time.Duration) error {
	valueByte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, string(valueByte), expiration).Err()
}

func Get[T any](ctx context.Context, client *redis.Client, key string) (res T, err error) {
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}

	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
