package payload

import (
	"bbs-go/internal/model"
	"log/slog"

	"bbs-go/common/jsons"
	"bbs-go/common/strs"
)

// Image
type ImageInfo struct {
	Url     string `json:"url"`
	Preview string `json:"preview"`
}

func BuildImageList(images []model.ImageDTO) (imageList []ImageInfo) {
	if len(images) > 0 {
		for _, image := range images {
			imageList = append(imageList, ImageInfo{
				Url:     image.URL,
				Preview: GetPreviewURL(image.URL),
			})
		}
	}
	return
}

func BuildImage(imageStr string) *ImageInfo {
	if strs.IsBlank(imageStr) {
		return nil
	}
	var img *model.ImageDTO
	if err := jsons.Parse(imageStr, &img); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return nil
	} else {
		return &ImageInfo{
			Url:     img.URL,
			Preview: GetPreviewURL(img.URL),
		}
	}
}
