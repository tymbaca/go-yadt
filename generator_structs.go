package yadt

import (
	"github.com/lukasjarosch/go-docx"
)

type FileGenerator struct {
	templateBytes []byte
	data          *parseData

	activeTemplate *docx.Document
	filenames      []string
	tmpDirectory   string
}

type parseData []fileData

type fileData struct {
	Filename string                `json:"filename", binding:"required"`
	Pages    []docx.PlaceholderMap `json:"pages", binding:"required"`
}
