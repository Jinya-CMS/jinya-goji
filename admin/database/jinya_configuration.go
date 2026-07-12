package database

type JinyaConfigurationType int

const (
	JinyaConfigurationTypeString = iota
	JinyaConfigurationTypeInteger
	JinyaConfigurationTypeBoolean
)

type JinyaConfiguration struct {
	Group string                 `db:"group"`
	Key   string                 `db:"key"`
	Value string                 `db:"value"`
	Type  JinyaConfigurationType `db:"type"`
}

func FindJinyaConfiguration(key string, group string) (*JinyaConfiguration, error) {
	return SelectOne[JinyaConfiguration](`
select "group", key, value, type
from jinya_configuration
where "group" = $1
  and key = $2
`, group, key)
}
