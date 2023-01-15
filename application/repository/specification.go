package repository

import (
	"bytes"
	"reflect"
)

const like string = " LIKE "
const ilike string = " ILIKE "
const equals string = " = "
const noEquals string = " != "
const openPercent string = "'%"
const closePercent string = "%'"
const quote string = "'"

type Specification struct {
	chains bytes.Buffer
}

func SpecBuilder() *Specification {
	s := &Specification{}
	s.chains.WriteString(WHERE)
	s.chains.WriteString(" ")
	return s
}

func (s *Specification) Like(value string, fieldName string, searchDto interface{}, isIntensive bool) *Specification {
	tagField, ok := getFieldName(fieldName, searchDto)
	if ok && len(value) > 0 {
		s.chains.WriteString(tagField)
		if isIntensive {
			s.chains.WriteString(ilike)
		} else {
			s.chains.WriteString(like)
		}
		s.aroundPercent(value)
	}
	return s
}

func (s *Specification) Equals(value string, fieldName string, searchDto interface{}) *Specification {
	tagField, ok := getFieldName(fieldName, searchDto)
	if ok && len(value) > 0 {
		s.chains.WriteString(tagField)
		s.chains.WriteString(equals)
		s.aroundQuote(value)
	}
	return s
}

func (s *Specification) NoEquals(value string, fieldName string, searchDto interface{}) *Specification {
	tagField, ok := getFieldName(fieldName, searchDto)
	if ok && len(value) > 0 {
		s.chains.WriteString(tagField)
		s.chains.WriteString(noEquals)
		s.aroundQuote(value)
	}
	return s
}

func (s *Specification) Or() *Specification {
	s.chains.WriteString(" OR ")
	return s
}

func (s *Specification) And() *Specification {
	s.chains.WriteString(" AND ")
	return s
}

func (s *Specification) Build() string {
	return s.chains.String()
}

func (s *Specification) aroundPercent(value string) {
	s.chains.WriteString(openPercent)
	s.chains.WriteString(value)
	s.chains.WriteString(closePercent)
}

func (s *Specification) aroundQuote(value string) {
	s.chains.WriteString(quote)
	s.chains.WriteString(value)
	s.chains.WriteString(quote)
}

func getFieldName(fieldName string, searchDto interface{}) (string, bool) {
	tag, ok := reflect.TypeOf(searchDto).FieldByName(fieldName)
	if ok {
		return tag.Tag.Get("pg"), ok
	}
	return "", ok
}
