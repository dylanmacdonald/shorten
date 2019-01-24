package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dylanmacdonald/shorten/api"
)

// Service is a container for external dependancies and concrete business logic
type Service struct {
	Store Store
}

type Store interface {
	Get(ctx context.Context, key int64) (string, error)
	Set(ctx context.Context, key int64, value string) error
	NotFound(error) bool
}

// Shorten creates a hash of the request url and returns it
// if the hash doesn't already exist in the store a pair of hash:URL is created
// if the hash does exist the url is returned
func (s Service) Shorten(ctx context.Context, req *api.ShortenRequest) (*url.URL, error) {
	key := hash(req.URL)

	uri, err := url.Parse(fmt.Sprintf("https://shorten.cluster.host-name.com/%d", key))
	if err != nil {
		return nil, err
	}

	_, err = s.Store.Get(ctx, key)
	if err == nil {
		return uri, nil
	}

	if s.Store.NotFound(err) {
		err = s.Store.Set(ctx, key, req.URL)
		return uri, err
	}

	return uri, err
}

/**
Shamelessly stolen from http://www.cse.yorku.ca/~oz/hash.html
this is a go implementation of djb2
In a production environment I'd recommend something like https://golang.org/pkg/hash/fnv/

unsigned long hash(unsigned char *str) {
	unsigned long hash = 5381;
	int c;

	while (c = *str++)
		hash = ((hash << 5) + hash) + c; /* hash * 33 + c
		return hash;
	}
}
*/
func hash(s string) int64 {
	var hash int64 = 5381

	for _, c := range s {
		// hash * 33 + c
		hash = ((hash << 5) + hash) + int64(c)
	}
	return hash
}

func (s Service) Redirect(ctx context.Context, path int64) (*url.URL, error) {
	u, err := s.Store.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	uri, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if uri.Scheme == "" {
		uri.Scheme = "http"
	}
	return uri, nil
}
