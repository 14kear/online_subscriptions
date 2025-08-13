package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/entity"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

var (
	ErrGetFailed    = errors.New("could not get record")
	ErrCreateFailed = errors.New("could not create record")
	ErrUpdateFailed = errors.New("could not update record")
	ErrDeleteFailed = errors.New("could not delete record")
	ErrSumFailed    = errors.New("could not sum records")
)

type Repository interface {
	SaveRecord(ctx context.Context, record *entity.Record) error
	DeleteRecordByID(ctx context.Context, id uint) error
	GetRecordByID(ctx context.Context, id uint) (*entity.Record, error)
	GetRecordsByUserID(ctx context.Context, userID string) ([]entity.Record, error)
	GetRecordByUserIDAndServiceName(ctx context.Context, userID string, serviceName string) (*entity.Record, error)
	UpdateRecord(ctx context.Context, record *entity.Record) error
	ListRecords(ctx context.Context, limit, offset int, userID, serviceName string) ([]entity.Record, error)
	SumPriceForPeriod(ctx context.Context, startTime, endTime time.Time, userID, serviceName string) (int, error)
}

type RecordService struct {
	log              *slog.Logger
	recordRepository Repository
}

func NewRecordService(log *slog.Logger, recordRepository Repository) *RecordService {
	return &RecordService{log: log, recordRepository: recordRepository}
}

func (s *RecordService) CreateRecord(ctx context.Context, record *entity.Record) error {
	const op = "recordService.CreateRecord"

	log := s.log.With(slog.String("operation", op))
	log.Info("creating new record...")

	now := time.Now()
	created := record.CreatedAt
	if created.IsZero() {
		created = now
	}

	if record.ExpiresAt.Before(record.CreatedAt) || record.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("%w: expires date must be after created date", ErrCreateFailed)
	}

	if err := s.recordRepository.SaveRecord(ctx, record); err != nil {
		log.Error("failed to save record", slog.Any("error", err))
		return fmt.Errorf("%w: %v", ErrCreateFailed, err)
	}

	log.Info("record successfully created")

	return nil
}

func (s *RecordService) DeleteRecordByID(ctx context.Context, id uint) error {
	const op = "recordService.DeleteRecordByID"

	log := s.log.With(slog.String("operation", op))

	log.Info("deleting record...")

	if err := s.recordRepository.DeleteRecordByID(ctx, id); err != nil {
		log.Error("failed to delete record", slog.Any("error", err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}

	log.Info("record successfully deleted")

	return nil
}

func (s *RecordService) UpdateRecord(ctx context.Context, record *entity.Record) error {
	const op = "recordService.UpdateRecord"

	log := s.log.With(slog.String("operation", op))
	log.Info("updating record...")

	now := time.Now()
	created := record.CreatedAt
	if created.IsZero() {
		created = now
	}

	if record.ExpiresAt.Before(record.CreatedAt) || record.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("%w: expires date must be after created date", ErrCreateFailed)
	}

	if err := s.recordRepository.UpdateRecord(ctx, record); err != nil {
		log.Error("failed to update record", slog.Any("error", err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
	}

	log.Info("record successfully updated")

	return nil
}

func (s *RecordService) GetRecordByID(ctx context.Context, id uint) (*entity.Record, error) {
	const op = "recordService.GetRecordByID"

	log := s.log.With(slog.String("operation", op))

	log.Info("getting record...")

	record, err := s.recordRepository.GetRecordByID(ctx, id)
	if err != nil {
		log.Error("failed to get record", slog.Any("error", err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("%w: %v", ErrGetFailed, err)
	}

	log.Info("record successfully retrieved")

	return record, nil
}

func (s *RecordService) GetRecordsByUserID(ctx context.Context, userID string) ([]entity.Record, error) {
	const op = "recordService.GetRecordsByUserID"

	log := s.log.With(slog.String("operation", op))
	log.Info("getting records...")

	records, err := s.recordRepository.GetRecordsByUserID(ctx, userID)
	if err != nil {
		log.Error("failed to get records", slog.Any("error", err))
		return nil, fmt.Errorf("%w: %v", ErrGetFailed, err)
	}

	log.Info("records successfully retrieved")

	return records, nil
}

func (s *RecordService) GetRecordByUserIDAndServiceName(ctx context.Context, userID, serviceName string) (*entity.Record, error) {
	const op = "recordService.GetRecordByUserIdAndServiceName"

	log := s.log.With(slog.String("operation", op))
	log.Info("getting record...")

	record, err := s.recordRepository.GetRecordByUserIDAndServiceName(ctx, userID, serviceName)
	if err != nil {
		log.Error("failed to get record", slog.Any("error", err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("%w: %v", ErrGetFailed, err)
	}

	log.Info("record successfully retrieved")

	return record, nil
}

func (s *RecordService) ListRecords(ctx context.Context, limit, offset int, userID, serviceName string) ([]entity.Record, error) {
	const op = "recordService.ListRecords"

	log := s.log.With(slog.String("operation", op))
	log.Info("getting records...")

	records, err := s.recordRepository.ListRecords(ctx, limit, offset, userID, serviceName)
	if err != nil {
		log.Error("failed to get records", slog.Any("error", err))
		return nil, fmt.Errorf("%w: %v", ErrGetFailed, err)
	}

	log.Info("records successfully retrieved")

	return records, nil
}

func (s *RecordService) SummaryPriceOfSelectedRecords(ctx context.Context, startTime, endTime time.Time, userID, serviceName string) (int, error) {
	const op = "recordService.SummaryPriceOfSelectedRecords"

	log := s.log.With(slog.String("operation", op))
	log.Info("summary records...")

	total, err := s.recordRepository.SumPriceForPeriod(ctx, startTime, endTime, userID, serviceName)
	if err != nil {
		log.Error("failed to get records", slog.Any("error", err))
		return 0, fmt.Errorf("%w: %v", ErrSumFailed, err)
	}

	log.Info("records successfully summary")

	return total, nil
}
