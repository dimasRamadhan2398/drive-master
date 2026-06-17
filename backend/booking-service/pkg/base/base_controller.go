package base

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// GetUUIDIDFromPath extracts a uuid.UUID ID from the URL path parameter
func GetUUIDIDFromPath(ctx *gin.Context, param string) (uuid.UUID, error) {
	idStr := ctx.Param(param)
	if idStr == "" {
		return uuid.Nil, errors.New("parameter not found")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid id format")
	}

	return id, nil
}
