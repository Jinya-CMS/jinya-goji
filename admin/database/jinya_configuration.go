package database

type JinyaConfiguration struct {
	Group string `db:"group"`
	Key   string `db:"key"`
	Value string `db:"value"`
	Type  string `db:"type"`
}

func FindJinyaConfiguration(key string, group string) (*JinyaConfiguration, error) {
	return SelectOne[JinyaConfiguration](`
SELECT "group", "key", value, type
FROM jinya_configuration
WHERE "group" = ? AND "key" = ?
`, group, key)
}
