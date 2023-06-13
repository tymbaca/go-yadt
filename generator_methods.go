package yadt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"

	"golang.org/x/sync/errgroup"

	"github.com/tymbaca/go-yadt/utils"

	"github.com/lukasjarosch/go-docx"
)

var err error

func New(templateStream io.Reader, jsonStream io.Reader) (*FileGenerator, error) {
	templateBytes, err := utils.StreamToBytes(templateStream)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := utils.StreamToBytes(jsonStream)
	if err != nil {
		return nil, err
	}
	fileGenerator, err := NewFromBytes(templateBytes, jsonBytes)
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

func NewFromBytes(templateBytes []byte, jsonBytes []byte) (*FileGenerator, error) {
	fileGenerator := new(FileGenerator)

	fileGenerator.templateBytes = templateBytes
	fileGenerator.data, err = parseJsonToData(jsonBytes)
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

func NewFromFiles(templateFilename string, jsonFilename string) (*FileGenerator, error) {

	templateBytes, err := os.ReadFile(templateFilename)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := os.ReadFile(jsonFilename)
	if err != nil {
		return nil, err
	}

	fileGenerator, err := NewFromBytes(templateBytes, jsonBytes)
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

func (s *FileGenerator) GenerateZip(filename string) error {
	var err error
	// s.tmpDirectory, err = os.MkdirTemp("./fixing/", "fixing-yadt") // TODO CHANGE MkdirTemp PARAMETERS TO ("", "") AFTER FIX
	// log.Printf("Created tmp directory: %s", s.tmpDirectory)
	// if err != nil {
	// 	panic(errors.New("Error while creating temporary directory: " + err.Error()))
	// }
	// TODO uncomment after fix
	// defer os.RemoveAll(s.tmpDirectory)

	err = s.generateFiles()
	if err != nil {
		return errors.New("Generation error: " + err.Error())
	}
	err = utils.CompressFiles(s.filenames, filename)
	if err != nil {
		return errors.New("Compression error: " + err.Error())
	}

	return nil
}

func (s *FileGenerator) generateFiles() error {
	s.filenames = []string{}
	errg, _ := errgroup.WithContext(context.Background())
	for i, fileData := range *s.data {
		tmpDirectory, err := os.MkdirTemp("./fixing/", fileData.Filename)
		log.Printf("Created tmp directory: %s", tmpDirectory)
		if err != nil {
			panic(errors.New("Error while creating temporary directory: " + err.Error()))
		}

		resultFilename := path.Join(tmpDirectory, (*s.data)[i].Filename+".docx")
		errg.Go(func() error {
			err := generateFile(&fileData, s.templateBytes, resultFilename, tmpDirectory)
			return err
		})
		s.filenames = append(s.filenames, resultFilename)
	}
	return errg.Wait()
}

func generateFile(fileData *fileData, templateBytes []byte, resultFilename string, tmpDirectory string) error {
	log.Printf("Enter. Data filename: %s. Result filename: %s", fileData.Filename, path.Base(resultFilename))
	var pageFilenames []string

	for i, pageData := range fileData.Pages {
		pageFilename := path.Join(tmpDirectory, fileData.Filename+"_"+fmt.Sprint(i)+".docx")
		pageFilenames = append(pageFilenames, pageFilename)

		generatePageFile(templateBytes, pageFilename, pageData)
	}

	if len(pageFilenames) >= 2 {
		err := mergePageFilesToFile(pageFilenames, resultFilename)
		if err != nil {
			return err
		}
		log.Printf("Merging. Result filename:  %s, files: %s", path.Base(resultFilename), pageFilenames)
	} else {
		err := os.Rename(pageFilenames[0], resultFilename)
		if err != nil {
			return err
		}
		log.Printf("Renamed %s to %s", path.Base(pageFilenames[0]), path.Base(resultFilename))
	}
	return nil
}

func generatePageFile(templateBytes []byte, outputFilename string, pageData docx.PlaceholderMap) error {
	template, err := docx.OpenBytes(templateBytes)
	if err != nil {
		return err
	}

	err = template.ReplaceAll(pageData)
	if err != nil {
		return err
	}

	err = template.WriteToFile(outputFilename)
	if err != nil {
		return err
	}
	return nil
}

func mergePageFilesToFile(targetFilenames []string, mergedFilename string) error {
	// args := append([]string{MERGER_PROGRAM_SET_PAGEBREAKS_OPTION, mergedFilename}, targetFilenames)
	if len(targetFilenames) < 1 {
		return errors.New("There is no specified files to merge. Pass 1 or more filenames.")
	}
	args := append([]string{mergerSetPageBreaksOption, mergedFilename}, targetFilenames...)
	mergerCommand := exec.Command(pageMergerName, args...)
	output, err := mergerCommand.Output()
	if err != nil {
		log.Println(string(output))
		panic(err)
	} else {
		return nil
	}
}

func parseJsonToData(json_data []byte) (*parseData, error) {
	parseData := new(parseData)
	err := json.Unmarshal(json_data, &parseData)
	if err != nil {
		return nil, err
	}
	return parseData, nil
}
