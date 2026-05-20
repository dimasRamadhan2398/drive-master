package utils

import (
	"strconv"
	"time"

	apperrors "user-service/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// String to Uint
func StringToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	return uint(val), err
}

func StringToUintPtr(s string) *uint {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}
	uintVal := uint(val)
	return &uintVal
}

// Uint to String
func UintToString(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

func UintToStringPtr(i *uint) string {
	if i == nil {
		return ""
	}
	return strconv.FormatUint(uint64(*i), 10)
}

// String to Int
func StringToInt(s string) (int, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	return int(val), err
}

func StringToIntPtr(s string) *int {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	intVal := int(val)
	return &intVal
}

// Int to String
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// String to Float64
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func StringToFloat64Ptr(s string) *float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &val
}

// Float64 to String
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// String to Bool
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func StringToBoolPtr(s string) *bool {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &val
}

// Bool to String
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// String to UUID
func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func StringToUUIDPtr(s string) *uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

// UUID to String
func UUIDToString(id uuid.UUID) string {
	return id.String()
}

func UUIDToStringPtr(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

// String to Time (various formats)
func StringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func StringToTimePtr(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}

func StringToTimeCustom(s, layout string) (time.Time, error) {
	return time.Parse(layout, s)
}

// Time to String (various formats)
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func TimeToStringCustom(t time.Time, layout string) string {
	return t.Format(layout)
}

// Time to Date only (YYYY-MM-DD)
func TimeToDateString(t time.Time) string {
	return t.Format("2006-01-02")
}

// Pointer to value helpers
func PtrToVal[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

func ValToPtr[T any](val T) *T {
	return &val
}

// Uint to Int
func UintToInt(u uint) int {
	return int(u)
}

// Int to Uint
func IntToUint(i int) uint {
	return uint(i)
}

// Uint64 to Uint
func Uint64ToUint(u uint64) uint {
	return uint(u)
}

// Uint to Uint64
func UintToUint64(u uint) uint64 {
	return uint64(u)
}

// String to Uint64
func StringToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// Uint64 to String
func Uint64ToString(u uint64) string {
	return strconv.FormatUint(u, 10)
}

// Bytes to String
func BytesToString(b []byte) string {
	return string(b)
}

// String to Bytes
func StringToBytes(s string) []byte {
	return []byte(s)
}

// Safe string operations
func StringPtrToString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func IntPtrToInt(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func UintPtrToUint(ptr *uint) uint {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// Helper to get user ID from path (returns UUID)
func getUserIDFromPath(c *gin.Context, paramName string) (uuid.UUID, error) {
	idStr := c.Param(paramName)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, apperrors.ErrBadRequest
	}
	return id, nil
}

// Helper to get uint ID from path
func getUintIDFromPath(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, apperrors.ErrBadRequest
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, apperrors.ErrBadRequest
	}
	return uint(id), nil
}