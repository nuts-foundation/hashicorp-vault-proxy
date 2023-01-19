package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "github.com/nuts-foundation/hashicorp-vault-proxy/api/v1"
	"github.com/nuts-foundation/hashicorp-vault-proxy/vault"
	"os"
)

const listenAddress = ":8210"

func main() {
	fmt.Println("Starting the Hashicorp Vault Proxy...")
	kv, err := vault.NewKVStore(loadConfig())
	if err != nil {
		panic(fmt.Errorf("unable to create Vault KVStore: %w", err))
	}
	handler := v1.NewStrictHandler(v1.NewWrapper(kv), nil)

	e := echo.New()
	loggerCfg := middleware.DefaultLoggerConfig
	loggerCfg.Skipper = func(c echo.Context) bool {
		return c.Path() == "/health"
	}
	e.Use(middleware.LoggerWithConfig(loggerCfg))
	e.HideBanner = true
	v1.RegisterHandlers(e, handler)
	err = e.Start(listenAddress)
	if err != nil {
		panic(fmt.Errorf("unable to start server: %w", err))
	}
	fmt.Println("Goodbye!")
}

func loadConfig() vault.Config {
	pathPrefix := os.Getenv("VAULT_PATHPREFIX")
	if pathPrefix == "" {
		pathPrefix = "kv"
	}
	cfg := vault.Config{
		PathPrefix: pathPrefix,
		Timeout:    0,
	}
	return cfg
}
