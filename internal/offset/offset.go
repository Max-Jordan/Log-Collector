package offset

type Offset struct {
	Source string `json:"source"`
	Offset int64 `json:"offset"`
}

type OffsetStorage interface {
	Save([]Offset) error
	Load() ([]Offset, error)
}

func NewOffset(src string, off int64) Offset {
	return Offset{Source: src, Offset: off}
}

func (o *Offset) SetOffset(off int64) {
	o.Offset = off
}

