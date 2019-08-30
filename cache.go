package cache

import (
	"time"
)

// Cache defines an arbitrary cache implementation
type Cache interface {
	Hash(data interface{}) string
	Set(key string, data interface{}, ttl time.Duration) error
	Get(key string, target interface{}) error
	GetKeys(keyPattern string) ([]string, error)
	Delete(keys ...string) error
}
