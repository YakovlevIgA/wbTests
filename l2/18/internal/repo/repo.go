package repo

import (
	"fmt"
	"time"
)

func newStorage() *Storage {
	return &Storage{
		lib: make(map[int]*Event),
	}
}

type Repository interface {
	Create(event *Event) (int, error)
	Update(id int, event *Event) error
	Delete(id int) error
	GetEventsForDay(date time.Time) ([]*Event, error)
	GetEventsForWeek(date time.Time) ([]*Event, error)
	GetEventsForMonth(date time.Time) ([]*Event, error)
	GetEventsForWeekByNumber(week int) ([]*Event, error)   // новый метод
	GetEventsForMonthByNumber(month int) ([]*Event, error) // новый метод
}

func NewRepository() Repository {
	return newStorage()
}

func (s *Storage) Create(event *Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	event.CreatedAt = time.Now()
	s.lib[s.lastID] = event
	return s.lastID, nil
}

func (s *Storage) Update(id int, event *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.lib[id]; !ok {
		return ErrNotFound
	}
	event.CreatedAt = s.lib[id].CreatedAt
	s.lib[id] = event
	return nil
}

func (s *Storage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.lib[id]; !ok {
		return ErrNotFound
	}
	delete(s.lib, id)
	return nil
}

func (s *Storage) GetEventsForDay(date time.Time) ([]*Event, error) {
	return s.filterEvents(func(e *Event) bool {
		return sameDay(e.Date, date)
	})
}

func (s *Storage) GetEventsForWeek(date time.Time) ([]*Event, error) {
	year, week := date.ISOWeek()
	return s.filterEvents(func(e *Event) bool {
		y, w := e.Date.ISOWeek()
		return y == year && w == week
	})
}

func (s *Storage) GetEventsForMonth(date time.Time) ([]*Event, error) {
	return s.filterEvents(func(e *Event) bool {
		return e.Date.Year() == date.Year() && e.Date.Month() == date.Month()
	})
}

func (s *Storage) filterEvents(match func(*Event) bool) ([]*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*Event
	for _, ev := range s.lib {
		if match(ev) {
			result = append(result, ev)
		}
	}
	return result, nil
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

var ErrNotFound = fmt.Errorf("event not found")

func (s *Storage) GetEventsForWeekByNumber(week int) ([]*Event, error) {
	return s.filterEvents(func(e *Event) bool {
		_, w := e.Date.ISOWeek()
		return w == week
	})
}

func (s *Storage) GetEventsForMonthByNumber(month int) ([]*Event, error) {
	return s.filterEvents(func(e *Event) bool {
		return int(e.Date.Month()) == month
	})
}
