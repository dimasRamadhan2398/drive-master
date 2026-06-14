package base

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUintIDFromPath extracts a uint ID from the URL path parameter
func GetUintIDFromPath(ctx *gin.Context, param string) (uint, error) {
	idStr := ctx.Param(param)
	if idStr == "" {
		return 0, errors.New("parameter not found")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid id format")
	}

	return uint(id), nil
}
