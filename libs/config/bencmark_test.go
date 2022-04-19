package config

import (
	"testing"

	"bitbucket.org/efishery/go-efishery/libs/cache"
)

func BenchmarkWithCache(b *testing.B) {
	testCases := []struct {
		config VaultConfig
		path   string
	}{
		{
			config: VaultConfig{
				Host:     vaultaddress,
				Username: "your username vault", // add your valid token here
				Password: "your password vault",
			},
			path: "efishery/data/service_payment_platform_mapping/staging",
		},
	}
	for _, tC := range testCases {
		var err error
		var svc RemoteConfigService
		svc, err = NewVault(tC.config)
		if err != nil {
			b.Error(err)
		}
		svc = WithCache(
			svc,
			cache.NewMemoryCache(cache.MemoryConfig{TTL: 60}),
		)
		// close auto renew
		defer svc.Close()
		for n := 0; n < b.N; n++ {
			for i := 0; i < 100; i++ {
				_, err := svc.Read(tC.path)
				if err != nil {
					b.Error(err)
				}
			}

		}

		defer svc.Close()
	}
}

func BenchmarkVault(b *testing.B) {
	testCases := []struct {
		config VaultConfig
		path   string
	}{
		{
			config: VaultConfig{
				Host:     vaultaddress,
				Username: "your username vault", // add your valid token here
				Password: "your password vault",
			},
			path: "efishery/data/service_payment_platform_mapping/staging",
		},
	}
	for _, tC := range testCases {
		svc, err := NewVault(tC.config)
		if err != nil {
			b.Error(err)
			return
		}
		for n := 0; n < b.N; n++ {
			for i := 0; i < 100; i++ {

				_, err := svc.Read(tC.path)
				if err != nil {
					b.Error(err)
				}
			}
		}

		defer svc.Close()
	}
}
