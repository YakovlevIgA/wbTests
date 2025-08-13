package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"time"

	"wbService/internal/dto"
	"wbService/internal/repo"
)

type Service interface {
	CreateEvent(ctx *fiber.Ctx) error
	UpdateEvent(ctx *fiber.Ctx) error
	DeleteEvent(ctx *fiber.Ctx) error
	GetEventsForDay(ctx *fiber.Ctx) error
	GetEventsForWeek(ctx *fiber.Ctx) error
	GetEventsForMonth(ctx *fiber.Ctx) error
}

type service struct {
	repo repo.Repository
	log  *zap.SugaredLogger
}

func NewService(r repo.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: r,
		log:  logger,
	}
}

func (s *service) CreateEvent(ctx *fiber.Ctx) error {
	var req EventRequest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Errorw("Invalid request body", "error", err)
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "invalid request body")
	}

	event := &repo.Event{
		Description: req.Description,
		Date:        req.Date,
	}

	id, err := s.repo.Create(event)
	if err != nil {
		s.log.Errorw("Failed to create event", "error", err)
		return dto.InternalServerError(ctx)
	}

	return dto.SuccessResponse(ctx, map[string]int{"event_id": id})
}

func (s *service) UpdateEvent(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Errorw("Invalid ID format", "error", err, "id", idStr)
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "invalid id format, expected integer")
	}

	var req EventUpdateRequest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Errorw("Invalid request body", "error", err)
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "invalid request body")
	}

	event := &repo.Event{
		Description: req.Description,
		Date:        req.Date,
	}

	if err := s.repo.Update(id, event); err != nil {
		s.log.Errorw("Failed to update event", "id", id, "error", err)
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "event not found")
	}

	return dto.SuccessResponse(ctx, map[string]int{"updated_id": id})
}

func (s *service) DeleteEvent(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Errorw("Invalid ID format", "error", err, "id", idStr)
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "invalid id format, expected integer")
	}

	if err := s.repo.Delete(id); err != nil {
		s.log.Errorw("Failed to delete event", "id", id, "error", err)
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "event not found")
	}

	return dto.SuccessResponse(ctx, map[string]int{"deleted_id": id})
}

func (s *service) GetEventsForDay(ctx *fiber.Ctx) error {
	dayStr := ctx.Params("day") // например "2025-08-13"
	date, err := time.Parse("2006-01-02", dayStr)
	if err != nil {
		s.log.Errorw("Invalid date format", "error", err, "day", dayStr)
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "invalid date format, expected YYYY-MM-DD")
	}

	events, err := s.repo.GetEventsForDay(date)
	if err != nil {
		s.log.Errorw("Failed to get events for day", "error", err)
		return dto.InternalServerError(ctx)
	}

	return dto.SuccessResponse(ctx, events)
}

func (s *service) GetEventsForWeek(ctx *fiber.Ctx) error {
	weekStr := ctx.Params("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		s.log.Errorw("Invalid week format", "error", err, "week", weekStr)
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "invalid week format, expected integer")
	}

	events, err := s.repo.GetEventsForWeekByNumber(week)
	if err != nil {
		s.log.Errorw("Failed to get events for week", "error", err)
		return dto.InternalServerError(ctx)
	}

	return dto.SuccessResponse(ctx, events)
}

func (s *service) GetEventsForMonth(ctx *fiber.Ctx) error {
	monthStr := ctx.Params("month")
	monthInt, err := strconv.Atoi(monthStr)
	if err != nil || monthInt < 1 || monthInt > 12 {
		s.log.Errorw("Invalid month format", "error", err, "month", monthStr)
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "invalid month format, expected 1-12")
	}

	events, err := s.repo.GetEventsForMonthByNumber(monthInt)
	if err != nil {
		s.log.Errorw("Failed to get events for month", "error", err)
		return dto.InternalServerError(ctx)
	}

	return dto.SuccessResponse(ctx, events)
}
