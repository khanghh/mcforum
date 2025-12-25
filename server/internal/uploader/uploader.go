package uploader

import (
	"strings"

	"bbs-go/common/strs"
	"bbs-go/common/urls"

	"bbs-go/internal/config"
)

type uploader interface {
	PutImage(data []byte, contentType string) (string, error)
	PutObject(key string, data []byte, contentType string) (string, error)
	CopyImage(originUrl string) (string, error)
}

var (
	local = &localUploader{}
)

func PutImage(data []byte, contentType string) (string, error) {
	return getUploader().PutImage(data, contentType)
}

func PutObject(key string, data []byte, contentType string) (string, error) {
	return getUploader().PutObject(key, data, contentType)
}

func CopyImage(url string) (string, error) {
	u1 := urls.ParseUrl(url).GetURL()
	u2 := urls.ParseUrl(config.Instance().BaseUrl).GetURL()
	// local host, no download
	if u1.Host == u2.Host {
		return url, nil
	}
	return getUploader().CopyImage(url)
}

func getUploader() uploader {
	return local
}

// IsEnabledOss whether to enable Aliyun OSS
func IsEnabledOss() bool {
	enable := config.Instance().Uploader.Enable
	return strs.EqualsIgnoreCase(enable, "aliyun") || strs.EqualsIgnoreCase(enable, "oss") ||
		strs.EqualsIgnoreCase(enable, "aliyunOss")
}

// IsOssImageUrl whether the image is stored in Aliyun OSS
func IsOssImageUrl(url string) bool {
	host := urls.ParseUrl(config.Instance().Uploader.AliyunOss.Host).GetURL().Host
	return strings.Contains(url, host)
}
