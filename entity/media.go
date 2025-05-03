package entity

import (
	"mime/multipart"
)

type Media struct {
	Key    string                `json:"-"`
	Url    string                `json:"url"`
	Base64 string                `json:"base64"`
	File   multipart.File        `json:"-"`
	Header *multipart.FileHeader `json:"-"`
}
