package model

type UploadType string

type Upload struct {
	Model
	UserID    int64  `gorm:"not null;index"`                         // ID of the user who uploaded the file
	URL       string `gorm:"type:varchar(512);not null;uniqueIndex"` // Static URL for accessing the file
	FileName  string `gorm:"type:varchar(256);not null"`             // Original filename
	Size      int64  `gorm:"not null"`                               // Size of the file in bytes
	MimeType  string `gorm:"type:varchar(128);not null"`             // MIME type of the file
	Status    int    `gorm:"not null;default:0"`                     // Status of the file (e.g., active, deleted)
	CreatedAt int64  `gorm:"createdAt"`                              // Timestamp of when the file was uploaded
}
