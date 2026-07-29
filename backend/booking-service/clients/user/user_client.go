package user

import (
	"booking-service/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// IUserClient defines the interface for user-service operations
type IUserClient interface {
	GetInstructorRecurringSchedules(ctx context.Context, instructorID uuid.UUID) ([]RecurringScheduleResponse, error)
	GetAllInstructors(ctx context.Context) ([]InstructorWithSchedules, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
	GetEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) (*EntitlementInfo, error)
	GetMemberEntitlements(ctx context.Context, memberID uuid.UUID) ([]EntitlementInfo, error)
	CreateEntitlement(ctx context.Context, memberID uuid.UUID, input CreateEntitlementInput) error
	UseSession(ctx context.Context, memberID, entitlementID uuid.UUID) error
	RateInstructor(ctx context.Context, instructorID uuid.UUID, rating float64) error
}

// UserClient implements IUserClient
type UserClient struct {
	baseURL    string
	httpClient *http.Client
	jwtSecret  string
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
func NewUserClient(baseURL string, jwtSecret string) IUserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		jwtSecret: jwtSecret,
	}
}

type Claims struct {
	User *UserInfo
	jwt.RegisteredClaims
}


// generateServiceToken generates a JWT token for service-to-service communication
func (c *UserClient) generateServiceToken() (string, error) {
	uuidResult, err := uuid.Parse("cf475ead-91bf-4a55-a47e-cf93279240b6")
	if err != nil {
		return "", fmt.Errorf("failed to parse user ID: %w", err)
	}
    claims := &Claims{
        User: &UserInfo{
            ID: uuidResult,
			Email: "admin@example.com",			
        },
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
        },
    }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.jwtSecret))
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

	// Add JWT token for service-to-service authentication
	token, err := c.generateServiceToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate service token: %w", err)
	}

	logger.Info("token adalah",logger.LogField("token", token))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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

	// Add JWT token for service-to-service authentication
	token, err := c.generateServiceToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate service token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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


// GetUserByID retrieves a user by ID from user-service
func (c *UserClient) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add JWT token for service-to-service authentication
	token, err := c.generateServiceToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate service token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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

	var user UserInfo
	if err := json.Unmarshal(apiResp.Data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetEntitlement retrieves an entitlement by member ID and entitlement ID from user-service
func (c *UserClient) GetEntitlement(ctx context.Context, memberID, entitlementID uuid.UUID) (*EntitlementInfo, error) {
	url := fmt.Sprintf("%s/api/v1/entitlements/members/%s/entitlements/%s", c.baseURL, memberID.String(), entitlementID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add JWT token for service-to-service authentication
	token, err := c.generateServiceToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate service token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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

	var entitlement EntitlementInfo
	if err := json.Unmarshal(apiResp.Data, &entitlement); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entitlement: %w", err)
	}

	return &entitlement, nil
}

// GetMemberEntitlements retrieves all entitlements of a member from user-service
func (c *UserClient) GetMemberEntitlements(ctx context.Context, memberID uuid.UUID) ([]EntitlementInfo, error) {
	url := fmt.Sprintf("%s/api/v1/entitlements/members/%s?limit=100", c.baseURL, memberID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add JWT token for service-to-service authentication
	token, err := c.generateServiceToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate service token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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

	var list struct {
		Data []EntitlementInfo `json:"data"`
	}
	if err := json.Unmarshal(apiResp.Data, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entitlements list: %w", err)
	}

	return list.Data, nil
}

// CreateEntitlementInput is used for POST /members/:id/entitlements
type CreateEntitlementInput struct {
	BookingID        uuid.UUID `json:"bookingId"`
	PackageID        uuid.UUID `json:"packageId"`
	PackageName      string    `json:"packageName"`
	TotalSessions    int       `json:"totalSessions"`
	IsNightSession   bool      `json:"isNightSession"`
	IsWeekendSession bool      `json:"isWeekendSession"`
	StartDate        string    `json:"startDate"`
	EndDate          string    `json:"endDate,omitempty"`
}

// CreateEntitlement creates a new entitlement in user-service
func (c *UserClient) CreateEntitlement(ctx context.Context, memberID uuid.UUID, input CreateEntitlementInput) error {
	url := fmt.Sprintf("%s/api/v1/entitlements/members/%s", c.baseURL, memberID.String())

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal create entitlement input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	token, err := c.generateServiceToken()
	if err != nil {
		return fmt.Errorf("failed to generate service token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user-service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// UseSession decrements remaining sessions and increments used sessions for a member entitlement
func (c *UserClient) UseSession(ctx context.Context, memberID, entitlementID uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/entitlements/members/%s/entitlements/%s/use-session", c.baseURL, memberID.String(), entitlementID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	token, err := c.generateServiceToken()
	if err != nil {
		return fmt.Errorf("failed to generate service token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call user-service use-session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user-service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// RateInstructor sends a session rating for an instructor to update their average rating
func (c *UserClient) RateInstructor(ctx context.Context, instructorID uuid.UUID, rating float64) error {
	url := fmt.Sprintf("%s/api/v1/instructors/%s/rate", c.baseURL, instructorID.String())

	bodyData := map[string]interface{}{"rating": rating}
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		return fmt.Errorf("failed to marshal rate payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	token, err := c.generateServiceToken()
	if err != nil {
		return fmt.Errorf("failed to generate service token: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call user-service rate instructor: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user-service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}