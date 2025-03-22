package util

import (
	"fmt"
	"strings"
)

type SqlBuilder struct {
	query strings.Builder
	args  []interface{}
	index int
}

type SqlBuilderResult struct {
	Query string
	Args  []interface{}
}

func NewSqlBuilder(clauses ...string) *SqlBuilder {
	sb := &SqlBuilder{
		index: 1,
	}

	for _, clause := range clauses {
		sb.writeClause(clause)
	}
	return sb
}

func (s *SqlBuilder) Concat(clause string, values ...interface{}) *SqlBuilder {
	if len(values) == 0 {
		s.query.WriteString(fmt.Sprintf("\n%s", clause))
		return s
	}

	replacementIndices := make([]interface{}, len(values))
	startIndex := s.index

	for i := range values {
		replacementIndices[i] = startIndex + i
	}

	formattedClause := fmt.Sprintf(clause, replacementIndices...)

	s.writeClause(formattedClause)

	s.args = append(s.args, values...)
	s.index += len(values)

	return s
}

func (s *SqlBuilder) Result() SqlBuilderResult {
	return SqlBuilderResult{
		Query: strings.TrimSpace(s.query.String()),
		Args:  s.args,
	}
}

// Deprecated: Use Result instead.
func (s *SqlBuilder) GetArgs() []interface{} {
	return s.args
}

// Deprecated: Use Result instead.
func (s *SqlBuilder) Build() string {
	return strings.TrimSpace(s.query.String())
}

func (s *SqlBuilder) writeClause(clause string) {
	s.query.WriteString(fmt.Sprintf("\n%s", clause))
}
