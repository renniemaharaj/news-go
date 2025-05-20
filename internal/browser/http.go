package browser

import (
	"bytes"
	"io"
)

// DeepCopyBody reads a body and returns n independent ReadClosers
func DeepCopyBody(body io.ReadCloser, n int) ([]io.ReadCloser, error) {
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	copies := make([]io.ReadCloser, n)
	for i := 0; i < n; i++ {
		copies[i] = io.NopCloser(bytes.NewReader(data))
	}

	return copies, nil
}
