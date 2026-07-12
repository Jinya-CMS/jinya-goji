package analyzer

import (
	"strings"
)

type Statement struct {
	Query string
}

type QueryAnalyzer struct{}

func (qa *QueryAnalyzer) GetQueryType(statement Statement) string {
	query := strings.TrimSpace(statement.Query)
	parts := strings.Split(query, " ")
	if len(parts) == 0 {
		return ""
	}

	return strings.ToUpper(parts[0])
}

func (qa *QueryAnalyzer) GetStatements(query string) []Statement {
	rawStatements := strings.Split(query, ";")
	var statements []Statement
	for _, raw := range rawStatements {
		trimmed := strings.TrimSpace(raw)
		if trimmed != "" {
			statements = append(statements, Statement{Query: trimmed})
		}
	}
	return statements
}
