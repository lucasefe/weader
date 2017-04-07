package util

// Cache is the data structure where the data is stored
type Cache struct {
	data map[string]interface{}
}

// NewCache constructs a new Cache
func NewCache() *Cache {
	cache := &Cache{}
	cache.data = make(map[string]interface{})

	return cache
}

// FetcherFunc defines the func fetcher type
type FetcherFunc func() (interface{}, error)

// Fetch retrieves the value from the cache, if available, or executes the fetcher to get the value.
func (c *Cache) Fetch(key string, fetcher FetcherFunc) (interface{}, error) {
	if value, ok := c.data[key]; ok {
		return value, nil
	}

	value, err := fetcher()
	if err != nil {
		return nil, err
	}

	c.Set(key, value)

	return value, nil
}

// Get retrieves the value stored by the key from the cache
func (c *Cache) Get(key string) interface{} {
	if value, ok := c.data[key]; ok {
		return value
	}

	return nil
}

// Set stores the value by the key on the cache
func (c *Cache) Set(key string, value interface{}) {
	c.data[key] = value
}
