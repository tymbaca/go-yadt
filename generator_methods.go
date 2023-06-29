package yadt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"golang.org/x/sync/errgroup"

	"github.com/tymbaca/go-yadt/utils"

	"github.com/lukasjarosch/go-docx"
)

var (
	leftDelimiter  = "{"
	rightDelimiter = "}"
)

// FileGenerator constructor. Takes template and json data as an io.Reader's.
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

// FileGenerator constructor. Takes template and json data in byte slices.
func NewFromBytes(templateBytes []byte, jsonBytes []byte) (*FileGenerator, error) {
	fileGenerator := new(FileGenerator)
	data, err := parseJsonToData(jsonBytes)
	if err != nil {
		return nil, err
	}

	fileGenerator.templateBytes = templateBytes
	fileGenerator.data = data

	err = fileGenerator.validateInput()
	if err != nil {
		return nil, err
	}

	return fileGenerator, nil
}

// FileGenerator constructor. Takes template and json data filenames and read them.
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

// Generates docx files and packs them into zip file at the specified path.
func (s *FileGenerator) GenerateZip(path string) error {
	var err error
	s.tmpDirectory, err = os.MkdirTemp("", "") // TODO CHANGE MkdirTemp PARAMETERS TO ("", "") AFTER FIX
	defer os.RemoveAll(s.tmpDirectory)
	if err != nil {
		panic(fmt.Errorf("error while creating temporary directory: " + err.Error()))
	}

	err = s.generateFiles()
	if err != nil {
		return err
	}
	err = utils.CompressFiles(s.filenames, path)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileGenerator) generateFiles() error {
	s.filenames = []string{}
	eg, _ := errgroup.WithContext(context.Background())

	for i, _fileData := range *s.data {
		fileData := _fileData

		resultFilename := path.Join(s.tmpDirectory, (*s.data)[i].Filename+".docx")
		eg.Go(func() error {
			err := generateFile(fileData, s.templateBytes, resultFilename, s.tmpDirectory)
			return err
		})
		s.filenames = append(s.filenames, resultFilename)
	}
	return eg.Wait()
}

func generateFile(fileData fileData, templateBytes []byte, resultFilename string, tmpDirectory string) error {
	pageFilenames, err := generatePageFiles(fileData, templateBytes, tmpDirectory)
	if err != nil {
		return err
	}

	err = mergeOrRenamePageFiles(pageFilenames, resultFilename)
	if err != nil {
		return err
	}
	return nil
}

func generatePageFiles(fileData fileData, templateBytes []byte, tmpDirectory string) ([]string, error) {
	var pageFilesPaths []string
	var eg errgroup.Group

	for i, _pageData := range fileData.Pages {
		pageData := _pageData // Needed for goroutines

		pageFilePath := path.Join(tmpDirectory, fileData.Filename+"_"+fmt.Sprint(i)+".docx")
		pageFilesPaths = append(pageFilesPaths, pageFilePath)

		eg.Go(func() error {
			err := generatePageFile(templateBytes, pageFilePath, pageData)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return pageFilesPaths, nil
}

func generatePageFile(templateBytes []byte, outputPath string, pageData docx.PlaceholderMap) error {
	template, err := docx.OpenBytes(templateBytes)
	if err != nil {
		return err
	}

	err = template.ReplaceAll(pageData)
	if err != nil {
		return err
	}

	err = template.WriteToFile(outputPath)
	if err != nil {
		return err
	}
	return nil
}

func mergeOrRenamePageFiles(pageFilesPaths []string, resultPath string) error {
	// Avoiding pagemerger call if it unnecessary
	if len(pageFilesPaths) >= 2 {
		err := mergePageFilesToFile(pageFilesPaths, resultPath)
		if err != nil {
			return err
		}
	} else if len(pageFilesPaths) == 1 {
		err := os.Rename(pageFilesPaths[0], resultPath)
		if err != nil {
			return err
		}
	} else {
		return ErrFileDataWithoutPages
	}
	return nil
}

func mergePageFilesToFile(targetFilenames []string, mergedFilename string) error {
	// args := append([]string{MERGER_PROGRAM_SET_PAGEBREAKS_OPTION, mergedFilename}, targetFilenames)
	if len(targetFilenames) < 1 {
		return errors.New("there is no specified files to merge, pass 1 or more filenames")
	}
	args := append([]string{mergerSetPageBreaksOption, mergedFilename}, targetFilenames...)
	mergerCommand := exec.Command(pageMergerName, args...)
	output, err := mergerCommand.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error while using pagemerger CLI: \n%s", string(output))
	} else {
		return nil
	}
}

func parseJsonToData(json_data []byte) (*parseData, error) {
	parseData := new(parseData)
	err := json.Unmarshal(json_data, &parseData)
	if err != nil {
		return nil, ErrBadData
	}
	return parseData, nil
}

func (g *FileGenerator) validateInput() error {
	err := g.validateTemplate()
	if err != nil {
		return err
	}

	err = g.validateData()
	if err != nil {
		return err
	}

	err = g.validateIsCompatible()
	if err != nil {
		return err
	}
	return nil
}

func (g *FileGenerator) validateTemplate() error {
	_, err := docx.OpenBytes(g.templateBytes)
	if err != nil {
		return ErrBadTemplate
	}
	_, err = utils.FindPlaceholders(g.templateBytes, leftDelimiter, rightDelimiter)
	if err != nil {
		return err
	}
	return nil
}

func (g *FileGenerator) validateData() error {
	if len(*g.data) == 0 {
		return ErrBadData
	}

	err := g.checkDataFieldsPagesExist()
	if err != nil {
		return err
	}
	err = g.checkDataFieldSameness()
	if err != nil {
		return err
	}

	return nil
}

func (g *FileGenerator) getDataFields() []string {
	fieldMap := (*g.data)[0].Pages[0]
	var dataFields []string
	for key := range fieldMap {
		dataFields = append(dataFields, key)
	}
	return dataFields
}

func (g *FileGenerator) checkDataFieldSameness() error {
	files := *g.data
	expectedFields := getPageDataFields(files[0].Pages[0])
	for _, file := range files {
		for _, page := range file.Pages {
			fields := getPageDataFields(page)
			if !utils.CompareSlicesAsSets[string](expectedFields, fields) {
				return ErrDataWithDifferentFields
			}
		}
	}
	return nil
}

func (g *FileGenerator) checkDataFieldsPagesExist() error {
	files := *g.data
	for _, file := range files {
		if len(file.Pages) == 0 {
			return ErrFileDataWithoutPages
		}
	}
	return nil
}

func (g *FileGenerator) validateIsCompatible() error {
	templateFields, err := utils.FindPlaceholders(g.templateBytes, leftDelimiter, rightDelimiter)
	if err != nil {
		return err
	}
	dataFields := getDataFields(g.data)
	if utils.CompareSlicesAsSets[string](templateFields, dataFields) {
		return nil
	} else {
		return ErrIncompatible
	}
}

func getDataFields(parseData *parseData) []string {
	data := *parseData
	fields := getPageDataFields(data[0].Pages[0])
	return fields
}

func getPageDataFields(pageData docx.PlaceholderMap) []string {
	var fields []string
	for field := range pageData {
		fields = append(fields, field)
	}
	return fields
}
