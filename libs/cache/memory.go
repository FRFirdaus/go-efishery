package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

type memoryService struct {
	config MemoryConfig
	client *bigcache.BigCache
	// client
}

type MemoryConfig struct {
	// TTL in second
	TTL int
}

func (ss *memoryService) IsExists(key string) bool {
	_, err := ss.client.Get(key)
	if err != nil {
		return false
	}
	return true
}
func (ss *memoryService) Get(key string, v interface{}) error {
	data, err := ss.client.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
func (ss *memoryService) Set(key string, v interface{}) error {
	var data []byte
	switch vv := v.(type) {
	case []byte:
		data = vv
	default:
		marshalled, err := json.Marshal(v)
		if err != nil {
			return err
		}
		data = marshalled
	}
	return ss.client.Set(key, data)
}
func (ss *memoryService) Delete(key string) error {
	return ss.client.Delete(key)
}
func (ss *memoryService) Close() error {
	return ss.client.Close()
}

func NewMemoryCache(config MemoryConfig) CacheService {
	if config.TTL < 1 {
		log.Panicln("TTL must be greater than 0")
	}
	// https://github.com/allegro/bigcache#custom-initialization
	bigMemoryConfig := bigcache.Config{
		Shards:     1024,
		LifeWindow: time.Duration(config.TTL) * time.Second,
		// 1/2 LifeWindow
		CleanWindow: time.Duration(config.TTL/2) * time.Second,
		// rps * lifeWindow
		// 1000 rps
		MaxEntriesInWindow: 1000 * config.TTL,
		MaxEntrySize:       500,
		HardMaxCacheSize:   8192,
	}

	if bigMemoryConfig.CleanWindow < 1 {
		bigMemoryConfig.CleanWindow = bigMemoryConfig.LifeWindow
	}

	client, err := bigcache.NewBigCache(bigMemoryConfig)
	if err != nil {
		log.Panic(err)
	}

	return &memoryService{
		config: config,
		client: client,
	}
}
