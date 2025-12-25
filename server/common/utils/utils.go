package utils

import (
	"bbs-go/common/strs"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetRequestIP attempts to determine the client's IP address.
// It parses X-Real-IP and X-Forwarded-For to support reverse proxies (nginx/haproxy).
func GetRequestIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func GetUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

// GetHtmlText extracts plain text from an HTML string.
func GetHtmlText(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return ""
	}
	return doc.Text()
}

// GetSummaryText returns a trimmed summary of s up to length runes.
func GetSummaryText(s string, length int) string {
	s = strings.TrimSpace(s)
	summary := strs.Substr(s, 0, length)
	if strs.RuneLen(s) > length {
		summary += "..."
	}
	return summary
}

// GetSummaryHtml extracts text from HTML and returns a summary of the given length.
func GetSummaryHtml(htmlStr string, summaryLen int) string {
	if summaryLen <= 0 || strs.IsEmpty(htmlStr) {
		return ""
	}
	return GetSummaryText(GetHtmlText(htmlStr), summaryLen)
}
