package main

import (
	"os"
	"user-service/cmd/cli"
)

// @title           User Service API
// @version         1.0
// @description     API documentation for User Service
// @host            localhost:8001
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey XApiKey
// @in header
// @name x-api-key

// @securityDefinitions.apikey XRequestAt
// @in header
// @name x-request-at

// @securityDefinitions.apikey XServiceName
// @in header
// @name x-service-name
func main() {
    if err := cli.Execute(); err != nil {
        os.Exit(1)
    }
}
