package yadt

import "errors"

var (
	ErrValidation              = errors.New("error while validating input data")
	ErrBadTemplate             = errors.New("template is bad, make sure you correctly setted placeholders")
	ErrFieldsNotMatch          = errors.New("placeholders in template doesn't match with placeholders in data")
	ErrBadData                 = errors.New("data is bad, make sure it matches with json structure descripted in go-yadt documentation")
	ErrDataWithDifferentFields = errors.New("not all pageData's in data have the same fields, make sure they are similar")
	ErrEmptyFile               = errors.New("file does not have any page data")
)
