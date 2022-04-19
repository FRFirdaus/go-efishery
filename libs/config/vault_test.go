package config

import (
	"fmt"
	"testing"
)

var vaultaddress = "https://jambal.service.efishery.com/vault"

func TestGetVault(t *testing.T) {
	testCases := []struct {
		desc    string
		config  VaultConfig
		path    string
		success bool
	}{
		{
			desc: "Testing get valid token",
			config: VaultConfig{
				Host:     vaultaddress,
				Username: "your username vault", // add your valid token here
				Password: "your password vault",
			},
			path:    "efishery/data/service_payment_platform_mapping/staging",
			success: true,
		},
		{
			desc: "Testing get invalid path",
			config: VaultConfig{
				Host:     vaultaddress,
				Username: "your username vault", // add your valid token here
				Password: "your password vault",
			},
			path:    "efishery/data/service_payment_platform_mapping/testos",
			success: false,
		},
		{
			desc: "Testing using invalid token",
			config: VaultConfig{
				Host:     vaultaddress,
				Username: "your username vault", // add your valid token here
				Password: "your password vault",
			},
			path: "efishery/data/service_payment_platform_mapping/staging",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			svc, err := NewVault(tC.config)
			if err != nil && tC.success {
				t.Error(err)
				return
			}
			if svc == nil {
				return
			}
			defer svc.Close()

			data, err := svc.Read(tC.path)
			if err != nil && tC.success {
				t.Error(err)
				return
			}

			fmt.Println(string(data.Data()))
		})
	}
}

// Using jambal cannot using with test case
// func TestAutoRenewToken(t *testing.T) {
// 	testCases := []struct {
// 		desc    string
// 		config  VaultConfig
// 		success bool
// 	}{
// 		{
// 			desc: "Testing Auto Renew Token",
// 			config: VaultConfig{
// 				Host:     vaultaddress,
// 				Username: "your username vault", // add your valid token here
// 				Password: "your password vault",
// 			},
// 			success: true,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			var success bool
// 			tC.config.OnTokenRenew = func(rm json.RawMessage, err error) {
// 				fmt.Println(string(rm))
// 				if err == nil {
// 					success = true
// 				}
// 			}
// 			svc, err := NewVault(tC.config)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			defer svc.Close()

// 			// wait ttl
// 			<-time.After(time.Second * time.Duration(tC.config.TokenTTL))
// 			if success != tC.success {
// 				t.Errorf("Expected %v, got %v", tC.success, success)
// 			}
// 			// reset to long live token
// 			svc.RenewToken()
// 		})
// 	}
// }
