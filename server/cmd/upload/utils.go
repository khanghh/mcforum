package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"bbs-go/common/dates"
	"bbs-go/common/strs"
)

// generateStorageKey generates a storage key for the uploaded file
func generateStorageKey(id string, name string, mimeType string) string {
	ext := ""
	if strs.IsNotBlank(mimeType) {
		exts, err := mime.ExtensionsByType(mimeType)
		if err == nil && len(exts) > 0 {
			ext = exts[0]
		}
	}

	fileName := name
	if fileExt := filepath.Ext(fileName); fileExt != ext {
		fileName = fileName[:len(fileName)-len(fileExt)] + ext
	}

	ret, _ := url.JoinPath("", dates.Format(time.Now(), "2006/01/02/"), id, fileName)
	return ret
}

// openDownloadStream opens a download stream from the given URL and returns:
// - stream: io.ReadCloser
// - filename: suggested name from header or URL
// - mimetype: from Content-Type header
func openDownloadStream(rawURL string) (io.ReadCloser, string, string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, "", "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get MIME type
	mimeType := resp.Header.Get("Content-Type")

	// Get filename from Content-Disposition header
	cd := resp.Header.Get("Content-Disposition")
	var filename string
	if cd != "" {
		const prefix = "filename="
		if idx := strings.Index(cd, prefix); idx != -1 {
			filename = strings.Trim(cd[idx+len(prefix):], `"`)
		}
	}

	// Fallback to last URL segment
	if filename == "" {
		u, err := url.Parse(rawURL)
		if err != nil {
			resp.Body.Close()
			return nil, "", "", err
		}
		filename = path.Base(u.Path)
	}

	// Fallback MIME type if missing
	if mimeType == "" {
		mimeType = mime.TypeByExtension(path.Ext(filename))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}

	return resp.Body, filename, mimeType, nil
}
