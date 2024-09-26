package database

import (
	"log/slog"
	"testing"
)

func TestDatabaseUnusedWords(t *testing.T) {
	d := NewDatabase(slog.Default())
	
	d.AddOccurences("hello", 1)
	
	if d.Get("world") != 0 {
		t.Error("expected 0, got", d.Get("world"))
	}
}

func TestDatabaseSetAndGet(t *testing.T) {
	d := NewDatabase(slog.Default())
	
	d.AddOccurences("hello", 1)
	d.AddOccurences("hello", 2)
	d.AddOccurences("hello", 3)
	d.AddOccurences("hello", 4)
	
	if d.Get("hello") != 10 {
		t.Error("expected 10, got", d.Get("hello"))
	}
}