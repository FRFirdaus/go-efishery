package sqldb

// Configuration is standart template
// configuration for sqldb connection library (this library)
type Configuration struct {
	Key    string `json:"key" yaml:"key"`       // Key is connection key identification
	Dsn    string `json:"dsn" yaml:"dsn"`       // Dsn configuration using uri
	Enable bool   `json:"enable" yaml:"enable"` // Flag database  should start pooling
}

type localConfiguration struct {
	// in second
	MaxConnLifeTime int // If d <= 0, connections are not closed due to a connection's age
	//  in second
	MaxIdleConnections int // If d <= 0, connections are not closed due to a connection's idle time.
	// int total connection
	MaxOpenConnections int // If n <= 0, then there is no limit on the number of open connections.
}
