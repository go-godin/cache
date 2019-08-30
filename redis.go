package cache

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type Cache struct {
	client *redis.Client
}

// CacheClient defines an arbitrary cache client
type CacheClient interface {
	Hash(data interface{}) string
	Set(key string, data interface{}, ttl time.Duration) error
	Get(key string, target interface{}) error
	GetKeys(keyPattern string) ([]string, error)
	Delete(keys ...string) error
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

// Hash returns an MD5 hash of the passed data, marshalled into bytes by MessagePack
func (c *Cache) Hash(data interface{}) string {
	b, _ := msgpack.Marshal(data)
	return fmt.Sprintf("%x", md5.Sum(b))
}

// Set a cache entry. The passed data is marshalled using MessagePack
func (c *Cache) Set(key string, data interface{}, ttl time.Duration) error {
	packed, err := msgpack.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "unable to marshal data with MessagePack")
	}

	status := c.client.Set(key, packed, ttl)
	if status.Err() != nil {
		return errors.Wrap(err, "unable to set cache-data")
	}

	return nil
}

// Get a cache entry given it's key and Unmarshal the data into the target.
func (c *Cache) Get(key string, target interface{}) error {
	status := c.client.Get(key)
	if status.Err() != nil {
		return errors.Wrap(status.Err(), "Get failed")
	}

	data, err := status.Bytes()
	if err != nil {
		return errors.Wrap(err, "unable to return data bytes")
	}

	if err := msgpack.Unmarshal(data, target); err != nil {
		return errors.Wrap(err, "unable to unmarshal data")
	}

	return nil
}

// GetKeys returns all keys which match the given keyPattern.
func (c *Cache) GetKeys(keyPattern string) ([]string, error) {
	status := c.client.Keys(keyPattern)
	if status.Err() != nil {
		return nil, errors.Wrap(status.Err(), "unable to get keys by pattern")
	}
	return status.Val(), nil
}

// Delete all passed keys
func (c *Cache) Delete(keys ...string) error {
	status := c.client.Del(keys...)
	if status.Err() != nil {
		return errors.Wrap(status.Err(), "unable to delete keys")
	}
	return nil
}
