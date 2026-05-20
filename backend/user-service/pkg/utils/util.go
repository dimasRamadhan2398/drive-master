package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"

	"github.com/dustin/go-humanize"
)

type PaginationParam struct {
	Count int64       `json:"count"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}

type PaginationResult struct {
	TotalPage    int         `json:"totalPage"`
	TotalData    int64       `json:"totalData"`
	NextPage     *int        `json:"nextPage"`
	PreviousPage *int        `json:"previousPage"`
	Page         int         `json:"page"`
	Limit        int         `json:"limit"`
	Data         interface{} `json:"data"`
}

func GeneratePagination(params PaginationParam) PaginationResult {
	totalPage := int(math.Ceil(float64(params.Count) / float64(params.Limit)))

	var (
		nextPage     int
		previousPage int
	)
	if params.Page < totalPage {
		nextPage = params.Page + 1
	}

	if params.Page > 1 {
		previousPage = params.Page - 1
	}

	result := PaginationResult{
		TotalPage:    totalPage,
		TotalData:    params.Count,
		NextPage:     &nextPage,
		PreviousPage: &previousPage,
		Page:         params.Page,
		Limit:        params.Limit,
		Data:         params.Data,
	}
	return result
}

func GenerateSHA256(inputString string) string {
	hash := sha256.New()
	hash.Write([]byte(inputString))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func RupiahFormat(amount *float64) string {
	stringValue := "0"
	if amount != nil {
		humanizeValue := humanize.CommafWithDigits(*amount, 0)
		stringValue = strings.ReplaceAll(humanizeValue, ",", ".")
	}
	return fmt.Sprintf("Rp. %s", stringValue)
}