package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/Just4Ease/nuMicro/store"
)

type Redis struct {
	client *redis.Client
}

func New(address, password string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	return &Redis{
		client: client,
	}
}

// List all the known records

func (r *Redis) Init(...store.Option) error {
	return nil
}

func (r *Redis) List() ([]*store.Record, error) {
	return nil, nil
}

func (r *Redis) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	_, err := r.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	//records := make([]*store.Record, 0)

	//return []*store.Record{
	//	Key:    key,
	//	Value:  []byte(val),
	//	Expiry: 0,
	//}, nil

	return nil, nil
}

func (r *Redis) Write(record *store.Record) error {
	expiry := time.Duration(record.Expiry)
	if err := r.client.Set(record.Key, string(record.Value), expiry).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Delete(key string) error {
	return nil
}

func (r *Redis) String() string {
	return "redis"
}
