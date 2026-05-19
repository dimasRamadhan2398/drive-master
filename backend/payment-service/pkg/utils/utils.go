package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("ORD-%s-%s", timestamp, randomHex)
}

func GeneratePaymentCode() string {
	timestamp := time.Now().Format("20060102")
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("PAY%s%s", timestamp, strings.ToUpper(randomHex))
}

func FormatCurrency(amount float64, currency string) string {
	return fmt.Sprintf("%s %.2f", currency, amount)
}

func IsExpired(expiryTime time.Time) bool {
	return time.Now().After(expiryTime)
}

func ParseTime(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}