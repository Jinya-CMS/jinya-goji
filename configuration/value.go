package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"jinya.de/goji/database"
)

func getConfigurationValue(group, key string, defaultValue any, valueType database.JinyaConfigurationType) any {
	if database.GetDbMap() != nil {
		fromDb, err := database.FindJinyaConfiguration(group, key)
		if err == nil && fromDb != nil {
			return fromDb.Value
		}
	}

	fromEnv := os.Getenv(strings.ToUpper(fmt.Sprintf("%s_%s", group, key)))
	if fromEnv != "" {
		switch valueType {
		case database.JinyaConfigurationTypeString:
			return fromEnv
		case database.JinyaConfigurationTypeInteger:
			res, err := strconv.Atoi(fromEnv)
			if err == nil {
				return res
			}
		case database.JinyaConfigurationTypeBoolean:
			return fromEnv == "true"
		}
	}

	return defaultValue
}

func GetString(group, key, defaultValue string) string {
	return getConfigurationValue(group, key, defaultValue, database.JinyaConfigurationTypeString).(string)
}

func GetInteger(group, key string, defaultValue int) int {
	val := getConfigurationValue(group, key, defaultValue, database.JinyaConfigurationTypeInteger)
	if s, ok := val.(string); ok {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}

	if i, ok := val.(int); ok {
		return i
	}

	return defaultValue
}

func SetString(group, key string, value *string) error {
	if value == nil {
		return database.RemoveJinyaConfiguration(key, group)
	}

	return database.SetJinyaConfiguration(key, group, *value, database.JinyaConfigurationTypeString)
}

func SetInteger(group, key string, value int) error {
	return database.SetJinyaConfiguration(key, group, strconv.Itoa(value), database.JinyaConfigurationTypeInteger)
}
