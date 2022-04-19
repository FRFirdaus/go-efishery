package config

import (
	"context"
	"log"

	"bitbucket.org/efishery/go-efishery/libs/cache"
)

type memoryCache struct {
	cacheService  cache.CacheService
	remoteService RemoteConfigService
}

type cacheValue struct {
	JsonStringData string                 `json:"json_string_data,omitempty"`
	Meta           map[string]interface{} `json:"meta,omitempty"`
}

func (v cacheValue) Data() []byte {
	return []byte(v.JsonStringData)
}

func (v cacheValue) Metadata() map[string]interface{} {
	return v.Meta
}

func (ss *memoryCache) Read(identifier string) (Value, error) {
	return ss.ReadWithContext(context.Background(), identifier)
}

func (ss *memoryCache) ReadWithContext(ctx context.Context, identifier string) (Value, error) {
	var cacheData cacheValue
	err := ss.cacheService.Get(identifier, &cacheData)
	if err == nil {
		return cacheData, nil
	}

	data, err := ss.remoteService.ReadWithContext(ctx, identifier)
	if err != nil {
		return cacheValue{}, err
	}
	err = ss.cacheService.Set(identifier, cacheValue{
		JsonStringData: string(data.Data()),
		Meta:           data.Metadata(),
	})

	if err != nil {
		return cacheValue{}, err
	}
	return data, nil

}

func (ss *memoryCache) Close() error {
	return ss.remoteService.Close()
}

// WithCache is wrapper cache for remote config
// TTL in Second
// You can extend using other cache implementations
func WithCache(remoteService RemoteConfigService, cacheService cache.CacheService) RemoteConfigService {
	if remoteService == nil {
		log.Panic("Remote config service is nil")
	}

	if cacheService == nil {
		log.Panic("Cache service is nil")
	}

	return &memoryCache{remoteService: remoteService, cacheService: cacheService}
}
