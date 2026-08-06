package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISalesRepository interface {
	Create(ctx context.Context, sale *models.Sale) error
	CreateWithItems(ctx context.Context, sale *models.Sale, items []models.SaleItem) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Sale, error)
	FindByIDWithItems(ctx context.Context, id uuid.UUID) (*models.Sale, error)
	FindByOrderNumber(ctx context.Context, orderNumber string) (*models.Sale, error)
	FindAll(ctx context.Context, opts *base.QueryOptions) ([]models.Sale, int64, error)
	Update(ctx context.Context, sale *models.Sale) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.SaleStatus) error
	Delete(ctx context.Context, sale *models.Sale) error

	// Aggregation methods
	GetMonthlySales(ctx context.Context, year, month int) (*models.MonthlySales, error)
	GetMonthlySalesYear(ctx context.Context, year int) ([]models.MonthlySales, error)
	GetSalesTrend(ctx context.Context, startDate, endDate string) ([]models.DailySales, error)
	GetSalesBySource(ctx context.Context, startDate, endDate string) ([]models.SalesBySource, error)
	GetSalesByPackageType(ctx context.Context, startDate, endDate string) ([]models.SalesByPackageType, error)
	GetSalesByPackage(ctx context.Context, startDate, endDate string) ([]models.SalesByPackage, error)
	GetOverviewStats(ctx context.Context, startDate, endDate string) (*models.SalesOverviewStats, error)
	CountByStatus(ctx context.Context, status models.SaleStatus, startDate, endDate string) (int64, error)
	SumRevenue(ctx context.Context, startDate, endDate string) (float64, error)
}

type SalesRepository struct {
	*base.BaseRepository
}

func NewSalesRepository(baseRepo *base.BaseRepository) ISalesRepository {
	return &SalesRepository{BaseRepository: baseRepo}
}

// Create creates a new sale
func (r *SalesRepository) Create(ctx context.Context, sale *models.Sale) error {
	return r.BaseRepository.Create(ctx, sale)
}

// CreateWithItems creates a sale with its items in a transaction
func (r *SalesRepository) CreateWithItems(ctx context.Context, sale *models.Sale, items []models.SaleItem) error {
	return r.BaseRepository.WithTx(func(tx *gorm.DB) error {
		if err := tx.Create(sale).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].SaleID = sale.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByID finds a sale by ID
func (r *SalesRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Sale, error) {
	var sale models.Sale
	if err := r.BaseRepository.FindByID(ctx, &sale, id); err != nil {
		return nil, err
	}
	return &sale, nil
}

// FindByIDWithItems finds a sale by ID with items preloaded
func (r *SalesRepository) FindByIDWithItems(ctx context.Context, id uuid.UUID) (*models.Sale, error) {
	var sale models.Sale
	if err := r.DB.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&sale).Error; err != nil {
		return nil, err
	}
	return &sale, nil
}

// FindByOrderNumber finds a sale by order number
func (r *SalesRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*models.Sale, error) {
	var sale models.Sale
	if err := r.BaseRepository.FindOne(ctx, &sale, "order_number = ?", orderNumber); err != nil {
		return nil, err
	}
	return &sale, nil
}

// FindAll retrieves sales with pagination and filtering
func (r *SalesRepository) FindAll(ctx context.Context, opts *base.QueryOptions) ([]models.Sale, int64, error) {
	var sales []models.Sale
	count, err := r.BaseRepository.CountWithOptions(ctx, &models.Sale{}, opts)
	if err != nil {
		return nil, 0, err
	}
	if err := r.BaseRepository.FindMany(ctx, &models.Sale{}, &sales, opts); err != nil {
		return nil, 0, err
	}
	return sales, count, nil
}

// Update updates a sale
func (r *SalesRepository) Update(ctx context.Context, sale *models.Sale) error {
	return r.BaseRepository.Update(ctx, sale)
}

// UpdateStatus updates the status of a sale
func (r *SalesRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.SaleStatus) error {
	return r.BaseRepository.DB.WithContext(ctx).Model(&models.Sale{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete deletes a sale
func (r *SalesRepository) Delete(ctx context.Context, sale *models.Sale) error {
	return r.BaseRepository.Delete(ctx, sale)
}

// GetMonthlySales retrieves aggregated monthly sales data
func (r *SalesRepository) GetMonthlySales(ctx context.Context, year, month int) (*models.MonthlySales, error) {
	var result models.MonthlySales
	sql := `
		SELECT
			EXTRACT(YEAR FROM created_at)::int as year,
			EXTRACT(MONTH FROM created_at)::int as month,
			COUNT(*)::bigint as total_sales,
			COALESCE(SUM(total_amount), 0)::float8 as total_revenue,
			COALESCE(SUM(discount_amount), 0)::float8 as total_discount,
			COALESCE(SUM(CASE WHEN status = 'refunded' THEN final_amount ELSE 0 END), 0)::float8 as total_refunds,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as net_revenue,
			COALESCE(AVG(CASE WHEN status = 'completed' THEN final_amount END), 0)::float8 as avg_order_value,
			COUNT(CASE WHEN status = 'canceled' THEN 1 END)::bigint as canceled_sales,
			COUNT(CASE WHEN status = 'pending' THEN 1 END)::bigint as pending_sales,
			COUNT(CASE WHEN status = 'completed' THEN 1 END)::bigint as completed_sales
		FROM sales
		WHERE EXTRACT(YEAR FROM created_at) = $1 AND EXTRACT(MONTH FROM created_at) = $2
		GROUP BY EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, float64(year), float64(month)).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMonthlySalesYear retrieves aggregated monthly sales data for a year
func (r *SalesRepository) GetMonthlySalesYear(ctx context.Context, year int) ([]models.MonthlySales, error) {
	var results []models.MonthlySales
	sql := `
		SELECT
			EXTRACT(YEAR FROM created_at)::int as year,
			EXTRACT(MONTH FROM created_at)::int as month,
			COUNT(*)::bigint as total_sales,
			COALESCE(SUM(total_amount), 0)::float8 as total_revenue,
			COALESCE(SUM(discount_amount), 0)::float8 as total_discount,
			COALESCE(SUM(CASE WHEN status = 'refunded' THEN final_amount ELSE 0 END), 0)::float8 as total_refunds,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as net_revenue,
			COALESCE(AVG(CASE WHEN status = 'completed' THEN final_amount END), 0)::float8 as avg_order_value,
			COUNT(CASE WHEN status = 'canceled' THEN 1 END)::bigint as canceled_sales,
			COUNT(CASE WHEN status = 'pending' THEN 1 END)::bigint as pending_sales,
			COUNT(CASE WHEN status = 'completed' THEN 1 END)::bigint as completed_sales
		FROM sales
		WHERE EXTRACT(YEAR FROM created_at) = $1
		GROUP BY EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)
		ORDER BY EXTRACT(MONTH FROM created_at)
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, float64(year)).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// GetSalesTrend retrieves daily sales data for trend analysis
func (r *SalesRepository) GetSalesTrend(ctx context.Context, startDate, endDate string) ([]models.DailySales, error) {
	var results []models.DailySales
	sql := `
		SELECT
			DATE(created_at)::text as date,
			COUNT(*)::bigint as total_sales,
			COALESCE(SUM(total_amount), 0)::float8 as total_revenue,
			COALESCE(SUM(discount_amount), 0)::float8 as total_discount,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as net_revenue,
			COUNT(CASE WHEN status = 'completed' THEN 1 END)::bigint as completed_sales,
			COUNT(CASE WHEN status = 'refunded' THEN 1 END)::bigint as refunded_sales,
			COUNT(CASE WHEN status = 'pending' THEN 1 END)::bigint as pending_sales
		FROM sales
		WHERE created_at >= $1::date AND created_at <= $2::date
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// GetSalesBySource retrieves sales breakdown by source
func (r *SalesRepository) GetSalesBySource(ctx context.Context, startDate, endDate string) ([]models.SalesBySource, error) {
	var results []models.SalesBySource
	sql := `
		SELECT
			COALESCE(source, 'unknown') as source,
			COUNT(*)::bigint as total_sales,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as total_revenue,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as percentage
		FROM sales
		WHERE created_at >= $1::date AND created_at <= $2::date
		GROUP BY COALESCE(source, 'unknown')
		ORDER BY total_revenue DESC
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// GetSalesByPackageType retrieves sales breakdown by package type
func (r *SalesRepository) GetSalesByPackageType(ctx context.Context, startDate, endDate string) ([]models.SalesByPackageType, error) {
	var results []models.SalesByPackageType
	sql := `
		SELECT
			COALESCE(package_type, 'unknown') as package_type,
			COUNT(*)::bigint as total_sales,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as total_revenue,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as percentage
		FROM sales
		WHERE created_at >= $1::date AND created_at <= $2::date
		GROUP BY COALESCE(package_type, 'unknown')
		ORDER BY total_revenue DESC
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// GetSalesByPackage retrieves sales breakdown by individual package (from sale_items)
func (r *SalesRepository) GetSalesByPackage(ctx context.Context, startDate, endDate string) ([]models.SalesByPackage, error) {
	var results []models.SalesByPackage
	sql := `
		SELECT
			si.package_id as package_id,
			COALESCE(si.package_name, 'Unknown Package') as package_name,
			COUNT(DISTINCT s.id)::bigint as total_sales,
			COALESCE(SUM(si.quantity), 0)::bigint as total_quantity,
			COALESCE(SUM(CASE WHEN s.status = 'completed' THEN si.subtotal ELSE 0 END), 0)::float8 as total_revenue,
			COALESCE(AVG(si.unit_price), 0)::float8 as avg_unit_price
		FROM sale_items si
		JOIN sales s ON s.id = si.sale_id
		WHERE s.created_at >= $1::date AND s.created_at <= $2::date
		GROUP BY si.package_id, si.package_name
		ORDER BY total_revenue DESC
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// GetOverviewStats retrieves overview statistics for sales
func (r *SalesRepository) GetOverviewStats(ctx context.Context, startDate, endDate string) (*models.SalesOverviewStats, error) {
	var result models.SalesOverviewStats
	sql := `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as total_revenue,
			COUNT(CASE WHEN status = 'completed' THEN 1 END)::bigint as total_sales,
			COALESCE(SUM(CASE WHEN status = 'refunded' THEN final_amount ELSE 0 END), 0)::float8 as total_refunds,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN final_amount ELSE 0 END), 0)::float8 as net_revenue,
			COALESCE(AVG(CASE WHEN status = 'completed' THEN final_amount END), 0)::float8 as avg_order_value,
			COUNT(CASE WHEN status = 'completed' THEN 1 END)::bigint as completed_sales,
			COUNT(CASE WHEN status = 'pending' THEN 1 END)::bigint as pending_sales,
			COUNT(CASE WHEN status = 'canceled' THEN 1 END)::bigint as canceled_sales,
			COUNT(CASE WHEN status = 'refunded' THEN 1 END)::bigint as refunded_sales
		FROM sales
		WHERE created_at >= $1::date AND created_at <= $2::date
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, startDate, endDate).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// CountByStatus counts sales by status
func (r *SalesRepository) CountByStatus(ctx context.Context, status models.SaleStatus, startDate, endDate string) (int64, error) {
	var count int64
	sql := `
		SELECT COUNT(*)::bigint
		FROM sales
		WHERE status = $1 AND created_at >= $2::date AND created_at <= $3::date
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, status, startDate, endDate).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SumRevenue calculates total revenue for a date range
func (r *SalesRepository) SumRevenue(ctx context.Context, startDate, endDate string) (float64, error) {
	var sum float64
	sql := `
		SELECT COALESCE(SUM(final_amount), 0)::float8
		FROM sales
		WHERE status = 'completed' AND created_at >= $1::date AND created_at <= $2::date
	`
	if err := r.BaseRepository.DB.WithContext(ctx).Raw(sql, startDate, endDate).Scan(&sum).Error; err != nil {
		return 0, err
	}
	return sum, nil
}
