package dbredis

import (
	"context"
	"fmt"

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
