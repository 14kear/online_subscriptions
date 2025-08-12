package repository

import (
	"context"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/entity"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveRecord(ctx context.Context, record *entity.Record) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *Repository) DeleteRecordByID(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.Record{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) GetRecordByID(ctx context.Context, id uint) (*entity.Record, error) {
	var record entity.Record

	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *Repository) GetRecordsByUserID(ctx context.Context, userID string) ([]entity.Record, error) {
	var records []entity.Record

	query := r.db.WithContext(ctx).Model(&entity.Record{}).Where("user_id = ?", userID)
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *Repository) GetRecordByUserIDAndServiceName(ctx context.Context, userID, serviceName string) (*entity.Record, error) {
	var record entity.Record

	if err := r.db.WithContext(ctx).Where("user_id = ? AND service_name = ?", userID, serviceName).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *Repository) UpdateRecord(ctx context.Context, record *entity.Record) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *Repository) ListRecords(ctx context.Context, limit, offset int, userID, serviceName string) ([]entity.Record, error) {
	var records []entity.Record

	query := r.db.WithContext(ctx).Model(&entity.Record{})

	// фильтрация
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	// пагинация
	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset >= 0 {
		query = query.Offset(offset)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *Repository) SumPriceForPeriod(ctx context.Context, startTime, endTime time.Time, userID, serviceName string) (int, error) {
	var total int

	query := r.db.WithContext(ctx).Model(&entity.Record{}).
		Where("created_at BETWEEN ? AND ?", startTime, endTime)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if err := query.Select("SUM(price)").Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
