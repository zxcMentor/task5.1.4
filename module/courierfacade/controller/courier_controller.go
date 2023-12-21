package controller

import (
	"context"
	"encoding/json"
	"geotask/module/courier/models"
	"geotask/module/courierfacade/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CourierController struct {
	courierService service.CourierFacer
}

func NewCourierController(courierService service.CourierFacer) *CourierController {
	return &CourierController{courierService: courierService}
}

func (c *CourierController) GetStatus(ctx *gin.Context) {
	// установить задержку в 50 миллисекунд

	// получить статус курьера из сервиса courierService используя метод GetStatus
	// отправить статус курьера в ответ
	// Установить задержку в 50 миллисекунд
	// Установить задержку в 50 миллисекунд
	time.Sleep(50 * time.Millisecond)

	// Получить статус курьера из сервиса courierService используя метод GetStatus
	status := c.courierService.GetStatus(ctx)

	// Отправить статус курьера в ответ
	ctx.JSON(http.StatusOK, gin.H{"status": status, "location": models.Point{}})
}

func (c *CourierController) MoveCourier(m webSocketMessage, dir, zoom int) {
	var cm CourierMove
	var err error
	// Получить данные из m.Data и десериализовать их в структуру CourierMove
	data, ok := m.Data.(string)
	if !ok {
		// Обработка ошибки: невозможно преобразовать данные в строку
		return
	}

	err = json.Unmarshal([]byte(data), &cm)
	if err != nil {
		// Обработка ошибки десериализации данных
		return
	}
	ctx := context.Background()
	// Вызвать метод MoveCourier у courierService
	c.courierService.MoveCourier(ctx, dir, zoom)
}
