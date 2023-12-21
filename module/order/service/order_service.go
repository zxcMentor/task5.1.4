package service

import (
	"context"
	"geotask/geo"
	"geotask/module/order/models"
	"geotask/module/order/storage"
	"time"
)

const (
	minDeliveryPrice = 100.00
	maxDeliveryPrice = 500.00

	maxOrderPrice = 3000.00
	minOrderPrice = 1000.00

	orderMaxAge = 2 * time.Minute
)

type Orderer interface {
	GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) // возвращает заказы через метод storage.GetByRadius
	Save(ctx context.Context, order models.Order) error                                             // сохраняет заказ через метод storage.Save с заданным временем жизни OrderMaxAge
	GetCount(ctx context.Context) (int, error)                                                      // возвращает количество заказов через метод storage.GetCount
	RemoveOldOrders(ctx context.Context) error                                                      // удаляет старые заказы через метод storage.RemoveOldOrders с заданным временем жизни OrderMaxAge
	GenerateOrder(ctx context.Context) error                                                        // генерирует заказ в случайной точке из разрешенной зоны, с уникальным id, ценой и ценой доставки
}

// OrderService реализация интерфейса Orderer
// в нем должны быть методы GetByRadius, Save, GetCount, RemoveOldOrders, GenerateOrder
// данный сервис отвечает за работу с заказами
type OrderService struct {
	storage       storage.OrderStorager
	allowedZone   geo.PolygonChecker
	disabledZones []geo.PolygonChecker
}

func NewOrderService(storage storage.OrderStorager, allowedZone geo.PolygonChecker, disallowedZone []geo.PolygonChecker) Orderer {
	return &OrderService{storage: storage, allowedZone: allowedZone, disabledZones: disallowedZone}
}
func (os *OrderService) GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) {
	// Реализация метода GetByRadius
	return nil, nil // Замените на реальную логику
}

func (os *OrderService) Save(ctx context.Context, order models.Order) error {
	// Реализация метода Save
	return nil // Замените на реальную логику
}

func (os *OrderService) GetCount(ctx context.Context) (int, error) {
	// Реализация метода GetCount
	return 0, nil // Замените на реальную логику
}

func (os *OrderService) RemoveOldOrders(ctx context.Context) error {
	// Реализация метода RemoveOldOrders
	return nil // Замените на реальную логику
}

func (os *OrderService) GenerateOrder(ctx context.Context) error {
	// Реализация метода GenerateOrder
	return nil // Замените на реальную логику
}
