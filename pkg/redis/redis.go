package redisclient

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Password string
	DB       int
	PoolSize int
}

type Client struct {
	cli *redis.Client
}

func NewRedisClient(c *Config) *Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Host,
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
		PoolSize: c.PoolSize, // 连接池连接数量
	})
	return &Client{redisClient}
}

func (c *Client) Close() {
	if c.cli == nil {
		return
	}
	c.cli.Close()
}
