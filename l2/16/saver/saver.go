package saver

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Save(u *url.URL, data []byte, contentType string, outDir string) string {
	contentType = strings.ToLower(strings.TrimSpace(contentType))

	localPath := urlToPath(u, contentType, outDir)

	if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
		fmt.Printf("Ошибка создания директории: %v\n", err)
		return ""
	}

	if strings.HasPrefix(contentType, "text/html") {
		data = rewriteLinks(u, data, outDir)
	}

	if err := os.WriteFile(localPath, data, 0644); err != nil {
		fmt.Printf("Ошибка сохранения %s: %v\n", localPath, err)
		return ""
	}

	return localPath
}

func urlToPath(u *url.URL, contentType, outDir string) string {
	h := md5.Sum([]byte(u.String()))
	hash := hex.EncodeToString(h[:])[:8]

	ext := extFromType(contentType)
	if ext == "dat" {
		ext = extFromURL(u)
	}

	path := u.Path
	if path == "" || strings.HasSuffix(path, "/") {
		path += "index"
	}

	filename := filepath.Base(path)

	if !strings.Contains(filename, ".") {
		filename += "." + ext
	} else {
		extFromName := filepath.Ext(filename)
		filename = strings.TrimSuffix(filename, extFromName)
		filename += "." + ext
	}

	dir := filepath.Join(outDir, u.Host, filepath.Dir(path))

	return filepath.Join(dir, fmt.Sprintf("%s_%s.%s", strings.TrimSuffix(filename, "."+ext), hash, ext))
}
func extFromType(contentType string) string {
	ct := strings.ToLower(contentType)
	switch {
	case strings.HasPrefix(ct, "text/html"):
		return "html"
	case strings.HasPrefix(ct, "text/css"):
		return "css"
	case strings.HasPrefix(ct, "application/javascript"), strings.HasPrefix(ct, "text/javascript"):
		return "js"
	case strings.HasPrefix(ct, "image/jpeg"):
		return "jpg"
	case strings.HasPrefix(ct, "image/png"):
		return "png"
	case strings.HasPrefix(ct, "image/gif"):
		return "gif"
	default:
		return "dat"
	}
}

func extFromURL(u *url.URL) string {
	ext := filepath.Ext(u.Path)
	if ext != "" {
		return strings.TrimPrefix(ext, ".")
	}
	return "dat"
}

func rewriteLinks(base *url.URL, html []byte, outDir string) []byte {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		fmt.Printf("Ошибка парсинга HTML для переписи ссылок: %v\n", err)
		return html
	}

	attrs := map[string]string{
		"a":      "href",
		"img":    "src",
		"script": "src",
		"link":   "href",
	}

	baseDir := filepath.Dir(urlToPath(base, "text/html", outDir))

	for tag, attr := range attrs {
		doc.Find(tag + "[" + attr + "]").Each(func(i int, s *goquery.Selection) {
			val, _ := s.Attr(attr)
			if val == "" {
				return
			}
			u, err := base.Parse(val)
			if err != nil {
				return
			}
			localPath := urlToPath(u, "", outDir)
			relPath, err := filepath.Rel(baseDir, localPath)
			if err != nil {
				return
			}
			relPath = filepath.ToSlash(relPath)
			s.SetAttr(attr, relPath)
		})
	}

	htmlStr, err := doc.Html()
	if err != nil {
		fmt.Printf("Ошибка при генерации HTML: %v\n", err)
		return html
	}
	return []byte(htmlStr)
}
