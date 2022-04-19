package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type VaultConfig struct {
	Host string `json:"host"`
	// Choose one authorization using token or username password.
	// Currently only support static username & password authentication
	Username string `json:"username"`
	Password string `json:"password"`
	// When using token will ignore username and password.
	Token string `json:"token"`
	// token ttl is ttl of token in second
	// if value 0 will set 30 days
	TokenTTL int `json:"token_ttl"`
	// disable auto renew token
	// by default token will renew before ttl expiration
	DisableAutoRenew bool `json:"disable_auto_renew"`

	// OnTokenRenew is hook when  auto renew token
	OnTokenRenew func(json.RawMessage, error)
}

type AuthResponse struct {
	ClientToken   string `json:"client_token"`
	LeaseDuration int    `json:"lease_duration"`
}

type VaultTokenResponse struct {
	Data AuthResponse `json:"data"`
}

type vaultService struct {
	config VaultConfig
	ctx    context.Context
	cancel context.CancelFunc
	ttl    int
}

// HTTP Helper
func (ss *vaultService) RequestWithContext(ctx context.Context, method, ep string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, ss.config.Host+ep, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if ss.config.Token != "" {
		req.Header.Add("X-Auth-Vault", ss.config.Token)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check error codes
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Status code %d , Body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (ss *vaultService) call(ctx context.Context, method, ep string, body []byte, v interface{}) error {
	resBody, err := ss.RequestWithContext(ctx, method, ep, body)
	if err != nil {
		return err

	}
	if err = json.Unmarshal(resBody, v); err != nil {
		return err
	}

	return nil
}

// End of http helper

// AuthWithContext is auth vault with passing context
func (ss *vaultService) AuthWithContext(ctx context.Context) (*VaultTokenResponse, error) {
	uri := "/auth/userpass"
	body := map[string]interface{}{
		"username": ss.config.Username,
		"password": ss.config.Password,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	var res VaultTokenResponse
	if err := ss.call(ctx, "POST", uri, payload, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// AuthWithContext is auth vault using back context
func (ss *vaultService) Auth() (*VaultTokenResponse, error) {
	return ss.AuthWithContext(context.Background())
}

// Renew Token renew
// parameter optional in second
func (ss *vaultService) RenewToken() error {
	return ss.RenewTokenWithContext(context.Background())
}

func (ss *vaultService) RenewTokenWithContext(ctx context.Context) error {
	// source
	// https://www.vaultproject.io/docs/concepts/lease#lease-durations-and-renewal
	response, err := ss.AuthWithContext(ctx)
	if err != nil {
		return err
	}

	// next renew 80% from lease duration
	ss.ttl = (response.Data.LeaseDuration * 80) / 100
	ss.config.Token = response.Data.ClientToken

	return nil
}

func (ss *vaultService) startAutoRenew() error {
	err := ss.RenewTokenWithContext(ss.ctx)
	if err != nil {
		return err
	}

	go func() {
		backoff := 0
		for {
			select {
			case <-ss.ctx.Done():
				return
			case <-time.After(time.Duration(ss.ttl) * time.Second):
				err := ss.RenewTokenWithContext(ss.ctx)

				if ss.config.OnTokenRenew != nil {
					ss.config.OnTokenRenew([]byte(ss.config.Token), err)
				}

				if err != nil {
					backoff++
					log.Println("Renew token error", err, "attempt", backoff)
					delay := (15 * backoff) + (60*2*backoff - 1)
					<-time.After(time.Duration(delay) * time.Second)
				}

				// timeout refresh logic
			case <-time.After(1 * time.Minute):
				continue
			}
		}
	}()

	return nil
}

type vaultValue map[string]interface{}

func (v vaultValue) Data() []byte {
	if v == nil || v["data"] == nil {
		return nil
	}

	resp, err := json.Marshal(v["data"])
	if err != nil {
		return nil
	}

	return resp
}

func (v vaultValue) Metadata() map[string]interface{} {
	if v == nil || v["metadata"] == nil {
		return nil
	}

	if vv, ok := v["metadata"].(map[string]interface{}); ok {
		return vv
	}

	return nil
}

// Read is read value from vault with passing context
func (ss *vaultService) ReadWithContext(ctx context.Context, identifier string) (Value, error) {
	var data map[string]interface{}
	err := ss.call(ctx, "GET", "/get/"+identifier, nil, &data)
	if err != nil {
		return vaultValue{}, err
	}
	return vaultValue(data), nil
}

// Read is read value from vault with context background
func (ss *vaultService) Read(identifier string) (Value, error) {
	return ss.ReadWithContext(context.Background(), identifier)
}

// Close periodic auto renew token
func (ss *vaultService) Close() error {
	ss.cancel()
	return nil
}

func NewVault(config VaultConfig) (*vaultService, error) {
	ctx := context.Background()
	ctxChild, cancelCtx := context.WithCancel(ctx)

	svc := &vaultService{
		config: config,
		ctx:    ctxChild,
		cancel: cancelCtx,
	}

	// when disable auto renew will bypass
	// periodic renew token
	if !config.DisableAutoRenew {
		err := svc.startAutoRenew()
		if err != nil {
			return nil, err
		}
	}

	return svc, nil
}
