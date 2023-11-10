package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	cache *redis.Client
}

func NewCache(ctx context.Context) (*Cache, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := cli.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &Cache{cache: cli}, nil
}

func (c *Cache) Clear() error {
	err := c.cache.FlushDB(context.Background()).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Close() {
	log.Println("Closing cache connection")

	err := c.Clear()
	if err != nil {
		panic(err)
	}

	er := c.cache.Close()
	if er != nil {
		panic(er)
	}
}

func (c *Cache) GetCache() *redis.Client {
	return c.cache
}
