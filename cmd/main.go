package main

import (
	"fmt"
	"log/slog"
	"os"

	fs "github.com/Max-Jordan/Log-Collector/internal/offset"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fos, err := fs.NewFileOffsetStorage(logger, "offset/offset.json")
	if err != nil {
		fmt.Println(err)
	}
	offs := []fs.Offset{fs.NewOffset("log-collector", 1), fs.NewOffset("second-collctor", 100)}
	fos.Save(offs)
	fmt.Println(fos.Load())
}
