package models

type LengthCounter struct {
	Total int // Total # of bytes transferred
}

// Write implements the io.Writer interface.
// Always completes and never returns an error.
func (wc *LengthCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += n
	return n, nil
}
