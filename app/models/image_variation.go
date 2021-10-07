package models

type ImageVariation struct {
	Base
	MinioKey string `json:"-"`
	Etag     string `json:"-"`
	Width    uint   `json:"width"`
	Height   uint   `json:"height"`
	MimeType string `json:"mimeType"`
	ImageID  uint   `json:"-"`
}
