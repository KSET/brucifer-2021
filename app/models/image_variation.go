package models

type ImageVariation struct {
	Base
	MinioKey string `json:"-" gorm:"size:255"`
	Etag     string `json:"-" gorm:"size:255"`
	Width    uint   `json:"width"`
	Height   uint   `json:"height"`
	MimeType string `json:"mimeType" gorm:"size:63"`
	ImageID  uint   `json:"-"`
}
