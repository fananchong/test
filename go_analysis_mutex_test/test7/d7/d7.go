package d7

import (
	"sync"
	"time"
)

/*
Usage:
	c := d7.New(10*time.Second, 1*time.Hour)

	c.Set("hello", true)
*/

type TimeoutCache struct {
	mutex sync.Mutex // data
	data  map[interface{}]bool
}

func New(expiration time.Duration, cleanDuration time.Duration) *TimeoutCache {
	c := &TimeoutCache{
		data: make(map[interface{}]bool),
	}

	{
		c.mutex.Lock()
		defer c.mutex.Unlock()
		for k := range c.data {
			delete(c.data, k)
		}
	}

	return c
}
