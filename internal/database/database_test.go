package database

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func getTestDatabase(persister SnapshotPersister) *Database {
	return &Database{
		ctx: context.Background(),
		logger: slog.Default(),
		persister: persister,
		data: make(map[string]int),
	}
}

func TestDatabaseUnusedWords(t *testing.T) {
	persister := MockPersister{}
	d := getTestDatabase(&persister)

	d.AddOccurences("hello", 1)

	if d.Get("world") != 0 {
		t.Error("expected 0, got", d.Get("world"))
	}
}


func TestDatabaseSetAndGet(t *testing.T) {
	persister := MockPersister{}
	d := getTestDatabase(&persister)

	d.AddOccurences("hello", 1)
	d.AddOccurences("hello", 2)
	d.AddOccurences("hello", 3)
	d.AddOccurences("hello", 4)

	if d.Get("hello") != 10 {
		t.Error("expected 10, got", d.Get("hello"))
	}
}

func TestDatabaseSaveLoadSnapshot(t *testing.T) {
	persister := MockPersister{}
	d := getTestDatabase(&persister)

	d.AddOccurences("hello", 1)
	d.AddOccurences("hello", 2)
	d.AddOccurences("hello", 3)
	d.AddOccurences("hello", 4)

	if err := d.SaveSnapshot(); err != nil {
		t.Error(err)
	}

	if persister.data["hello"] != 10 {
		t.Error("[PERSISTER] expected 10, got", persister.data["hello"])
	}

	recoveredDatabase := getTestDatabase(&persister)
	if err := recoveredDatabase.LoadSnapshot(); err != nil {
		t.Error(err)
	}

	if recoveredDatabase.Get("hello") != 10 {
		t.Error("expected 10, got", d.Get("hello"))
	}
}
