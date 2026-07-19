package common

import (
	"context"

	"gorm.io/gorm"
)

// BaseRepository defines standard CRUD operations that all repositories should implement.
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id string) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, pagination *PaginationQuery) ([]T, int64, error)
}

// GormRepository provides a generic implementation of BaseRepository using GORM.
type GormRepository[T any] struct {
	db *gorm.DB
}

// NewGormRepository creates a new generic GORM repository.
func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

// Create inserts a new record.
func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByID retrieves a record by its primary key.
func (r *GormRepository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Or return a custom not found error
		}
		return nil, err
	}
	return &entity, nil
}

// Update saves changes to an existing record.
func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete marks a record as deleted (soft delete if struct has gorm.DeletedAt).
func (r *GormRepository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
}

// List returns a paginated list of records and the total count.
func (r *GormRepository[T]) List(ctx context.Context, pagination *PaginationQuery) ([]T, int64, error) {
	var items []T
	var total int64

	query := r.db.WithContext(ctx).Model(new(T))
	
	// Count total records (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	limit, offset := pagination.GetLimitOffset()
	order := pagination.GetOrder()
	
	err := query.Order(order).Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
