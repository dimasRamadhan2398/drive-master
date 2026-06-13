package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// IUserClient defines the interface for user-service operations
type IUserClient interface {
	GetInstructorRecurringSchedules(ctx context.Context, instructorID uuid.UUID) ([]RecurringScheduleResponse, error)
	GetAllInstructors(ctx context.Context) ([]InstructorWithSchedules, error)
	GetUserByID(ctx context.Context, userID uint) (*UserResponse, error)
}

// UserClient implements IUserClient
type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

// InstructorWithSchedules represents an instructor with their recurring schedules
type InstructorWithSchedules struct {
	ID                 uuid.UUID                  `json:"id"`
	FirstName          string                    `json:"firstName"`
	LastName           string                    `json:"lastName"`
	Email              string                    `json:"email"`
	RecurringSchedules []RecurringScheduleDTO    `json:"recurringSchedules"`
}

// RecurringScheduleDTO represents a recurring schedule for API calls
type RecurringScheduleDTO struct {
	ID           uuid.UUID `json:"id"`
	InstructorID uuid.UUID `json:"instructorId"`
	DayOfWeek    int       `json:"dayOfWeek"`
	DayName      string    `json:"dayName"`
	StartTime    string    `json:"startTime"`
	EndTime      string    `json:"endTime"`
	IsActive     bool      `json:"isActive"`
}

// NewUserClient creates a new user service client
func NewUserClient(baseURL string) IUserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetInstructorRecurringSchedules retrieves all active recurring schedules for an instructor
func (c *UserClient) GetInstructorRecurringSchedules(ctx context.Context, instructorID uuid.UUID) ([]RecurringScheduleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/instructors/%s/recurring-schedules", c.baseURL, instructorID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var schedules []RecurringScheduleResponse
	if err := json.Unmarshal(apiResp.Data, &schedules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedules: %w", err)
	}

	return schedules, nil
}

// GetAllInstructors retrieves all instructors with their recurring schedules
func (c *UserClient) GetAllInstructors(ctx context.Context) ([]InstructorWithSchedules, error) {
	url := fmt.Sprintf("%s/api/v1/instructors/with-schedules", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var instructors []InstructorWithSchedules
	if err := json.Unmarshal(apiResp.Data, &instructors); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instructors: %w", err)
	}

	return instructors, nil
}

// doRequest is a helper for making HTTP requests (not currently used but available for future endpoints)
func (c *UserClient) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = io.NopCloser(bytes.NewReader(jsonData))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.httpClient.Do(req)
}

// UserResponse represents a user from user-service
type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
}

// GetUserByID retrieves a user by ID from user-service
func (c *UserClient) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	var user UserResponse
	if err := json.Unmarshal(apiResp.Data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}