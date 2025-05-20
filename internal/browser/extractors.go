package browser

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Extracts paragraph text from an HTTP response
func ExtractTextContent(body io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			sb.WriteString(text + "\n")
		}
	})

	return sb.String(), nil
}

// Extracts likely thumbnails from an HTTP response
func ExtractThumbnails(body io.Reader, base string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var images []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && isLikelyThumbnail(src) {
			images = append(images, resolveURL(src, base))
		}
	})

	return images, nil
}
