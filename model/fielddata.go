package model

import "strings"

type fieldNames []string

func NewFieldName(p *[]string) *fieldNames {
	return (*fieldNames)(p)
}

func (f *fieldNames) String() string {
	*f = fieldNames([]string{"Score1", "Score2", "Score3", "Score4", "Score5"})
	return ""
}

func (f *fieldNames) Set(val string) error {
	*f = fieldNames(strings.Split(val, ":"))
	return nil
}
