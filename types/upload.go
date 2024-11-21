package types

import "os"

type FileType string

type Acl string

const (
	Text          FileType = "text/plain"
	CSV           FileType = "text/csv"
	Docx          FileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	PDF           FileType = "application/pdf"
	XLSX          FileType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ImageJPG      FileType = "image/jpeg"
	ImagePNG      FileType = "image/png"
	VideoMP4      FileType = "video/mp4"
	AudioMP3      FileType = "audio/mpeg"
	CompressedZIP FileType = "application/zip"
	XML           FileType = "application/xml"
	JSON          FileType = "application/json"
	Public        Acl      = "public-read"
	Private       Acl      = "private"
)

type PrepareUpload struct {
	FileName           string   `json:"fileName"`
	FileSize           float64  `json:"fileSize"`
	FileType           FileType `json:"fileType"`
	CustomId           string   `json:"customId"`
	ContentDisposition string   `json:"contentDisposition"`
	Acl                Acl      `json:"acl"`
	ExpiresIn          int      `json:"expiresIn"`
}

type PrepareUploadResponse struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

type UploadFile struct {
	File *os.File `json:"file"`
}

type UploadBody struct {
	File any `json:"file"`
}
