package config

import (
	"fmt"
	"log"
	"sync"
)

var _svcRemoteConfig RemoteConfigService
var _mutex sync.Mutex

func InitRemoteConfig(svc RemoteConfigService) error {
	_mutex.Lock()
	defer _mutex.Unlock()

	if _svcRemoteConfig != nil {
		return fmt.Errorf("Already initialized")
	}
	_svcRemoteConfig = svc
	return nil

}

func CloseRemoteConfig() error {
	_mutex.Lock()
	defer _mutex.Unlock()

	if _svcRemoteConfig == nil {
		return fmt.Errorf("Remote config already closed")
	}

	err := _svcRemoteConfig.Close()
	if err != nil {
		return err
	}

	// Clean
	_svcRemoteConfig = nil

	return nil
}

func RemoteConfig() RemoteConfigService {
	_mutex.Lock()
	defer _mutex.Unlock()

	if _svcRemoteConfig == nil {
		log.Panic("Please init fisrt")
	}
	return _svcRemoteConfig

}
