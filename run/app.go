package run

import (
	"context"
	"geotask/cache"
	"geotask/geo"
	cservice "geotask/module/courier/service"
	cstorage "geotask/module/courier/storage"
	"geotask/module/courierfacade/controller"
	cfservice "geotask/module/courierfacade/service"
	oservice "geotask/module/order/service"
	ostorage "geotask/module/order/storage"
	"geotask/router"
	"geotask/server"
	"geotask/workers/order"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	// получение хоста и порта redis
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	// инициализация клиента redis
	rclient := cache.NewRedisClient(host, port)

	// инициализация контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// проверка доступности redis
	_, err := rclient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	// инициализация разрешенной зоны
	allowedZone, _ := geo.NewAllowedZone()
	// инициализация запрещенных зон
	disAllowedZones := []geo.PolygonChecker{geo.NewDisAllowedZone1(), geo.NewDisAllowedZone2()}

	// инициализация хранилища заказов
	orderStorage := ostorage.NewOrderStorage(rclient)
	// инициализация сервиса заказов
	orderService := oservice.NewOrderService(orderStorage, allowedZone, disAllowedZones)

	orderGenerator := order.NewOrderGenerator(orderService)
	orderGenerator.Run()

	oldOrderCleaner := order.NewOrderCleaner(orderService)
	oldOrderCleaner.Run()

	// инициализация хранилища курьеров
	courierStorage := cstorage.NewCourierStorage(rclient)
	// инициализация сервиса курьеров
	courierSevice := cservice.NewCourierService(courierStorage, allowedZone, disAllowedZones)

	// инициализация фасада сервиса курьеров
	courierFacade := cfservice.NewCourierFacade(courierSevice, orderService)

	// инициализация контроллера курьеров
	courierController := controller.NewCourierController(courierFacade)

	// инициализация роутера
	routes := router.NewRouter(courierController)
	// инициализация сервера
	r := server.NewHTTPServer()
	// инициализация группы роутов
	api := r.Group("/api")
	// инициализация роутов
	routes.CourierAPI(api)

	mainRoute := r.Group("/")

	routes.Swagger(mainRoute)
	// инициализация статических файлов
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("public"))))

	// запуск сервера
	//serverPort := os.Getenv("SERVER_PORT")

	if os.Getenv("ENV") == "prod" {
		certFile := "/app/certs/cert.pem"
		keyFile := "/app/certs/private.pem"
		return r.RunTLS(":443", certFile, keyFile)
	}

	return r.Run()
}
