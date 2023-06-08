package utils

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

func CompressFiles(filenames []string, resultFilename string) error {
	archive, err := os.Create(resultFilename)
	if err != nil {
		return err
	}
	defer archive.Close()

	writer := zip.NewWriter(archive)
	for _, filename := range filenames {
		err := putFileInArchive(filename, writer)
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

func putFileInArchive(filename string, zipWriter *zip.Writer) error {
	// Open target file
	targetFile, err := os.Open(filename)
	if err != nil {
		return err
	}

	cleanFilename := path.Base(filename)

	// Init placeholder in zip
	archiveFile, err := zipWriter.Create(cleanFilename)
	if err != nil {
		return err
	}
	// Copy target file to zip
	_, err = io.Copy(archiveFile, targetFile)
	if err != nil {
		return err
	}

	return err
}
