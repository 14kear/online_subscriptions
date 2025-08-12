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
	var req struct {
		ServiceName string    `json:"service_name" binding:"required"`
		Price       int       `json:"price" binding:"required"`
		UserID      string    `json:"user_id" binding:"required"`
		ExpiresAt   time.Time `json:"expires_at" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record := entity.Record{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		ExpiresAt:   req.ExpiresAt,
	}

	err := h.RecordService.CreateRecord(ctx.Request.Context(), &record)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"record": record})
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
		ServiceName string    `json:"service_name" binding:"required"`
		Price       int       `json:"price" binding:"required"`
		UserID      string    `json:"user_id" binding:"required"`
		ExpiresAt   time.Time `json:"expires_at" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record := entity.Record{
		ID:          uri.ID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		ExpiresAt:   req.ExpiresAt,
	}

	err := h.RecordService.UpdateRecord(ctx.Request.Context(), &record)
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

	ctx.JSON(http.StatusOK, gin.H{"record": record})
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

	ctx.JSON(http.StatusOK, gin.H{"records": records})
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

	ctx.JSON(http.StatusOK, gin.H{"record": record})
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

	ctx.JSON(http.StatusOK, gin.H{"records": records})
}

func (h *RecordHandler) SumPriceForPeriod(ctx *gin.Context) {
	var req struct {
		StartTime   time.Time `json:"start_time" binding:"required"`
		EndTime     time.Time `json:"end_time" binding:"required"`
		UserID      string    `json:"user_id" binding:"required"`
		ServiceName string    `json:"service_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EndTime.Before(req.StartTime) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "end time must be greater than start time"})
		return
	}

	total, err := h.RecordService.SummaryPriceOfSelectedRecords(
		ctx.Request.Context(),
		req.StartTime,
		req.EndTime,
		req.UserID,
		req.ServiceName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"total_price": total})
}
