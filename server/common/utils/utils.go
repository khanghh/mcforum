package utils

import (
	"bbs-go/common/strs"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetRequestIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
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

// GetHtmlText 获取html文本
func GetHtmlText(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return ""
	}
	return doc.Text()
}

func GetSummaryText(s string, length int) string {
	s = strings.TrimSpace(s)
	summary := strs.Substr(s, 0, length)
	if strs.RuneLen(s) > length {
		summary += "..."
	}
	return summary
}

func GetSummaryHtml(htmlStr string, summaryLen int) string {
	if summaryLen <= 0 || strs.IsEmpty(htmlStr) {
		return ""
	}
	return GetSummaryText(GetHtmlText(htmlStr), summaryLen)
}
