package offset

type OffsetStorage interface {
	Save([]Offset) error
	Load() ([]Offset, error)
}