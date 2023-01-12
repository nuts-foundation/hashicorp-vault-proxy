package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	v1 "github.com/nuts-foundation/hashicorp-vault-proxy/api/v1"
)

func main() {
	fmt.Println("Starting the Hashicorp Vault Proxy...")
	e := echo.New()
	handler := v1.NewStrictHandler(v1.Wrapper{}, nil)
	v1.RegisterHandlers(e, handler)
	e.Start(":7863")
	fmt.Println("Done!")
}
