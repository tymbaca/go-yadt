package yadt

import (
	"errors"

	"github.com/tymbaca/go-yadt/utils"
)

var (
	ErrPlaceholderWithWhitespaces   = utils.ErrPlaceholderWithWhitespaces
	ErrTemplatePlaceholdersNotFound = utils.ErrTemplatePlaceholdersNotFound
	ErrBadTemplate                  = utils.ErrBadTemplate
	ErrValidation                   = errors.New("error while validating input data")
	ErrIncompatible                 = errors.New("placeholders in template doesn't match with placeholders in data")
	ErrBadData                      = errors.New("data is bad, make sure it matches with json structure descripted in go-yadt documentation")
	ErrDataWithDifferentFields      = errors.New("not all pageData's in data have the same fields, make sure they are similar")
	ErrFileDataWithoutPages         = errors.New("file does not have any page data")
)
