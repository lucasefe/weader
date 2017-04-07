package util

// Cache is the data structure where the data is stored
type Cache struct {
	Data map[string]interface{}
}

// NewCache constructs a new Cache
func NewCache() *Cache {
	cache := &Cache{}
	cache.Data = make(map[string]interface{})

	return cache
}

// FetcherFunc defines the func fetcher type
type FetcherFunc func() (interface{}, error)

// Fetch retrieves the value from the cache, if available, or executes the fetcher to get the value.
func (c *Cache) Fetch(key string, fetcher FetcherFunc) (interface{}, error) {
	if value, ok := c.Data[key]; ok {
		return value, nil
	}

	value, err := fetcher()
	if err != nil {
		return nil, err
	}

	c.Data[key] = value

	return value, nil
}
