package offset

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type FileOffsetStorage struct {
	logger *slog.Logger
	path   string
}

func NewFileOffsetStorage(logger *slog.Logger, path string) (*FileOffsetStorage, error) {
	fos := &FileOffsetStorage{
		logger: logger,
		path:   path,
	}
	err := fos.init()
	if err != nil {
		logger.Error("Error of init file storage")
		return nil, err
	}
	logger.Info("File storage has been created")
	return fos, nil
}

func (fos *FileOffsetStorage) init() error {
	dir, _ := filepath.Split(fos.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fos.logger.Error("Error of creating dir")
		return err
	}

	file, err := os.OpenFile(fos.path, os.O_CREATE, 0655)
	if err != nil {
		fos.logger.Error("File open error")
		return err
	}
	fos.logger.Info("File storage has been initialized")
	return file.Close()
}

func (fos *FileOffsetStorage) Save(off []Offset) error {
	file, err := os.OpenFile(fos.path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fos.logger.Error("File open error")
		return err
	}

	data, err := json.Marshal(off)
	if err != nil {
		fos.logger.Error("Marshaling error")
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	return file.Close()
}

func (fos *FileOffsetStorage) Load() ([]Offset, error) {
	data, err := os.ReadFile(fos.path)
	if err != nil {
		return nil, err
	}
	off := []Offset{}
	err = json.Unmarshal(data, &off)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return off, nil
		}
		fos.logger.Error("Load offsets error")
		return nil, err
	}
	fos.logger.Info("Offsets successfully loaded")
	return off, nil
}