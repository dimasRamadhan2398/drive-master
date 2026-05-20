package main

import (
	"core-service/cmd/cli"
	"os"
)

// @title           Core Service API
// @version         1.0
// @description     API documentation for Core Service - manages cars, packages, and regions
// @host            localhost:8002
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
