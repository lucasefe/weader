package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

const cacheFileSuffix = "cache"

// FetcherFunc defines the func fetcher type
type FetcherFunc func() (interface{}, error)

// Cache is the data structure where the data is stored
type Cache struct {
	name string
	data map[string]interface{}
	sync.RWMutex
}

// NewCache constructs a new Cache
func NewCache(name string) (*Cache, error) {
	cache := &Cache{name: name}
	cache.data = make(map[string]interface{})

	if err := cache.load(); err != nil {
		return cache, err
	}

	return cache, nil
}

// Fetch retrieves the value from the cache, if available, or executes the fetcher to get the value.
func (c *Cache) Fetch(key string, fetcher FetcherFunc) (interface{}, error) {
	if value, ok := c.data[key]; ok {
		return value, nil
	}

	value, err := fetcher()
	if err != nil {
		return nil, err
	}

	if err := c.Set(key, value); err != nil {
		return nil, err
	}

	return value, nil
}

// Get retrieves the value stored by the key from the cache
func (c *Cache) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()

	if value, ok := c.data[key]; ok {
		return value
	}

	return nil
}

// Set stores the value by the key on the cache
func (c *Cache) Set(key string, value interface{}) error {
	c.Lock()
	defer c.Unlock()

	c.data[key] = value

	if err := c.sync(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) load() error {
	if _, err := os.Stat(c.cachePath()); os.IsNotExist(err) {
		return err
	}

	data, err := ioutil.ReadFile(c.cachePath())
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(data), c.data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) sync() error {
	data, err := yaml.Marshal(&c.data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.cachePath(), data, os.FileMode(0644))
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) cachePath() string {
	return fmt.Sprintf(".%s-%s.yml", c.name, cacheFileSuffix)
}
