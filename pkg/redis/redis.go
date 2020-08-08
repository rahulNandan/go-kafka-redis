package redis

import (
	"fmt"

	redisImpl "github.com/go-redis/redis"
)

//Redis - client
type Redis struct {
	ID     string
	client *redisImpl.Client
}

//New - return redis obj with connect
func New(id, host, port, password string) (*Redis, error) {
	opts := &redisImpl.Options{
		Addr:     host + ":" + port,
		Password: password,
	}
	client := redisImpl.NewClient(opts)

	pong, err := client.Ping().Result()

	r := Redis{
		ID:     id,
		client: client,
	}

	if err != nil && pong == "" {
		return nil, err
	}

	fmt.Printf("Connected to Redis %s:%s\n", host, port)
	return &r, nil
}

//Close - close redis connections
func (r *Redis) Close() error {
	fmt.Printf("Closing %s Redis connection... \n", r.ID)
	return r.client.Close()
}

//Read - read from redis.
func (r *Redis) Read(keys []string) ([]interface{}, error) {
	values, err := r.client.MGet(keys...).Result()

	if err != nil {
		return nil, err
	}

	return values, nil
}

//Incr - incr a key in redis
func (r *Redis) Incr(key string) error {
	_, err := r.client.Incr(key).Result()
	if err != nil {
		return err
	}
	return nil
}
