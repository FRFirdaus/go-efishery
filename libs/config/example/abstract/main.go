package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"bitbucket.org/efishery/go-efishery/libs/config"
)

func main() {
	ctx := context.Background()
	var vaultaddress = "https://jambal.service.efishery.com/vault"

	vaultConfig := config.VaultConfig{
		Host:     vaultaddress,
		Username: "your username vault", // add your valid token here
		Password: "your password vault",
		OnTokenRenew: func(rm json.RawMessage, e error) {
			// add logic on token renew here
			// e.g when error hook slack notif, hook whatsapp notif etc
			fmt.Println("is error", e)
			fmt.Println("resp", string(rm))
		},
	}

	vaultService, err := config.NewVault(vaultConfig)
	if err != nil {
		log.Panic(err)
	}

	// init for global reusable function
	// and high level abstraction
	config.InitRemoteConfig(vaultService)
	// close auto renew
	defer config.CloseRemoteConfig()

	// call this in any package/module
	resp, err := config.RemoteConfig().ReadWithContext(ctx, "efishery/data/service_payment_platform_mapping/staging")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("resp", string(resp.Data()))
}
