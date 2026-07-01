package offset

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)



type FileOffsetStorage struct {
	logger *slog.Logger
	path string
}

func NewFileOffsetStorage(logger *slog.Logger, path string) (*FileOffsetStorage, error) {
	fos := &FileOffsetStorage{
		logger: logger,
		path: path,
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
	file, err := os.OpenFile(fos.path, os.O_RDWR, 0644)
	if err != nil {
		fos.logger.Error("File open error")
		return err
	}

	data, err := json.Marshal(off)
	if err != nil {
		fos.logger.Error("Marshaling error")
		return err
	}
	file.Write(data)
	return file.Close()
}

func (fos *FileOffsetStorage) Load() ([]Offset, error) {
	data, err := os.ReadFile(fos.path)
	if err != nil {
		fmt.Println(err)
	}
	off := []Offset{}
	err = json.Unmarshal(data, &off)
	if err != nil {
		fos.logger.Error("Load offsers error")
		return nil, err
	}
	fos.logger.Info("Offsets successfully loaded")
	return off, nil
}

