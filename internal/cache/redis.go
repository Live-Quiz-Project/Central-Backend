package cache

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	cache *redis.Client
}

func NewCache(ctx context.Context) (*Cache, error) {
	host := os.Getenv("CACHE_HOST")
	port := os.Getenv("CACHE_PORT")
	pass := os.Getenv("CACHE_PASS")
	
	address := fmt.Sprintf("%s:%s",host,port)

	cli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pass,
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
