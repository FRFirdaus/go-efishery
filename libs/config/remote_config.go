package config

import "context"

type RemoteConfigService interface {
	// Read is get value from remote config
	// identifier mean id / path / dir of secret value
	Read(identifier string) (Value, error)
	ReadWithContext(ctx context.Context, identifier string) (Value, error)

	// Close is close pool connection to remote config
	Close() error
}

// abstraction of value
type Value interface {
	// return byte of data
	Data() []byte
	// return meta data of value
	Metadata() map[string]interface{}
}
