package database

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Database struct {
	ctx       context.Context
	logger    *slog.Logger
	persister SnapshotPersister
	data      map[string]int
	lock      sync.Mutex
}

func NewDatabase(
	ctx context.Context,
	logger *slog.Logger,
	persister SnapshotPersister,
) *Database {

	database := &Database{
		ctx:       ctx,
		logger:    logger,
		persister: persister,
		lock:      sync.Mutex{},
		data:      make(map[string]int),
	}

	database.LoadSnapshot()
	database.StartPeriodicFlush()

	return database
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

func (d *Database) LoadSnapshot() error {
	data, err := d.persister.LoadSnapshot()
	d.logger.Debug("Loaded snapshot")

	if err != nil {
		d.logger.Error("Failed to load snapshot", "error", err)
		return err
	}

	d.data = data
	return nil
}

func (d *Database) SaveSnapshot() error {
	data := d.Export()

	err := d.persister.SaveSnapshot(data)

	if err != nil {
		d.logger.Error("Failed to save snapshot", "error", err)
		return err
	}

	d.logger.Debug("Saved snapshot")

	return nil
}

func (d *Database) StartPeriodicFlush() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for {
			select {
			case <-ticker.C:
				d.SaveSnapshot()
			case <-d.ctx.Done():
				d.logger.Debug("Stopping flush thread...")
				return
			}
		}
	}()
}

func (d *Database) Close(ctx context.Context) error {
	var done = make(chan bool, 1)

	go func() {
		d.SaveSnapshot()
		done <- true
	}()

	for {
		select {
		case <-ctx.Done():
			d.logger.Warn("Snapshot save timed out")
			return ctx.Err()
		case <-done:
			d.logger.Debug("Snapshot save on close successful")
			return nil
		}
	}

}

func (d *Database) Export() map[string]int {
	d.lock.Lock()
	defer d.lock.Unlock()
	
	copy := make(map[string]int)

	for k, v := range d.data {
		copy[k] = v
	}

	return copy
}

func (d *Database) Import(data map[string]int) {
	d.lock.Lock()

	d.data = data

	defer d.lock.Unlock()
}
			