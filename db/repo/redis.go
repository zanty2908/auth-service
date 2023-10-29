package repo

import (
	"auth-service/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Caching interface {
	Get(c context.Context, key string) (interface{}, error)
	Set(c context.Context, key string, value interface{}, exp time.Duration) (interface{}, error)
	Exists(c context.Context, key ...string) (bool, error)
}

func NewRedisClient(config *config.Config) *redis.Client {
	redisCfg := config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%d`, redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		// DB:           redisCfg.Database, // use default DB
		// Protocol:     3,                 // specify 2 for RESP 2 or 3 for RESP 3
		// PoolSize:     redisCfg.PoolSize,
		// WriteTimeout: redisCfg.WriteTimeout,
		// ReadTimeout:  redisCfg.ReadTimeout,
		// DialTimeout:  redisCfg.DialTimeout,
	})

	ctx := context.Background()
	res, err := rdb.Ping(ctx).Result()
	log.Info().Msgf("Connect redis: %s - err: %v - instance: %v", res, err, rdb)

	return rdb
}

func (s *RepoImpl) Get(c context.Context, key string) (interface{}, error) {
	return s.Redis.Get(c, key).Result()
}

func (s *RepoImpl) Set(c context.Context, key string, value interface{}, exp time.Duration) (interface{}, error) {
	return s.Redis.Set(c, key, value, exp).Result()
}

func (s *RepoImpl) Exists(c context.Context, key ...string) (bool, error) {
	res, err := s.Redis.Exists(c, key...).Result()
	return res > 0, err
}
