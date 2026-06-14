package base

import (
	"context"
	"fmt"
	"strings"

	apperrors "core-service/pkg/errors"

	"gorm.io/gorm"
)

// QueryOptions holds common query parameters for filtering, searching, and pagination
type QueryOptions struct {
	Where        map[string]any // Key-value pairs for WHERE conditions
	Search       string         // Search term for ILIKE queries
	SearchFields []string       // Fields to search in when Search is provided
	Offset       int            // Number of records to skip (pagination)
	Limit        int            // Maximum number of records to return
	Order        string         // Order clause, e.g., "created_at DESC"
	Preloads     []string       // Relationships to preload
}

// NewQueryOptions creates a QueryOptions with default values
func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Limit:  10,
		Offset: 0,
		Order:  "created_at DESC",
	}
}

// WithWhere sets the Where condition
func (o *QueryOptions) WithWhere(where map[string]any) *QueryOptions {
	o.Where = where
	return o
}

// WithSearch sets the search term and fields
func (o *QueryOptions) WithSearch(search string, fields ...string) *QueryOptions {
	o.Search = search
	o.SearchFields = fields
	return o
}

// WithPagination sets offset and limit
func (o *QueryOptions) WithPagination(offset, limit int) *QueryOptions {
	o.Offset = offset
	o.Limit = limit
	return o
}

// WithOrder sets the order clause
func (o *QueryOptions) WithOrder(order string) *QueryOptions {
	o.Order = order
	return o
}

// WithPreloads sets the relationships to preload
func (o *QueryOptions) WithPreloads(preloads ...string) *QueryOptions {
	o.Preloads = preloads
	return o
}

// BaseRepository provides common database operations for all repositories
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

// Create creates a new record
func (r *BaseRepository) Create(ctx context.Context, entity any) error {
	if err := r.DB.WithContext(ctx).Create(entity).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// CreateTx creates a new record within a transaction
func (r *BaseRepository) CreateTx(ctx context.Context, tx *gorm.DB, entity any) error {
	if err := tx.WithContext(ctx).Create(entity).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// Update updates a record
func (r *BaseRepository) Update(ctx context.Context, entity any) error {
	if err := r.DB.WithContext(ctx).Save(entity).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// UpdateTx updates a record within a transaction
func (r *BaseRepository) UpdateTx(ctx context.Context, tx *gorm.DB, entity any) error {
	if err := tx.WithContext(ctx).Save(entity).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// Delete deletes a record (soft delete if using DeletedAt)
func (r *BaseRepository) Delete(ctx context.Context, entity any) error {
	if err := r.DB.WithContext(ctx).Delete(entity).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// HardDelete permanently deletes a record
func (r *BaseRepository) HardDelete(ctx context.Context, entity any) error {
	if err := r.DB.WithContext(ctx).Unscoped().Delete(entity).Error; err != nil {
		return apperrors.ErrDatabase
	}
	return nil
}

// FindByID finds a record by ID
func (r *BaseRepository) FindByID(ctx context.Context, entity any, id any) error {
	if err := r.DB.WithContext(ctx).First(entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return nil
}

// FindByIDWithPreload finds a record by ID with preloaded relationships
func (r *BaseRepository) FindByIDWithPreload(ctx context.Context, entity any, id any, preloads ...string) error {
	db := r.DB.WithContext(ctx)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	if err := db.First(entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return nil
}

// FindOne finds a single record matching the given conditions
func (r *BaseRepository) FindOne(ctx context.Context, entity any, condition any, args ...any) error {
	if err := r.DB.WithContext(ctx).Where(condition, args...).First(entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return nil
}

// applyQueryOptions applies QueryOptions to a GORM query
func (r *BaseRepository) applyQueryOptions(query *gorm.DB, opts *QueryOptions) *gorm.DB {
	// Apply Where conditions
	if len(opts.Where) > 0 {
		query = query.Where(opts.Where)
	}

	// Apply Search (ILIKE)
	if opts.Search != "" && len(opts.SearchFields) > 0 {
		searchPattern := "%" + opts.Search + "%"
		conditions := make([]string, len(opts.SearchFields))
		args := make([]any, len(opts.SearchFields))
		for i, field := range opts.SearchFields {
			conditions[i] = fmt.Sprintf("%s ILIKE ?", field)
			args[i] = searchPattern
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	// Apply Preloads
	for _, preload := range opts.Preloads {
		query = query.Preload(preload)
	}

	// Apply pagination
	if opts.Offset > 0 {
		query = query.Offset(opts.Offset)
	}
	if opts.Limit > 0 {
		query = query.Limit(opts.Limit)
	}

	// Apply Order
	if opts.Order != "" {
		query = query.Order(opts.Order)
	}

	return query
}

// FindWithOptions finds records using QueryOptions
func (r *BaseRepository) FindWithOptions(ctx context.Context, model any, results any, opts *QueryOptions) error {
	query := r.DB.WithContext(ctx).Model(model)
	query = r.applyQueryOptions(query, opts)

	if err := query.Find(results).Error; err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return nil
}

// FindMany finds multiple records with QueryOptions
func (r *BaseRepository) FindMany(ctx context.Context, model any, results any, opts *QueryOptions) error {
	return r.FindWithOptions(ctx, model, results, opts)
}

// CountWithOptions counts records using QueryOptions
func (r *BaseRepository) CountWithOptions(ctx context.Context, model any, opts *QueryOptions) (int64, error) {
	var count int64
	query := r.DB.WithContext(ctx).Model(model)

	// Apply Where conditions
	if len(opts.Where) > 0 {
		query = query.Where(opts.Where)
	}

	// Apply Search (ILIKE)
	if opts.Search != "" && len(opts.SearchFields) > 0 {
		searchPattern := "%" + opts.Search + "%"
		conditions := make([]string, len(opts.SearchFields))
		args := make([]any, len(opts.SearchFields))
		for i, field := range opts.SearchFields {
			conditions[i] = fmt.Sprintf("%s ILIKE ?", field)
			args[i] = searchPattern
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return count, nil
}

// Count counts records using QueryOptions
func (r *BaseRepository) Count(ctx context.Context, model any, opts *QueryOptions) (int64, error) {
	var count int64
	query := r.DB.WithContext(ctx).Model(model)

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return count, nil
}

// Exists checks if a record exists
func (r *BaseRepository) Exists(ctx context.Context, model any, condition any, args ...any) (bool, error) {
	var count int64
	if err := r.DB.WithContext(ctx).Model(model).Where(condition, args...).Count(&count).Error; err != nil {
		return false, fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return count > 0, nil
}

// WithTx starts a new transaction
func (r *BaseRepository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.DB.Transaction(fn)
}

// GetDB returns the underlying DB instance
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.DB
}

// Raw executes a raw SQL query
func (r *BaseRepository) Raw(ctx context.Context, result any, sql string, args ...any) error {
	if err := r.DB.WithContext(ctx).Raw(sql, args...).Scan(result).Error; err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return nil
}

// Exec executes a raw SQL statement
func (r *BaseRepository) Exec(ctx context.Context, sql string, args ...any) error {
	if err := r.DB.WithContext(ctx).Exec(sql, args...).Error; err != nil {
		return fmt.Errorf("%w: %v", apperrors.ErrDatabase, err)
	}
	return nil
}