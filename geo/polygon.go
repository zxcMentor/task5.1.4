package geo

import (
	"fmt"
	geo "github.com/kellydunn/golang-geo"
	"math/rand"
)

// Point представляет собой географическую точку с широтой и долготой.
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// PolygonChecker - интерфейс для проверки положения точки и разрешения входа в полигон.
type PolygonChecker interface {
	Contains(point Point) bool // проверить, находится ли точка внутри полигона
	Allowed() bool             // разрешено ли входить в полигон
	RandomPoint() Point        // сгенерировать случайную точку внутри полигона
}

// Polygon представляет собой полигон и его разрешение.
type Polygon struct {
	polygon *geo.Polygon
	allowed bool
}

// NewPolygon создает новый полигон с использованием библиотеки golang-geo.
func NewPolygon(points []Point, allowed bool) (*Polygon, error) {
	// Создаем массив указателей на точки для библиотеки golang-geo.
	geoPoints := make([]*geo.Point, len(points))
	for i, p := range points {
		geoPoints[i] = geo.NewPoint(p.Lat, p.Lng)
	}

	polygon := geo.NewPolygon(geoPoints)
	if !isPolygonValid(polygon) {
		return nil, fmt.Errorf("invalid polygon")
	}

	return &Polygon{
		polygon: polygon,
		allowed: allowed,
	}, nil
}

// isPolygonValid проверяет, что у полигона не более одной положительной и отрицательной вершины.
func isPolygonValid(polygon *geo.Polygon) bool {
	posCount := 0
	negCount := 0

	for _, point := range polygon.Points() {
		if point.Lng() > 0 {
			posCount++
		} else if point.Lng() < 0 {
			negCount++
		}
	}

	return posCount <= 1 && negCount <= 1
}

// Contains проверяет, находится ли точка внутри полигона.
func (p *Polygon) Contains(point Point) bool {
	return p.polygon.Contains(geo.NewPoint(point.Lat, point.Lng))
}

// Allowed возвращает, разрешен ли вход в полигон.
func (p *Polygon) Allowed() bool {
	return p.allowed
}

// RandomPoint генерирует случайную точку внутри полигона.
func (p *Polygon) RandomPoint() Point {
	for {
		// Генерируем случайную точку в пределах общих границ карты.
		lat := rand.Float64()*(90.0-(-90.0)) + (-90.0)
		lng := rand.Float64()*(180.0-(-180.0)) + (-180.0)

		// Создаем новую точку для проверки.
		randomPoint := geo.NewPoint(lat, lng)

		// Проверяем, находится ли точка внутри полигона.
		if p.polygon.Contains(randomPoint) {
			return Point{Lat: lat, Lng: lng}
		}
	}
}

// CheckPointIsAllowed проверяет, находится ли точка в разрешенной зоне.
func CheckPointIsAllowed(point Point, allowedZone PolygonChecker, disabledZones []PolygonChecker) bool {
	if allowedZone.Contains(point) && allowedZone.Allowed() {
		for _, zone := range disabledZones {
			if zone.Contains(point) && zone.Allowed() {
				return false
			}
		}
		return true
	}
	return false
}

// GetRandomAllowedLocation возвращает случайную точку в разрешенной зоне.
func GetRandomAllowedLocation(allowedZone PolygonChecker, disabledZones []PolygonChecker) Point {
	for {
		randomPoint := allowedZone.RandomPoint()
		if CheckPointIsAllowed(randomPoint, allowedZone, disabledZones) {
			return randomPoint
		}
	}
}

// NewDisAllowedZone1 создает полигон для запрещенной зоны 1.
func NewDisAllowedZone1() (*Polygon, error) {
	// Замените points на координаты полигона для запрещенной зоны 1.
	points := []Point{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 0}, {Lat: 0, Lng: 0}} // Пример
	return NewPolygon(points, false)
}

// NewDisAllowedZone2 создает полигон для запрещенной зоны 2.
func NewDisAllowedZone2() (*Polygon, error) {
	// Замените points на координаты полигона для запрещенной зоны 2.
	points := []Point{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 0}, {Lat: 0, Lng: 0}} // Пример
	return NewPolygon(points, false)
}

// NewAllowedZone создает полигон для разрешенной зоны.
func NewAllowedZone() (*Polygon, error) {
	// Замените points на координаты полигона для разрешенной зоны.
	points := []Point{{Lat: 0, Lng: 0}, {Lat: 0, Lng: 0}, {Lat: 0, Lng: 0}} // Пример
	return NewPolygon(points, true)
}
