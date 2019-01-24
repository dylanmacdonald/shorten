package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/dylanmacdonald/shorten/service"
	"github.com/dylanmacdonald/shorten/store"

	"github.com/Sirupsen/logrus"
	"github.com/dylanmacdonald/shorten/handlers"
)

func main() {
	logger := logrus.WithField("service", "image-service")
	redisHost, found := os.LookupEnv("REDIS_HOST")
	if !found {
		logger.Fatal("REDIS_HOST was not found in env")
	}
	redisPort, found := os.LookupEnv("REDIS_PORT")
	if !found {
		logger.Fatal("REDIS_PORT was not found in env")
	}

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			address := fmt.Sprintf("%s:%s", redisHost, redisPort)
			c, err := redis.Dial("tcp", address)

			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	s := service.Service{
		Store: store.Store{
			Pool: redisPool,
		},
	}
	logger.Info("listening on port 8080")
	logger.Fatal(http.ListenAndServe(":8080", handlers.InitHandlers(
		logger,
		s,
	)))
}
