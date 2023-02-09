/*
 * Copyright (C) 2023 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
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
