package payload

type UploadResponse struct {
	URL       string `json:"url"`
	FileName  string `json:"fileName"`
	MimeType  string `json:"mimeType"`
	Size      int64  `json:"size"`
	CreatedAt int64  `json:"createdAt"`
}
