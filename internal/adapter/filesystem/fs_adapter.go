package filesystem

import (
	"bufio"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	ad "github.com/Max-Jordan/Log-Collector/internal/adapter"
	"github.com/Max-Jordan/Log-Collector/internal/offset"
	src "github.com/Max-Jordan/Log-Collector/internal/source"
)

const scanInterval = 5 * time.Second

type fsAdapter struct {
	logger        *slog.Logger
	offsetStorage offset.OffsetStorage
	dir           string
}

type ReadResult struct {
	Source src.Source
	Logs   []ad.LogEntity
	Offset int64
	Err    error
}

func NewFSAdapter(logger *slog.Logger, offsetStorage offset.OffsetStorage, dir string) *fsAdapter {
	return &fsAdapter{
		logger:        logger,
		offsetStorage: offsetStorage,
		dir:           dir,
	}
}

func (fs fsAdapter) ScanSources() {
	fs.scanSourcesOnce()

	ticker := time.NewTicker(scanInterval)
	defer ticker.Stop()

	for range ticker.C {
		fs.scanSourcesOnce()
	}
}

func (fs fsAdapter) scanSourcesOnce() {
	files, err := os.ReadDir(fs.dir)
	if err != nil {
		fs.logger.Error("read directory", slog.String("error", err.Error()))
		return
	}

	offsets := make(map[string]int64)
	for _, file := range files {
		if file.Type().IsRegular() {
			offset, err := fs.offsetForSource(file.Name())
			if err != nil {
				fs.logger.Error("load offset", slog.String("source", file.Name()), slog.String("error", err.Error()))
				continue
			}
			offsets[file.Name()] = offset
		}
	}

	fs.processSources(offsets)
}

func (fs fsAdapter) ReadLog(source src.Source, offsetValue int64) ([]ad.LogEntity, int64, error) {
	filePath := filepath.Join(fs.dir, string(source))
	file, err := os.Open(filePath)
	if err != nil {
		fs.logger.Error("open log file", slog.String("file", filePath), slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer file.Close()

	if _, err := file.Seek(offsetValue, io.SeekStart); err != nil {
		return nil, 0, err
	}

	logs := make([]ad.LogEntity, 0)
	nextOffset := offsetValue
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		nextOffset += int64(len(line)) + 1

		if len(line) == 0 {
			continue
		}

		var log ad.LogEntity
		if err := json.Unmarshal(line, &log); err != nil {
			fs.logger.Error("parse log line", slog.String("file", filePath), slog.String("error", err.Error()))
			continue
		}

		logs = append(logs, log)
	}

	if err := scanner.Err(); err != nil {
		return nil, offsetValue, err
	}

	return logs, nextOffset, nil
}

func (fs fsAdapter) ReadLogs(sources map[string]int64) <-chan ReadResult {
	results := make(chan ReadResult, len(sources))
	go func() {
		defer close(results)

		var wg sync.WaitGroup
		wg.Add(len(sources))
		for sourceName, offsetValue := range sources {
			source := src.Source(sourceName)
			offsetValue := offsetValue

			go func() {
				defer wg.Done()
				logs, nextOffset, err := fs.ReadLog(source, offsetValue)
				results <- ReadResult{
					Source: source,
					Logs:   logs,
					Offset: nextOffset,
					Err:    err,
				}
			}()
		}

		wg.Wait()
	}()

	return results
}

func (fs fsAdapter) processSources(sources map[string]int64) {
	for result := range fs.ReadLogs(sources) {
		if result.Err != nil {
			fs.logger.Error(
				"read source",
				slog.String("source", string(result.Source)),
				slog.String("error", result.Err.Error()),
			)
			continue
		}
		if len(result.Logs) == 0 {
			continue
		}
		fs.logger.Info(
			"read logs",
			slog.String("source", string(result.Source)),
			slog.Int("count", len(result.Logs)),
			slog.Int64("offset", result.Offset),
		)
		if err := fs.SaveOffset(result.Source, result.Offset); err != nil {
			fs.logger.Error(
				"save offset",
				slog.String("source", string(result.Source)),
				slog.String("error", err.Error()),
			)
			continue
		}
	}
}

func (fs fsAdapter) offsetForSource(src src.Source) (int64, error) {
	offsets, err := fs.offsetStorage.Load()
	if err != nil {
		return 0, err
	}

	for _, off := range offsets {
		if off.Source == string(src) {
			return off.Offset, nil
		}
	}

	return 0, nil
}

func (fs fsAdapter) SaveOffset(source src.Source, offsetValue int64) error {
	offsets, err := fs.offsetStorage.Load()
	if err != nil {
		return err
	}

	for i := range offsets {
		if offsets[i].Source == string(source) {
			offsets[i].SetOffset(offsetValue)
			return fs.offsetStorage.Save(offsets)
		}
	}

	offsets = append(offsets, offset.NewOffset(string(source), offsetValue))
	return fs.offsetStorage.Save(offsets)
}
