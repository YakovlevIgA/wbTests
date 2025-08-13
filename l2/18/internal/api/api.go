package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"wbService/cmd/middleware" // todo: check
	"wbService/internal/service"
)

// Routers - структура для хранения зависимостей роутов
type Routers struct {
	Service service.Service
}

// NewRouters - конструктор для настройки API
func NewRouters(r *Routers, log *zap.SugaredLogger) *fiber.App {
	app := fiber.New()
	app.Use(middleware.RequestLogger(log))
	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET, POST, DELETE",
		AllowHeaders:     "Accept, Content-Type, X-REQUEST-SomeID",
		ExposeHeaders:    "Link",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Группа маршрутов с авторизацией
	apiGroup := app.Group("/v1")

	// Роут для создания задачи
	apiGroup.Post("/create_event", r.Service.CreateEvent)
	apiGroup.Post("/update_event/:id", r.Service.UpdateEvent)
	apiGroup.Delete("/events/:id", r.Service.DeleteEvent)
	apiGroup.Get("/events_for_day/:day", r.Service.GetEventsForDay)
	apiGroup.Get("/events_for_week/:week", r.Service.GetEventsForWeek)
	apiGroup.Get("/events_for_month/:month", r.Service.GetEventsForMonth)

	return app
}
