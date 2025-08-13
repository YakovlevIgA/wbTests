package service_test

import (
	"bytes"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"wbService/internal/repo"
	"wbService/internal/service"
)

func setupService() (service.Service, *repo.Storage) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	storage := repo.NewRepository().(*repo.Storage) // Для тестов нам нужен конкретный storage
	svc := service.NewService(storage, sugar)
	return svc, storage
}

// ---- Helper для выполнения HTTP-запроса ----
func performRequest(app *fiber.App, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	rec := httptest.NewRecorder()
	if resp != nil {
		rec.Code = resp.StatusCode
		_, _ = rec.Body.ReadFrom(resp.Body)
	}
	return rec
}

// ---- Тесты ----
func TestCreateEvent(t *testing.T) {
	svc, _ := setupService()
	app := fiber.New()
	app.Post("/event", svc.CreateEvent)

	body := map[string]interface{}{
		"description": "Test Event",
		"date":        "2025-08-13T10:00:00Z",
	}

	resp := performRequest(app, http.MethodPost, "/event", body)
	if resp.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.Code)
	}
}

func TestCreateEventBadRequest(t *testing.T) {
	svc, _ := setupService()
	app := fiber.New()
	app.Post("/event", svc.CreateEvent)

	// Некорректный JSON
	resp := performRequest(app, http.MethodPost, "/event", "bad json")
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.Code)
	}
}

func TestUpdateEventNotFound(t *testing.T) {
	svc, _ := setupService()
	app := fiber.New()
	app.Put("/event/:id", svc.UpdateEvent)

	body := map[string]interface{}{
		"description": "Updated Event",
		"date":        "2025-08-14T10:00:00Z",
	}

	resp := performRequest(app, http.MethodPut, "/event/999", body)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for not found, got %d", resp.Code)
	}
}

func TestDeleteEventNotFound(t *testing.T) {
	svc, _ := setupService()
	app := fiber.New()
	app.Delete("/event/:id", svc.DeleteEvent)

	resp := performRequest(app, http.MethodDelete, "/event/999", nil)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for not found, got %d", resp.Code)
	}
}

func TestGetEventsForDayInvalidDate(t *testing.T) {
	svc, _ := setupService()
	app := fiber.New()
	app.Get("/events/day/:day", svc.GetEventsForDay)

	resp := performRequest(app, http.MethodGet, "/events/day/invalid-date", nil)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid date, got %d", resp.Code)
	}
}

func TestGetEventsForWeekInvalidWeek(t *testing.T) {
	svc, _ := setupService()
	app := fiber.New()
	app.Get("/events/week/:week", svc.GetEventsForWeek)

	resp := performRequest(app, http.MethodGet, "/events/week/notanint", nil)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid week, got %d", resp.Code)
	}
}

func TestGetEventsForMonthInvalidMonth(t *testing.T) {
	svc, _ := setupService()
	app := fiber.New()
	app.Get("/events/month/:month", svc.GetEventsForMonth)

	resp := performRequest(app, http.MethodGet, "/events/month/99", nil)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid month, got %d", resp.Code)
	}
}
