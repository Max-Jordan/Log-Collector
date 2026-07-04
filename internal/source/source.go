package source

import ad "github.com/Max-Jordan/Log-Collector/internal/adapter"

type Source = string

type SourcesScanner interface {
	ScanSources()
}

type LogReader interface {
	ReadLog(src Source, offset int64) ([]ad.LogEntity, int64, error)
}