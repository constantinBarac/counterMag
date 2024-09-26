package database

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SnapshotPersister interface {
	SaveSnapshot(data map[string]int) error
	LoadSnapshot() (map[string]int, error)
}

type FileSnapshotPersister struct {
	logPath string
}

func (p *FileSnapshotPersister) SaveSnapshot(data map[string]int) error {
	logFile, err := os.OpenFile(p.logPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()

	if err != nil {
		return err
	}

	buffer := make([]string, 0, 50)
	for k, v := range data {
		buffer = append(buffer, fmt.Sprintf("%s %d", k, v))
		mergedRecords := strings.Join(buffer, "\n")
		logFile.WriteString(mergedRecords)
		buffer = make([]string, 0, 50)
	}

	mergedRecords := strings.Join(buffer, "\n")
	logFile.WriteString(mergedRecords)

	return nil
}

func (p *FileSnapshotPersister) LoadSnapshot() (map[string]int, error) {
	logFile, err := os.OpenFile(p.logPath, os.O_RDONLY, 0644)
	defer logFile.Close()

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(logFile)

	data := make(map[string]int)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		tokens := strings.Split(line, " ")
		key := tokens[0]
		value, _ := strconv.Atoi(tokens[1])
		data[key] = value
	}

	return data, nil
}

type MockPersister struct {
	logPath string
	data    map[string]int
}

func (p *MockPersister) SaveSnapshot(data map[string]int) error {
	p.data = data
	return nil
}

func (p *MockPersister) LoadSnapshot() (map[string]int, error) {
	return p.data, nil
}