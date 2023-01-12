package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "github.com/nuts-foundation/hashicorp-vault-proxy/api/v1"
	"github.com/nuts-foundation/hashicorp-vault-proxy/vault"
)

func main() {
	fmt.Println("Starting the Hashicorp Vault Proxy...")
	kv, err := vault.NewKVStore(vault.Config{
		Token:      "unsafe",
		Address:    "http://localhost:8200",
		PathPrefix: "kv",
		Timeout:    0,
	})
	if err != nil {
		panic(fmt.Errorf("unable to create Vault KVStore: %w", err))
	}
	handler := v1.NewStrictHandler(v1.NewWrapper(kv), nil)

	e := echo.New()
	e.Use(middleware.Logger())
	e.HideBanner = true
	v1.RegisterHandlers(e, handler)
	err = e.Start(":7863")
	if err != nil {
		panic(fmt.Errorf("unable to start server: %w", err))
	}
	fmt.Println("Done!")
}
