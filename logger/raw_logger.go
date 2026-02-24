package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type DecisionRawLog struct {
	Timestamp string `json:"timestamp"`
	Strategy  string `json:"strategy"`
	CoT       string `json:"cot"`
	Result    string `json:"result"`
}

var rawLoggerMutex sync.Mutex

// SaveRawDecision 异步保存 AI 输出和思维链到 HTML/前端 可读取的 YYYYMMDD_HH.json 文件
func SaveRawDecision(strategyName, cot, result string) {
	go func() {
		rawLoggerMutex.Lock()
		defer rawLoggerMutex.Unlock()

		now := time.Now()
		fileName := fmt.Sprintf("decisions_%s.json", now.Format("20060102_15"))
		dirPath := filepath.Join("data", "raw")

		if err := os.MkdirAll(dirPath, 0755); err != nil {
			Errorf("Failed to create raw data directory: %v", err)
			return
		}

		filePath := filepath.Join(dirPath, fileName)

		var logs []DecisionRawLog

		// Read existing file if it exists
		if data, err := os.ReadFile(filePath); err == nil {
			_ = json.Unmarshal(data, &logs)
		}

		// Append new log entry
		logs = append(logs, DecisionRawLog{
			Timestamp: now.Format(time.RFC3339),
			Strategy:  strategyName,
			CoT:       cot,
			Result:    result,
		})

		// Write back gracefully
		if newData, err := json.MarshalIndent(logs, "", "  "); err == nil {
			if err := os.WriteFile(filePath, newData, 0644); err != nil {
				Errorf("Failed to write raw decision log: %v", err)
			}
		}

		// GC old files, keeping max 12 hours rolling window
		cleanupOldRawLogs(dirPath, 12)
	}()
}

func cleanupOldRawLogs(dirPath string, maxFiles int) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "decisions_") && strings.HasSuffix(entry.Name(), ".json") {
			files = append(files, entry.Name())
		}
	}

	if len(files) <= maxFiles {
		return
	}

	// Because file format is decisions_YYYYMMDD_HH.json, default string sort acts as chronological sort
	sort.Strings(files)

	// Delete excessive old files
	filesToDelete := len(files) - maxFiles
	for i := 0; i < filesToDelete; i++ {
		os.Remove(filepath.Join(dirPath, files[i]))
	}
}
