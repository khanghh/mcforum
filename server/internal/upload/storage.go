package upload

import (
	"bbs-go/internal/config"
	"io"
)

type uploadStorage interface {
	// SaveStream saves data from reader to dstPath with the given MIME type.
	// Returns the public URL and an error, if any.
	SaveStream(dstPath string, reader io.Reader, mimeType string) (string, int, error)

	// SaveFile saves data to dstPath with the given MIME type.
	// Returns the public URL and an error, if any.
	SaveFile(dstPath string, data []byte, mimeType string) (string, error)

	// SaveImage reduces image quality, save images and thumbnails based on its mimeType.
	// Thumbnail saved as dstPath + "_thumb" + <extension>.
	// Returns the public URL and an error, if any.
	SaveImage(dstPath string, data []byte, mimeType string) (string, error)
}

var (
	local = &localFileStorage{}
)

func SaveImage(id string, name string, data []byte, mimeType string) (string, error) {
	dst := generateStorageKey(id, name, mimeType)
	return getStorage().SaveImage(dst, data, mimeType)
}

func SaveStream(id string, name string, reader io.Reader, mimeType string) (string, int, error) {
	dst := generateStorageKey(id, name, mimeType)
	return getStorage().SaveStream(dst, reader, mimeType)
}

func SaveFile(id string, name string, data []byte, mimeType string) (string, error) {
	dst := generateStorageKey(id, name, mimeType)
	return getStorage().SaveFile(dst, data, mimeType)
}

func SaveFromURL(id string, url string) (string, int, error) {
	stream, filename, mimeType, err := openDownloadStream(url)
	if err != nil {
		return "", 0, err
	}
	return SaveStream(id, filename, stream, mimeType)
}

func SaveImageFromURL(id string, url string) (string, error) {
	stream, filename, mimeType, err := openDownloadStream(url)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(stream)
	if err != nil {
		return "", err
	}
	return SaveImage(id, filename, data, mimeType)
}

func PutImage(data []byte, contentType string) (string, error) {
	return "", nil
}

func PutObject(key string, data []byte, contentType string) (string, error) {
	return "", nil
}

func getStorage() uploadStorage {
	return local
}

// IsSUploadEnabeld whether to enable Aliyun OSS
func IsSUploadEnabeld() bool {
	return config.Instance().Uploader.Enable == "supload"
}
