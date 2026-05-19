package main

import (
	"os"
	"payment-service/cmd/cli"
)

// @title           Payment Service API
// @version         1.0
// @description     API documentation for Payment Service - manages payments, transactions, and payment methods
// @host            localhost:8004
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