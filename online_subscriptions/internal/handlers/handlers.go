package handlers

import (
	"errors"
	_ "github.com/14kear/effective_mobile/online_subscriptions/docs"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/entity"
	"github.com/14kear/effective_mobile/online_subscriptions/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// RecordCreateUpdateRequest для создания и обновления записи
type RecordCreateUpdateRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	ExpiresAt   string `json:"expires_at" binding:"required,datetime=02-01-2006"`
	CreatedAt   string `json:"created_at" binding:"omitempty,datetime=02-01-2006"`
}

// RecordQuery для поиска записи по пользователю и сервису
type RecordQuery struct {
	ServiceName string `json:"service_name" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
}

// SumPeriodQuery для суммирования платежей за период
type SumPeriodQuery struct {
	StartTime   string `json:"start_time" binding:"required,datetime=02-01-2006"`
	EndTime     string `json:"end_time" binding:"required,datetime=02-01-2006"`
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
}

// RecordHandler обрабатывает запросы для записей подписок
type RecordHandler struct {
	RecordService *services.RecordService
}

// NewRecordHandler создает новый экземпляр RecordHandler
func NewRecordHandler(recordService *services.RecordService) *RecordHandler {
	return &RecordHandler{RecordService: recordService}
}

// CreateRecord создает новую запись подписки
// @Summary Создать запись подписки
// @Description Создает новую запись онлайн подписки
// @Tags Подписки
// @Accept json
// @Produce json
// @Param input body RecordCreateUpdateRequest true "Данные подписки"
// @Success 201 {object} entity.Record "Созданная запись"
// @Failure 400 {object} map[string]string "Неверный формат данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /create [post]
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

// DeleteRecord удаляет запись подписки
// @Summary Удалить запись подписки
// @Description Удаляет запись подписки по указанному ID
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path int true "ID записи" example(1)
// @Success 204 "Запись успешно удалена"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Запись не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /delete/{id} [delete]
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

// UpdateRecord обновляет запись подписки
// @Summary Обновить запись подписки
// @Description Обновляет существующую запись подписки
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path int true "ID записи" example(1)
// @Param input body RecordCreateUpdateRequest true "Данные для обновления"
// @Success 204 "Запись успешно обновлена"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Запись не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /update/{id} [put]
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

// GetRecordByID получает запись подписки по ID
// @Summary Получить запись подписки
// @Description Возвращает запись подписки по указанному ID
// @Tags Подписки
// @Accept json
// @Produce json
// @Param id path int true "ID записи" example(1)
// @Success 200 {object} entity.Record "Запись подписки"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Запись не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /record/{id} [get]
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

// GetRecordsByUserID получает записи подписок пользователя
// @Summary Получить подписки пользователя
// @Description Возвращает список подписок для указанного пользователя
// @Tags Подписки
// @Accept json
// @Produce json
// @Param user_id query string true "ID пользователя" example(user123)
// @Success 200 {array} entity.Record "Список подписок"
// @Failure 400 {object} map[string]string "Неверный ID пользователя"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /records/user [get]
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

// GetRecordByUserIDAndServiceName получает запись подписки
// @Summary Найти подписку пользователя
// @Description Возвращает конкретную подписку пользователя по названию сервиса
// @Tags Подписки
// @Accept json
// @Produce json
// @Param user_id query string true "ID пользователя" example(user123)
// @Param service_name query string true "Название сервиса" example(Netflix)
// @Success 200 {object} entity.Record "Запись подписки"
// @Failure 400 {object} map[string]string "Неверные параметры"
// @Failure 404 {object} map[string]string "Запись не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /record/user_service [get]
func (h *RecordHandler) GetRecordByUserIDAndServiceName(ctx *gin.Context) {
	var req struct {
		ServiceName string `form:"service_name" binding:"required"`
		UserID      string `form:"user_id" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
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

// ListRecords получает список записей
// @Summary Список подписок
// @Description Возвращает список подписок с фильтрацией и пагинацией
// @Tags Подписки
// @Accept json
// @Produce json
// @Param user_id query string false "Фильтр по ID пользователя" example(user123)
// @Param service_name query string false "Фильтр по названию сервиса" example(Netflix)
// @Param limit query int false "Лимит записей (макс. 100)" minimum(1) maximum(100) default(20)
// @Param offset query int false "Смещение" minimum(0) default(0)
// @Success 200 {array} entity.Record "Список подписок"
// @Failure 400 {object} map[string]string "Неверные параметры"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /records [get]
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

// SumPriceForPeriod вычисляет сумму платежей
// @Summary Сумма платежей за период
// @Description Возвращает сумму платежей за указанный период с возможностью фильтрации
// @Tags Аналитика
// @Accept json
// @Produce json
// @Param start_time query string true "Начальная дата (DD-MM-YYYY)" example(01-01-2023)
// @Param end_time query string true "Конечная дата (DD-MM-YYYY)" example(31-12-2023)
// @Param user_id query string false "Фильтр по ID пользователя" example(user123)
// @Param service_name query string false "Фильтр по названию сервиса" example(Netflix)
// @Success 200 {object} map[string]int "{"total_price": 1500}"
// @Failure 400 {object} map[string]string "Неверные параметры"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /records/summary [get]
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
