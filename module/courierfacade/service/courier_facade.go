package service

import (
	"context"
	cservice "geotask/module/courier/service"
	cfm "geotask/module/courierfacade/models"
	oservice "geotask/module/order/service"
	"log"
)

const (
	CourierVisibilityRadius = 2800 // 2.8km
)

type CourierFacer interface {
	MoveCourier(ctx context.Context, direction, zoom int) // отвечает за движение курьера по карте direction - направление движения, zoom - уровень зума
	GetStatus(ctx context.Context) cfm.CourierStatus      // отвечает за получение статуса курьера и заказов вокруг него
}

// CourierFacade фасад для курьера и заказов вокруг него (для фронта)
type CourierFacade struct {
	courierService cservice.Courierer
	orderService   oservice.Orderer
}

func NewCourierFacade(courierService cservice.Courierer, orderService oservice.Orderer) CourierFacer {
	return &CourierFacade{courierService: courierService, orderService: orderService}
}

func (cf *CourierFacade) MoveCourier(ctx context.Context, direction, zoom int) {
	// Реализация метода MoveCourier
	// Получение статуса курьера
	courierStatus, err := cf.courierService.GetCourier(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	cf.courierService.MoveCourier(*courierStatus, direction, zoom)
}

func (cf *CourierFacade) GetStatus(ctx context.Context) cfm.CourierStatus {
	// Реализация метода GetStatus
	//courierStatus, err := cf.courierService.GetCourier(ctx)
	//if err != nil {
	//	log.Printf("Error getting courier status: %v", err)
	//	return cfm.CourierStatus{}
	//}
	return cfm.CourierStatus{}
	//return cfm.CourierStatus{} // Пример, нужно заменить на реальную логику
}
