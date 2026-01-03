package payload

import (
	"bbs-go/internal/upload"
	"fmt"
	"net/url"
	"path"
)

func GetAvatarURL(url string) string {
	return getThumbnailURL(url, 100, 100)
}

func GetSmallCoverURL(url string) string {
	return getThumbnailURL(url, 400, 300)
}

func GetPreviewURL(url string) string {
	return getThumbnailURL(url, 200, 200)
}

func getThumbnailURL(rawURL string, width, height int) string {
	if !upload.IsSUploadEnabeld() {
		return rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	filename := path.Base(u.Path)
	thumbSize := fmt.Sprintf("%dx%d", width, height)
	u.Path = path.Join(path.Dir(u.Path), thumbSize, filename)
	return u.String()
}
