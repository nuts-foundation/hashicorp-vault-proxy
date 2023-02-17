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
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	v1 "github.com/nuts-foundation/hashicorp-vault-proxy/api/v1"
	"github.com/nuts-foundation/hashicorp-vault-proxy/vault"
)

const listenAddress = ":8210"

func main() {
	logFormat := os.Getenv("LOG_FORMAT")
	switch logFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	logrus.Infof("Starting the Hashicorp Vault Proxy on %s", listenAddress)

	// pathPrefix should always be set
	pathPrefix := os.Getenv("VAULT_PATHPREFIX")
	if pathPrefix == "" {
		pathPrefix = "kv"
	}

	// pathName is optional
	pathName, isSet := os.LookupEnv("VAULT_PATHNAME")
	if !isSet {
		pathName = "nuts-private-keys"
	}

	// JoinPath will only add a slash if pathName is set
	path, err := url.JoinPath(pathPrefix, pathName)
	if err != nil {
		panic(fmt.Errorf("unable to assemble vault secret path: %w", err))
	}

	kv, err := vault.NewKVStore(path)
	if err != nil {
		panic(fmt.Errorf("unable to create Vault KVStore: %w", err))
	}

	handler := v1.NewStrictHandler(v1.NewWrapper(kv), nil)

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health"
		},
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogRemoteIP: true,
		LogError:    true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			status := values.Status
			if values.Error != nil {
				switch errWithStatus := values.Error.(type) {
				case *echo.HTTPError:
					status = errWithStatus.Code
				default:
					status = http.StatusInternalServerError
				}
			}

			logrus.WithFields(logrus.Fields{
				"remote_ip": values.RemoteIP,
				"method":    values.Method,
				"uri":       values.URI,
				"status":    status,
			}).Info("HTTP request")

			return nil
		},
	}))
	e.HideBanner = true
	e.HidePort = true
	v1.RegisterHandlers(e, handler)
	err = e.Start(listenAddress)
	if err != nil {
		panic(fmt.Errorf("unable to start server: %w", err))
	}
	logrus.Info("Goodbye!")
}
