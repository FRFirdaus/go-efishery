package sqldb

import (
	"net/url"
	"strconv"
	"time"
)

const maxConnLifeTime = "max-ttl-conn"
const maxIdleConnections = "max-idle-conn"
const maxOpenConn = "max-conn"

// ParseURL is parse extra query from dsn
// Original sourcecode https://github.com/lib/pq/blob/master/url.go#L32
func ParseURL(dsn string) (*localConfiguration, error) {
	// default is zero
	defaultLocalConfig := localConfiguration{}

	u, err := url.Parse(dsn)
	if err != nil {
		return &defaultLocalConfig, err
	}

	accrue := func(k, v string) {
		if v == "" {
			return
		}

		switch k {
		case maxOpenConn:
			defaultLocalConfig.MaxOpenConnections, _ = strconv.Atoi(v)
		case maxIdleConnections:
			defaultLocalConfig.MaxIdleConnections, _ = strconv.Atoi(v)
			defaultLocalConfig.MaxIdleConnections = defaultLocalConfig.MaxIdleConnections * int(time.Second)
		case maxConnLifeTime:
			defaultLocalConfig.MaxConnLifeTime, _ = strconv.Atoi(v)
			defaultLocalConfig.MaxConnLifeTime = defaultLocalConfig.MaxConnLifeTime * int(time.Second)
		}
	}

	q := u.Query()
	for k := range q {
		accrue(k, q.Get(k))
	}

	return &defaultLocalConfig, nil
}
