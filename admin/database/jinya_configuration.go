package database

type JinyaConfiguration struct {
	Group string `db:"group"`
	Key   string `db:"key"`
	Value string `db:"value"`
	Type  string `db:"type"`
}
