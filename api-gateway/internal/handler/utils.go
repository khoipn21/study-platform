package handler

import (
	"io"
)

// ReadAll reads all data from an io.Reader
func ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}