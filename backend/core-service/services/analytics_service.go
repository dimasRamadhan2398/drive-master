package services

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"core-service/pkg/config"
	"core-service/pkg/logger"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

type GAOverview struct {
	Date      string `json:"date"`
	Users     int64  `json:"users"`
	PageViews int64  `json:"pageviews"`
}

type GAFunnelStep struct {
	EventName string `json:"event_name"`
	Count     int64  `json:"count"`
}

type GATrafficSource struct {
	Source   string `json:"source"`
	Medium   string `json:"medium"`
	Sessions int64  `json:"sessions"`
}

type IAnalyticsService interface {
	GetOverviewReport(ctx context.Context, startDate, endDate string) ([]GAOverview, error)
	GetFunnelReport(ctx context.Context) ([]GAFunnelStep, error)
}

type AnalyticsService struct {
	analyticsService *analyticsdata.Service
	propertyID       string
	isMock           bool
}

func NewAnalyticsService() IAnalyticsService {
	cfg := config.Get();
	measurementID := cfg.Analytics.GA4MeasurementID
	propertyID := cfg.Analytics.GA4PropertyID
	credentialsFile := cfg.Analytics.GA4CredentialsFile

	if measurementID == "" || propertyID == "" {
		logger.Warn("GA4_MEASUREMENT_ID or GA4_PROPERTY_ID environment variables are not set. Core-service will run analytics in MOCK mode with simulated data.")
		return &AnalyticsService{
			isMock: true,
		}
	}

	// If credentials file is not set, fall back to mock mode
	if credentialsFile == "" {
		logger.Warn("GA4_CREDENTIALS_FILE is not set. Core-service will run analytics in MOCK mode.")
		return &AnalyticsService{
			isMock: true,
		}
	}

	// Check if credentials file exists
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		logger.Warn("GA4 credentials file does not exist at: " + credentialsFile + ". Core-service will run analytics in MOCK mode.", logger.NewLogField("error", err))
		return &AnalyticsService{
			isMock: true,
		}
	}

	ctx := context.Background()
	srv, err := analyticsdata.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		logger.Error("Failed to initialize Google Analytics API service from file. Falling back to MOCK mode.", logger.NewLogField("error", err))
		return &AnalyticsService{
			isMock: true,
		}
	}

	logger.Info("Google Analytics 4 service successfully initialized in core-service.")
	return &AnalyticsService{
		analyticsService: srv,
		propertyID:       propertyID,
		isMock:           false,
	}
}

func (s *AnalyticsService) GetOverviewReport(ctx context.Context, startDate, endDate string) ([]GAOverview, error) {
	if s.isMock {
		return s.generateMockOverview(startDate, endDate), nil
	}

	req := &analyticsdata.RunReportRequest{
		Property: "properties/" + s.propertyID,
		DateRanges: []*analyticsdata.DateRange{
			{StartDate: startDate, EndDate: endDate},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "activeUsers"},
			{Name: "screenPageViews"},
		},
		Dimensions: []*analyticsdata.Dimension{
			{Name: "date"},
		},
		OrderBys: []*analyticsdata.OrderBy{
			{Dimension: &analyticsdata.DimensionOrderBy{DimensionName: "date"}},
		},
	}

	resp, err := s.analyticsService.Properties.RunReport(req.Property, req).Do()
	if err != nil {
		return nil, err
	}

	var results []GAOverview
	for _, row := range resp.Rows {
		dateStr := row.DimensionValues[0].Value
		usersVal := row.MetricValues[0].Value
		pvsVal := row.MetricValues[1].Value

		var users, pvs int64
		fmt.Sscanf(usersVal, "%d", &users)
		fmt.Sscanf(pvsVal, "%d", &pvs)

		parsedTime, err := time.Parse("20060102", dateStr)
		formattedDate := dateStr
		if err == nil {
			formattedDate = parsedTime.Format("2006-01-02")
		}

		results = append(results, GAOverview{
			Date:      formattedDate,
			Users:     users,
			PageViews: pvs,
		})
	}
	return results, nil
}

func (s *AnalyticsService) GetFunnelReport(ctx context.Context) ([]GAFunnelStep, error) {
	if s.isMock {
		return s.generateMockFunnel(), nil
	}

	req := &analyticsdata.RunReportRequest{
		Property: "properties/" + s.propertyID,
		DateRanges: []*analyticsdata.DateRange{
			{StartDate: "30daysAgo", EndDate: "today"},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "eventCount"},
		},
		Dimensions: []*analyticsdata.Dimension{
			{Name: "eventName"},
		},
	}

	resp, err := s.analyticsService.Properties.RunReport(req.Property, req).Do()
	if err != nil {
		return nil, err
	}

	eventCounts := map[string]int64{
		"page_view":      0,
		"view_item":      0,
		"begin_checkout": 0,
		"purchase":       0,
	}

	for _, row := range resp.Rows {
		name := row.DimensionValues[0].Value
		countVal := row.MetricValues[0].Value
		var count int64
		fmt.Sscanf(countVal, "%d", &count)

		if _, exists := eventCounts[name]; exists {
			eventCounts[name] = count
		}
	}

	steps := []string{"page_view", "view_item", "begin_checkout", "purchase"}
	var results []GAFunnelStep
	for _, step := range steps {
		results = append(results, GAFunnelStep{
			EventName: step,
			Count:     eventCounts[step],
		})
	}
	return results, nil
}

func (s *AnalyticsService) generateMockOverview(startDate, endDate string) []GAOverview {
	var results []GAOverview
	// Parsir tanggal mulai & akhir untuk menghasilkan range data (default 30 hari terakhir)
	days := 30
	now := time.Now()

	// Mulai dari 30 hari yang lalu hingga hari ini
	for i := days; i >= 0; i-- {
		t := now.AddDate(0, 0, -i)
		dateStr := t.Format("2006-01-02")
		
		// generate realistic fluctuating data
		r := rand.New(rand.NewSource(t.UnixNano()))
		users := int64(150 + r.Intn(100))
		pvs := users * int64(2 + r.Intn(3))

		results = append(results, GAOverview{
			Date:      dateStr,
			Users:     users,
			PageViews: pvs,
		})
	}
	return results
}

func (s *AnalyticsService) generateMockFunnel() []GAFunnelStep {
	// Simulated conversion numbers
	return []GAFunnelStep{
		{EventName: "page_view", Count: 1450},
		{EventName: "view_item", Count: 920},
		{EventName: "begin_checkout", Count: 310},
		{EventName: "purchase", Count: 120},
	}
}
