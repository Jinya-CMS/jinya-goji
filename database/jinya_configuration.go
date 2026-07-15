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

func FindJinyaConfiguration(group, key string) (*JinyaConfiguration, error) {
	return SelectOne[JinyaConfiguration](`
select "group", key, value, type
from jinya_configuration
where "group" = $1
  and key = $2
`, group, key)
}

func SetJinyaConfiguration(key, group string, value string, configType JinyaConfigurationType) error {
	_, err := dbMap.Exec(`
insert into jinya_configuration ("group", key, value, type)
values ($1, $2, $3, $4)
on conflict ("group", key) do update set value = $3`, group, key, value, configType)
	return err
}

func RemoveJinyaConfiguration(key, group string) error {
	_, err := dbMap.Exec(`
delete
from jinya_configuration
where "group" = $1
  and key = $2`, group, key)
	return err
}
