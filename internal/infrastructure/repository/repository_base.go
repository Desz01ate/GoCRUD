package repository

import (
	"context"

	"gorm.io/gorm"
)

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PaginationResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type Repository[T any, TKey any] interface {
	GetByID(ctx context.Context, id TKey) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	GetPaginated(ctx context.Context, req PaginationRequest) (*PaginationResponse[T], error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id TKey) error
}

type GormRepository[T any, TKey any] struct {
	db *gorm.DB
}

func NewGormRepository[T any, TKey any](db *gorm.DB) *GormRepository[T, TKey] {
	return &GormRepository[T, TKey]{db: db}
}

func (r *GormRepository[T, TKey]) GetByID(ctx context.Context, id TKey) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T, TKey]) GetAll(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GormRepository[T, TKey]) GetPaginated(ctx context.Context, req PaginationRequest) (*PaginationResponse[T], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var entities []T
	var total int64

	var entity T
	if err := r.db.WithContext(ctx).Model(&entity).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(req.PageSize).Find(&entities).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &PaginationResponse[T]{
		Data:       entities,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *GormRepository[T, TKey]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository[T, TKey]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[T, TKey]) Delete(ctx context.Context, id TKey) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}
