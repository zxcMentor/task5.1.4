package service

import (
	"context"
	"geotask/geo"
	"geotask/module/courier/models"
	"geotask/module/courier/storage"
	"math"
	"math/rand"
	"time"
)

// Направления движения курьера
const (
	DirectionUp    = 0
	DirectionDown  = 1
	DirectionLeft  = 2
	DirectionRight = 3
)

const (
	DefaultCourierLat = 59.9311
	DefaultCourierLng = 30.3609
)
const (
	allowedZoneMaxLatitude  = 90.0
	allowedZoneMinLatitude  = -90.0
	allowedZoneMaxLongitude = 180.0
	allowedZoneMinLongitude = -180.0
)

type Courierer interface {
	GetCourier(ctx context.Context) (*models.Courier, error)
	MoveCourier(courier models.Courier, direction, zoom int) error
}

type CourierService struct {
	courierStorage storage.CourierStorager
	allowedZone    geo.PolygonChecker
	disabledZones  []geo.PolygonChecker
}

func NewCourierService(courierStorage storage.CourierStorager, allowedZone geo.PolygonChecker, disbledZones []geo.PolygonChecker) Courierer {
	return &CourierService{courierStorage: courierStorage, allowedZone: allowedZone, disabledZones: disbledZones}
}

func (c *CourierService) GetCourier(ctx context.Context) (*models.Courier, error) {
	// получаем курьера из хранилища используя метод GetOne из storage/courier.go
	// проверяем, что курьер находится в разрешенной зоне
	// если нет, то перемещаем его в случайную точку в разрешенной зоне
	// сохраняем новые координаты курьера
	courierStorage := &storage.CourierStorage{} // Пример: создание экземпляра хранилища
	courier, err := courierStorage.GetOne(ctx)
	if err != nil {
		return nil, err
	}

	// Проверяем, что курьер находится в разрешенной зоне
	if !isInAllowedZone(courier) {
		// Если курьер не в разрешенной зоне, перемещаем его в случайную точку в разрешенной зоне
		moveCourierToRandomPointInAllowedZone(courier)

		// Сохраняем новые координаты курьера
		err := courierStorage.Save(ctx, *courier)
		if err != nil {
			return nil, err
		}
	}

	return courier, nil
}

// Функция для проверки, находится ли курьер в разрешенной зоне
func isInAllowedZone(courier *models.Courier) bool {
	// Реализуйте логику проверки зоны в соответствии с вашими требованиями
	// Верните true, если курьер в разрешенной зоне, и false в противном случае
	// Предположим, что курьер хранит свои координаты в полях Latitude и Longitude
	latitude := courier.Location.Lat
	longitude := courier.Location.Lng

	// Проверяем, находится ли курьер в разрешенной зоне
	if latitude >= allowedZoneMinLatitude &&
		latitude <= allowedZoneMaxLatitude &&
		longitude >= allowedZoneMinLongitude &&
		longitude <= allowedZoneMaxLongitude {
		return true
	}
	return false
}

// Функция для перемещения курьера в случайную точку в разрешенной зоне
func moveCourierToRandomPointInAllowedZone(courier *models.Courier) {
	// Реализуйте логику перемещения курьера в случайную точку в разрешенной зоне
	// Обновите координаты курьера со случайными значениями в разрешенной зоне
	rand.Seed(time.Now().UnixNano())
	courier.Location.Lat = rand.Float64()*(allowedZoneMaxLatitude-allowedZoneMinLatitude) + allowedZoneMaxLatitude
	courier.Location.Lng = rand.Float64()*(allowedZoneMaxLongitude-allowedZoneMinLongitude) + allowedZoneMinLongitude
}

// MoveCourier : direction - направление движения курьера, zoom - зум карты
func (c *CourierService) MoveCourier(courier models.Courier, direction, zoom int) error {
	// точность перемещения зависит от зума карты использовать формулу 0.001 / 2^(zoom - 14)
	// 14 - это максимальный зум карты

	// далее нужно проверить, что курьер не вышел за границы зоны
	// если вышел, то нужно переместить его в случайную точку внутри зоны

	// далее сохранить изменения в хранилище
	// Точность перемещения зависит от зума карты, используем формулу 0.001 / 2^(zoom - 14)
	precision := 0.001 / math.Pow(2, float64(zoom-14))

	// Рассчитываем новые координаты курьера в зависимости от направления и точности
	newLatitude, newLongitude := calculateNewCoordinates(courier.Location.Lat, courier.Location.Lng, direction, precision)

	// Проверяем, что курьер не вышел за границы зеленой зоны
	if !isInGreenZone(newLatitude, newLongitude) {
		// Если вышел, перемещаем его в случайную точку внутри зеленой зоны
		moveCourierToRandomPointInGreenZone(&courier)
	} else {
		// Иначе обновляем координаты курьера
		courier.Location.Lat = newLatitude
		courier.Location.Lng = newLongitude
	}

	// Сохраняем изменения в хранилище
	err := c.courierStorage
	if err != nil {
		return nil
	}

	return nil
}

// Функция для рассчета новых координат на основе направления и точности
func calculateNewCoordinates(latitude, longitude float64, direction int, precision float64) (float64, float64) {
	// Реализуйте логику рассчета новых координат на основе направления и точности
	// Зависит от вашей конкретной логики, например, можно использовать тригонометрию
	// Возвращаем новые координаты
	// Преобразование направления в радианы
	radDirection := float64(direction) * (math.Pi / 180.0)

	// Рассчет новых координат с использованием тригонометрии
	newLatitude := latitude + precision*math.Cos(radDirection)
	newLongitude := longitude + precision*math.Sin(radDirection)

	return newLatitude, newLongitude
}

// Функция для проверки, находится ли курьер в зеленой зоне
func isInGreenZone(latitude, longitude float64) bool {
	// Реализуйте логику проверки, находится ли курьер в зеленой зоне
	// Возвращаем true, если курьер в зеленой зоне, и false в противном случае

	// Проверяем, находится ли курьер в разрешенной зоне
	if latitude >= allowedZoneMinLatitude &&
		latitude <= allowedZoneMaxLatitude &&
		longitude >= allowedZoneMinLongitude &&
		longitude <= allowedZoneMaxLongitude {
		return true
	}
	return false

}

// Функция для перемещения курьера в случайную точку внутри зеленой зоны
func moveCourierToRandomPointInGreenZone(courier *models.Courier) {
	// Реализуйте логику перемещения курьера в случайную точку внутри зеленой зоны
	// Обновите координаты курьера со случайными значениями в зеленой зоне
	rand.Seed(time.Now().UnixNano())
	courier.Location.Lat = rand.Float64()*(allowedZoneMaxLatitude-allowedZoneMinLatitude) + allowedZoneMaxLatitude
	courier.Location.Lng = rand.Float64()*(allowedZoneMaxLongitude-allowedZoneMinLongitude) + allowedZoneMinLongitude
}
