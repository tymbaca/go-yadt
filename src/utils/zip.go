package utils

import (
	"archive/zip"
	"io"
	"os"
)

func CompressFiles(filenames []string, resultFilename string) error {
	archive, err := os.Create(resultFilename)
	if err != nil {
		return err
	}
	defer archive.Close()
	writer := zip.NewWriter(archive)

	for _, filename := range filenames {
		// Open target file
		targetFile, err := os.Open(filename)
		if err != nil {
			return err
		}
		// Init placeholder in zip
		archiveFile, err := writer.Create(filename)
		if err != nil {
			return err
		}
		// Copy target file to zip
		_, err = io.Copy(archiveFile, targetFile)
		if err != nil {
			return err
		}
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	return nil
}
