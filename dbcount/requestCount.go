package dbcount

import (
	"sync"
)

var (
	readCounts  = make(map[string]int)
	writeCounts = make(map[string]int)
	mu          sync.Mutex
)

func InitializeRequestCounts() {
	mu.Lock()
	defer mu.Unlock()
	readCounts = make(map[string]int)
	writeCounts = make(map[string]int)
}

func IncrementRequestCount(tableName string, isWrite bool) {
	mu.Lock()
	defer mu.Unlock()
	if isWrite {
		writeCounts[tableName]++
	} else {
		readCounts[tableName]++
	}
}

func GetReadCount(tableName string) int {
	mu.Lock()
	defer mu.Unlock()
	return readCounts[tableName]
}

func GetWriteCount(tableName string) int {
	mu.Lock()
	defer mu.Unlock()
	return writeCounts[tableName]
}

func PrintRequestCounts() {
	mu.Lock()
	defer mu.Unlock()
	for tableName, count := range readCounts {
		println("Read", tableName, count)
	}
	for tableName, count := range writeCounts {
		println("Write", tableName, count)
	}
}
