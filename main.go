package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "github.com/nuts-foundation/hashicorp-vault-proxy/api/v1"
	"github.com/nuts-foundation/hashicorp-vault-proxy/vault"
	"github.com/sirupsen/logrus"
	"os"
)

const listenAddress = ":8210"

func main() {
	logrus.Info("Starting the Hashicorp Vault Proxy...")
	pathPrefix := os.Getenv("VAULT_PATHPREFIX")
	if pathPrefix == "" {
		pathPrefix = "kv"
	}

	kv, err := vault.NewKVStore(pathPrefix)
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
	logrus.Info("Goodbye!")
}
