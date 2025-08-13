package service

import "time"

type EventRequest struct {
	Description string    `json:"description" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}

type EventUpdateRequest struct {
	ID          int       `json:"id" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}

type IDRequest struct {
	ID int `json:"id" validate:"required"`
}

type DateRequest struct {
	Date time.Time `json:"date" validate:"required"`
}
