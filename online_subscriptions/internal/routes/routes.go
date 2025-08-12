package routes

import (
	"github.com/14kear/effective_mobile/online_subscriptions/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *handlers.RecordHandler) {
	router.POST("/create", handler.CreateRecord)
	router.DELETE("/delete/:id", handler.DeleteRecord)
	router.PUT("/update/:id", handler.UpdateRecord)

	// получить по id
	router.GET("/record/:id", handler.GetRecordByID)

	// получить по user_id
	router.GET("/records/user", handler.GetRecordsByUserID)

	// получить по user_id + service_name
	router.GET("/record/user_service", handler.GetRecordByUserIDAndServiceName)

	// получить список с фильтрацией
	router.GET("/records", handler.ListRecords)

	// сумма за период
	router.GET("/records/summary", handler.SumPriceForPeriod)
}
