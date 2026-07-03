package source

type SourceType string

const (
	TypeFile SourceType = "file"
	TypeDB SourceType = "database"
	TypeHTTP SourceType = "http"
)

type Source struct {
	Name string
	Type SourceType
	Path string
}

type Record struct {
	Source string
	Data []byte
}

type Scanner interface {
	Scan() ([]Source, error)
}

type Reader interface {
	Read(src Source, offset int64) ([]Record, int64, error)
}