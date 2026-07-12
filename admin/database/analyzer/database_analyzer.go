package analyzer

import (
	"fmt"

	"jinya.de/goji/database"
)

func GetTables() (map[string]any, error) {
	db := database.GetDbMap()
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]any)
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}

		tableData := make(map[string]any)

		// Structure
		structure, err := executeSqlString(fmt.Sprintf("EXPLAIN %s", tableName))
		if err == nil {
			tableData["structure"] = structure
		}

		// Entry count
		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
		if err == nil {
			tableData["entryCount"] = count
		}

		// Size and Engine
		var bytes int64
		var engine string
		err = db.QueryRow(`
			SELECT ROUND(DATA_LENGTH + INDEX_LENGTH) AS bytes, ENGINE
			FROM information_schema.TABLES
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		`, tableName).Scan(&bytes, &engine)
		if err == nil {
			tableData["size"] = bytes
			tableData["engine"] = engine
		}

		// Constraints
		constraints, err := getTableConstraints(tableName)
		if err == nil {
			tableData["constraints"] = constraints
		}

		// Indexes
		indexes, err := executeSqlString(fmt.Sprintf("SHOW INDEXES FROM %s", tableName))
		if err == nil {
			tableData["indexes"] = indexes
		}

		result[tableName] = tableData
	}

	return result, nil
}

func getTableConstraints(tableName string) ([]map[string]any, error) {
	query := `
		SELECT
			kcu.CONSTRAINT_NAME AS CONSTRAINT_NAME,
			kcu.REFERENCED_TABLE_NAME AS REFERENCED_TABLE_NAME,
			kcu.REFERENCED_COLUMN_NAME AS REFERENCED_COLUMN_NAME,
			kcu.POSITION_IN_UNIQUE_CONSTRAINT AS POSITION_IN_UNIQUE_CONSTRAINT,
			tc.CONSTRAINT_TYPE AS CONSTRAINT_TYPE,
			rc.DELETE_RULE AS DELETE_RULE,
			rc.UPDATE_RULE AS UPDATE_RULE,
			kcu.COLUMN_NAME AS COLUMN_NAME
		FROM information_schema.TABLE_CONSTRAINTS tc
		INNER JOIN information_schema.KEY_COLUMN_USAGE kcu
			ON kcu.CONSTRAINT_NAME = tc.CONSTRAINT_NAME AND tc.TABLE_NAME = kcu.TABLE_NAME
		LEFT JOIN information_schema.REFERENTIAL_CONSTRAINTS rc
			ON rc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
		WHERE tc.TABLE_NAME = ? AND tc.TABLE_SCHEMA = DATABASE()
	`
	return executeSqlString(query, tableName)
}

func GetServerInfo() (map[string]string, error) {
	variables, err := GetVariables(VariablesTypeGlobal)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"version":        fmt.Sprint(variables["version"]),
		"comment":        fmt.Sprint(variables["version_comment"]),
		"compileMachine": fmt.Sprint(variables["version_compile_machine"]),
		"compileOs":      fmt.Sprint(variables["version_compile_os"]),
	}, nil
}

func GetVariables(variableType VariablesType) (map[string]any, error) {
	query := fmt.Sprintf("SHOW %s VARIABLES", variableType.String())
	rows, err := executeSqlString(query)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	for _, row := range rows {
		if name, ok := row["Variable_name"].(string); ok {
			result[name] = row["Value"]
		}
	}

	return result, nil
}

func executeSqlString(query string, args ...any) ([]map[string]any, error) {
	db := database.GetDbMap()
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		result = append(result, row)
	}

	return result, nil
}
