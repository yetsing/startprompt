package linemode

type LineMode string

//goland:noinspection ALL
const (
	Normal            LineMode = "normal"
	IncrementalSearch LineMode = "incremental-search"
	Complete          LineMode = "complete"
)

func (m LineMode) In(modes ...LineMode) bool {
	for _, mode := range modes {
		if m == mode {
			return true
		}
	}
	return false
}

func (m LineMode) Is(mode LineMode) bool {
	return m == mode
}
