package parser

import (
	"bytes"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractLinksAndResources(base *url.URL, html []byte) (links []*url.URL, resources []*url.URL) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if u := parseAndResolve(href, base); u != nil && isSameDomain(u, base) {
			links = append(links, u)
		}
	})

	doc.Find("img[src], script[src], link[href]").Each(func(i int, s *goquery.Selection) {
		var src string
		if s.Is("link") {
			rel, _ := s.Attr("rel")
			if !strings.Contains(rel, "stylesheet") {
				return
			}
			src, _ = s.Attr("href")
		} else {
			src, _ = s.Attr("src")
		}
		if u := parseAndResolve(src, base); u != nil {
			resources = append(resources, u)
		}
	})

	return
}

func parseAndResolve(href string, base *url.URL) *url.URL {
	href = strings.TrimSpace(href)
	if href == "" || strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "javascript:") {
		return nil
	}
	u, err := url.Parse(href)
	if err != nil {
		return nil
	}
	return base.ResolveReference(u)
}

func isSameDomain(u *url.URL, base *url.URL) bool {
	return u.Host == base.Host
}
