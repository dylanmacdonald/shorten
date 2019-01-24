package store

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type Store struct {
	Pool *redis.Pool
}

// Get retrieves a value from redis
// redigo.ErrNil is returned if the value isn't found
func (s Store) Get(ctx context.Context, key int64) (string, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	return reply, err
}

// Set stores a value in redis
func (s Store) Set(ctx context.Context, key int64, value string) error {
	conn := s.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) NotFound(err error) bool {
	return err == redis.ErrNil
}
