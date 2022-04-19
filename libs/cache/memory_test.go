package cache

import (
	"testing"
	"time"
)

func TestMemory(t *testing.T) {
	client := NewMemoryCache(MemoryConfig{TTL: 60})
	defer client.Close()
	testCases := []struct {
		desc    string
		data    interface{}
		success bool
	}{
		{
			desc:    "Test using byte",
			data:    map[string]interface{}{"k": "v"},
			success: true,
		},
		{
			desc:    "Test using byte",
			data:    []byte(`{"k":"v"}`),
			success: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := client.Set(tC.desc, tC.data)
			if err != nil && tC.success {
				t.Error(err)
				return
			}

			if !client.IsExists(tC.desc) {
				t.Error("Data Should exist")
				return

			}

			var data map[string]interface{}
			err = client.Get(tC.desc, &data)
			if err != nil && tC.success {
				t.Error(err)
				return

			}

			err = client.Delete(tC.desc)
			if err != nil && tC.success {
				t.Error(err)
				return

			}

			if client.IsExists(tC.desc) {
				t.Error("Data Should not exist")
				return

			}

		})
	}
}

func TestMemoryExpiration(t *testing.T) {
	client := NewMemoryCache(MemoryConfig{TTL: 3})
	defer client.Close()
	testCases := []struct {
		desc    string
		data    interface{}
		success bool
	}{
		{
			desc:    "Test using byte",
			data:    map[string]interface{}{"k": "v"},
			success: true,
		},
		{
			desc:    "Test using byte",
			data:    []byte(`{"k":"v"}`),
			success: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := client.Set(tC.desc, tC.data)
			if err != nil && tC.success {
				t.Error(err)
				return
			}

			if !client.IsExists(tC.desc) {
				t.Error("Data Should exist")
				return

			}

			<-time.After(5 * time.Second)

			if client.IsExists(tC.desc) {
				t.Error("Data Should not exist")
				return

			}

		})
	}
}

func TestMemoryNotFound(t *testing.T) {
	client := NewMemoryCache(MemoryConfig{TTL: 60})
	defer client.Close()
	testCases := []struct {
		desc    string
		data    interface{}
		success bool
	}{
		{
			desc:    "Test using byte",
			data:    map[string]interface{}{"k": "v"},
			success: false,
		},
		{
			desc:    "Test using byte",
			data:    []byte(`{"k":"v"}`),
			success: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var data map[string]interface{}
			err := client.Get(tC.desc, &data)
			if err == nil && tC.success {
				t.Error(err)
				return

			}

		})
	}
}
