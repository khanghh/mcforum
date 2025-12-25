package bbsurls

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"bbs-go/common/base62"
	"bbs-go/internal/config"
)

// Whether it is an internal link
func IsInternalUrl(href string) bool {
	if IsAnchor(href) {
		return true
	}
	u, err := url.Parse(config.Instance().BaseUrl)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return false
	}
	return strings.Contains(href, u.Host)
}

// Whether it is an anchor link
func IsAnchor(href string) bool {
	return strings.Index(href, "#") == 0
}

func AbsUrl(path string) string {
	return config.Instance().BaseUrl + path
}

func UserUrl(username string) string {
	return "/u/" + username
}

// User homepage
func AbsUserUrl(username string) string {
	return AbsUrl(UserUrl(username))
}
func TagUrl(tagName string) string {
	return "/tags/" + tagName
}

// Tag article list
func AbsTagUrl(tagName string) string {
	return AbsUrl(TagUrl(tagName))
}

func TopicUrl(slug string, id int64) string {
	return fmt.Sprintf("/t/%s.%s", slug, base62.Encode(id))
}

// Topic details
func AbsTopicUrl(slug string, id int64) string {
	return AbsUrl(TopicUrl(slug, id))
}

func UrlJoin(parts ...string) string {
	sep := "/"
	var ss []string
	for i, part := range parts {
		part = strings.TrimSpace(part)
		var (
			from = 0
			to   = len(part)
		)
		if strings.Index(part, sep) == 0 {
			from = 1
		}
		if strings.LastIndex(part, sep) == len(part)-1 {
			to = len(part) - 1
		}
		part = part[from:to]

		ss = append(ss, part)
		if i != len(parts)-1 {
			ss = append(ss, sep)
		}
	}
	return strings.Join(ss, "")
}
