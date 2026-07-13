package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"jinya.de/goji/database"
)

func GetConfigurationValue(group, key string, defaultValue any, valueType database.JinyaConfigurationType) any {
	if database.GetDbMap() != nil {
		fromDb, err := database.FindJinyaConfiguration(group, key)
		if err == nil && fromDb != nil {
			return *fromDb
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
	return GetConfigurationValue(group, key, defaultValue, database.JinyaConfigurationTypeString).(string)
}

func GetInteger(group, key string, defaultValue int) int {
	return GetConfigurationValue(group, key, defaultValue, database.JinyaConfigurationTypeInteger).(int)
}
