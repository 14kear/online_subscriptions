package handlers

import (
	"errors"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/entity"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type RecordHandler struct {
	RecordService *services.RecordService
}

func NewRecordHandler(recordService *services.RecordService) *RecordHandler {
	return &RecordHandler{RecordService: recordService}
}

func (h *RecordHandler) CreateRecord(ctx *gin.Context) {
	var createdAt time.Time

	var req struct {
		ServiceName string `json:"service_name" binding:"required"`
		Price       int    `json:"price" binding:"required"`
		UserID      string `json:"user_id" binding:"required"`
		ExpiresAt   string `json:"expires_at" binding:"required,datetime=02-01-2006"`
		CreatedAt   string `json:"created_at" binding:"omitempty,datetime=02-01-2006"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expiresAt, err := time.Parse("02-01-2006", req.ExpiresAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date must be DD-MM-YYYY"})
		return
	}

	if req.CreatedAt != "" {
		createdAt, err = time.Parse("02-01-2006", req.CreatedAt)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date must be DD-MM-YYYY"})
			return
		}
	}

	record := entity.Record{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		ExpiresAt:   expiresAt,
		CreatedAt:   createdAt,
	}

	err = h.RecordService.CreateRecord(ctx.Request.Context(), &record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, record)
}

func (h *RecordHandler) DeleteRecord(ctx *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.RecordService.DeleteRecordByID(ctx.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *RecordHandler) UpdateRecord(ctx *gin.Context) {
	var createdAt time.Time

	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if uri.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be > 0"})
		return
	}

	var req struct {
		ServiceName string `json:"service_name" binding:"required"`
		Price       int    `json:"price" binding:"required"`
		UserID      string `json:"user_id" binding:"required"`
		ExpiresAt   string `json:"expires_at" binding:"required,datetime=02-01-2006"`
		CreatedAt   string `json:"created_at" binding:"omitempty,datetime=02-01-2006"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expiresAt, err := time.Parse("02-01-2006", req.ExpiresAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date must be DD-MM-YYYY"})
		return
	}

	if req.CreatedAt != "" {
		createdAt, err = time.Parse("02-01-2006", req.CreatedAt)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date must be DD-MM-YYYY"})
			return
		}
	}

	record := entity.Record{
		ID:          uri.ID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		ExpiresAt:   expiresAt,
		CreatedAt:   createdAt,
	}

	err = h.RecordService.UpdateRecord(ctx.Request.Context(), &record)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *RecordHandler) GetRecordByID(ctx *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if uri.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be > 0"})
		return
	}

	record, err := h.RecordService.GetRecordByID(ctx.Request.Context(), uri.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, record)
}

func (h *RecordHandler) GetRecordsByUserID(ctx *gin.Context) {
	// в ТЗ указано, что управление пользователями вне зоны ответственности моего сервиса,
	// поэтому в данный момент общаемся json, uri, query, в полной версии я бы доставал uuid из контекста
	var req struct {
		UserID string `form:"user_id" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	records, err := h.RecordService.GetRecordsByUserID(ctx.Request.Context(), req.UserID)
	if err != nil {
		// 404 возвращать не буду, логичнее здесь просто вернуть пустой список

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

func (h *RecordHandler) GetRecordByUserIDAndServiceName(ctx *gin.Context) {
	var req struct {
		ServiceName string `json:"service_name" binding:"required"`
		UserID      string `json:"user_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := h.RecordService.GetRecordByUserIDAndServiceName(ctx.Request.Context(), req.UserID, req.ServiceName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, record)
}

func (h *RecordHandler) ListRecords(ctx *gin.Context) {
	var req struct {
		UserID      string `form:"user_id"`
		ServiceName string `form:"service_name"`
		Limit       int    `form:"limit"`
		Offset      int    `form:"offset"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// дефолт для лимита
	if req.Limit == 0 {
		req.Limit = 20
	}
	// защита от слишком большого лимита
	if req.Limit > 100 {
		req.Limit = 100
	}

	records, err := h.RecordService.ListRecords(ctx.Request.Context(), req.Limit, req.Offset, req.UserID, req.ServiceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, records)
}

func (h *RecordHandler) SumPriceForPeriod(ctx *gin.Context) {
	var req struct {
		StartTime   string `json:"start_time" binding:"required,datetime=02-01-2006"`
		EndTime     string `json:"end_time" binding:"required,datetime=02-01-2006"`
		UserID      string `json:"user_id"`
		ServiceName string `json:"service_name"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime, err := time.Parse("02-01-2006", req.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date must be DD-MM-YYYY"})
		return
	}

	endTime, err := time.Parse("02-01-2006", req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date must be DD-MM-YYYY"})
		return
	}

	if endTime.Before(startTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "end time must be greater than start time"})
		return
	}

	total, err := h.RecordService.SummaryPriceOfSelectedRecords(
		ctx.Request.Context(),
		startTime,
		endTime,
		req.UserID,
		req.ServiceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, total)
}
