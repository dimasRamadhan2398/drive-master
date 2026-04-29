package main

import (
	"os"
	"user-service/cmd/cli"
)

// @title           User Service API
// @version         1.0
// @description     API documentation for User Service
// @host            localhost:8081
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
    if err := cli.Execute(); err != nil {
        os.Exit(1)
    }
}
