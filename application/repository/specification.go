package repository

import (
	"bytes"
	"fmt"
	"reflect"
)

const like string = "LIKE"

// TODO перейти на байты со стринги
type Specification struct {
	chains []string
}

func CreateSpec() *Specification {
	return &Specification{}
}

func (s *Specification) Like(value string, fieldName string, searchDto interface{}) *Specification {
	tag, ok := reflect.TypeOf(searchDto).FieldByName(fieldName)
	if ok && len(value) > 0 {
		s.chains = append(s.chains, fmt.Sprint(tag.Tag.Get("pg"), " ", like, " '%", value, "%'"))
	}
	return s
}

func (s *Specification) Or() *Specification {
	s.chains = append(s.chains, " OR ")
	return s
}

func (s *Specification) And() *Specification {
	s.chains = append(s.chains, " AND ")
	return s
}

func (s *Specification) Build() string {
	buffer := bytes.Buffer{}
	for i, chain := range s.chains {
		if i == 0 {
			buffer.WriteString(WHERE)
			buffer.WriteString(" ")
			buffer.WriteString(chain)
		} else {
			buffer.WriteString(fmt.Sprint(" ", chain))
		}
	}
	return buffer.String()
}
