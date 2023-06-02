package generator

import (
	"github.com/lukasjarosch/go-docx"
)

type FileGenerator struct {
	TempateFilename string
	activeTemplate  *docx.Document
	data            ParseData
}

type ParseData []FileData

type FileData struct {
	Filename string                `json:"filename"`
	Pages    []docx.PlaceholderMap `json:"pages"`
}
