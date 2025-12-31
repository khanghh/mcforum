package payload

import (
	"bbs-go/internal/upload"
)

func HandleOssImageStyleAvatar(url string) string {
	if !upload.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, "")
}

func HandleOssImageStyleDetail(url string) string {
	if !upload.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, "")
}

func HandleOssImageStyleSmall(url string) string {
	if !upload.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, "")
}

func HandleOssImageStylePreview(url string) string {
	if !upload.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, "")
}

func HandleOssImageStyle(url, style string) string {
	return url
}
