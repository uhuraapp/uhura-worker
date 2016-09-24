package parser

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func FindLinks(body []byte) []string {
	links := make([]string, 0)

	r := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(r)
	if err == nil {
		doc.Find("link[rel=alternate][type='application/rss+xml']").Each(func(_ int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok {
				links = append(links, href)
			}
		})
	}

	return links
}
