package database

import (
	"log/slog"
	"sync"
)

type Database struct {
	logger  *slog.Logger
	data    map[string]int
	lock    sync.Mutex
}

func NewDatabase(logger *slog.Logger) *Database {
	return &Database{
		logger: logger,
		data:   make(map[string]int),
	}
}

func (d *Database) Get(key string) int {
	value, exists := d.data[key]

	if !exists {
		d.logger.Debug("[GET] key not found", "key", key)
		return 0
	}

	return value
}

func (d *Database) AddOccurences(key string, extraOccurences int) {
	d.lock.Lock()
	defer d.lock.Unlock()
	
	currentValue, exists := d.data[key]

	if !exists {
		currentValue = 0
	}
	
	newValue := currentValue + extraOccurences
	
	d.logger.Debug("[SET] Key", "key", key, "newValue", newValue)
	d.data[key] = newValue
}