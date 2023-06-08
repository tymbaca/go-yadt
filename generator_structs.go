package yadt

import (
	"github.com/lukasjarosch/go-docx"
)

type FileGenerator struct {
	TempateFilename string
	Data            *ParseData

	activeTemplate *docx.Document
	filenames      []string
	tmpDirectory   string
}

type ParseData []FileData

type FileData struct {
	Filename string                `json:"filename"`
	Pages    []docx.PlaceholderMap `json:"pages"`
}
