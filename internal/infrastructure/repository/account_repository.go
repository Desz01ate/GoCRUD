package repository

import (
	"arise_tech_assetment/internal/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Repository[domain.Account, uuid.UUID]
	FindByNumber(ctx context.Context, number string) (*domain.Account, error)
	FindByStatus(ctx context.Context, status domain.AccountStatus) ([]domain.Account, error)
	FindByStatusPaginated(ctx context.Context, status domain.AccountStatus, req PaginationRequest) (*PaginationResponse[domain.Account], error)
	FindByHolderName(ctx context.Context, holderName string) ([]domain.Account, error)
	FindByHolderNamePaginated(ctx context.Context, holderName string, req PaginationRequest) (*PaginationResponse[domain.Account], error)
}

type accountRepository struct {
	*GormRepository[domain.Account, uuid.UUID]
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		GormRepository: NewGormRepository[domain.Account, uuid.UUID](db),
	}
}

func (r *accountRepository) FindByNumber(ctx context.Context, number string) (*domain.Account, error) {
	var account domain.Account
	if err := r.GormRepository.db.WithContext(ctx).Where("number = ?", number).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByStatus(ctx context.Context, status domain.AccountStatus) ([]domain.Account, error) {
	var accounts []domain.Account
	if err := r.GormRepository.db.WithContext(ctx).Where("status = ?", status).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) FindByHolderName(ctx context.Context, holderName string) ([]domain.Account, error) {
	var accounts []domain.Account
	if err := r.GormRepository.db.WithContext(ctx).Where("holder_name ILIKE ?", "%"+holderName+"%").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) FindByStatusPaginated(ctx context.Context, status domain.AccountStatus, req PaginationRequest) (*PaginationResponse[domain.Account], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var accounts []domain.Account
	var total int64

	query := r.GormRepository.db.WithContext(ctx).Where("status = ?", status)

	if err := query.Model(&domain.Account{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&accounts).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &PaginationResponse[domain.Account]{
		Data:       accounts,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *accountRepository) FindByHolderNamePaginated(ctx context.Context, holderName string, req PaginationRequest) (*PaginationResponse[domain.Account], error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var accounts []domain.Account
	var total int64

	query := r.GormRepository.db.WithContext(ctx).Where("holder_name ILIKE ?", "%"+holderName+"%")

	if err := query.Model(&domain.Account{}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&accounts).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &PaginationResponse[domain.Account]{
		Data:       accounts,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}
