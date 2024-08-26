package redis

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrInvalidRedisOption = errors.New("redis option is invalid")

func NewRedisClient(opt *redis.Options) (*redis.Client, error) {
	if opt == nil {
		return nil, ErrInvalidRedisOption
	}
	return redis.NewClient(opt), nil
}
