package repo

import (
	"sync"
	"time"
)

type Event struct {
	Description string
	CreatedAt   time.Time
	Date        time.Time
}

type Storage struct {
	lib    map[int]*Event
	mu     sync.RWMutex
	lastID int
}
