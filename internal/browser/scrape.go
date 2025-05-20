package browser

import (
	"bytes"
	"io"
	"net/http"
)

// Combines both text and thumbnail extraction safely
func Scrape(resp *http.Response) (string, []string, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	resp.Body.Close()

	base := resp.Request.URL.String()
	body1 := io.NopCloser(bytes.NewReader(bodyBytes))
	body2 := io.NopCloser(bytes.NewReader(bodyBytes))

	text, err := ExtractTextContent(body1)
	if err != nil {
		return "", nil, err
	}

	thumbnails, err := ExtractThumbnails(body2, base)
	if err != nil {
		return "", nil, err
	}

	return text, thumbnails, nil
}
