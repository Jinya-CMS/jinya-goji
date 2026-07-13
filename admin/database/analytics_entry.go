package database

import (
	"encoding/json"
	"strconv"
	"time"
)

type DeviceType int

const (
	DeviceTypeSmartphone DeviceType = iota
	DeviceTypeComputer
	DeviceTypeTablet
	DeviceTypeOther
)

func (d DeviceType) String() string {
	switch d {
	case DeviceTypeSmartphone:
		return "smartphone"
	case DeviceTypeComputer:
		return "computer"
	case DeviceTypeTablet:
		return "tablet"
	case DeviceTypeOther:
		return "other"
	default:
		return "other"
	}
}

type EntityType int

const (
	EntityTypeBlogPost EntityType = iota
	EntityTypeBlogCategory
	EntityTypeForm
	EntityTypeGallery
	EntityTypeModernPage
	EntityTypeClassicPage
	EntityTypeArtist
)

func (e EntityType) String() string {
	switch e {
	case EntityTypeBlogPost:
		return "blog-post"
	case EntityTypeBlogCategory:
		return "blog-category"
	case EntityTypeForm:
		return "form"
	case EntityTypeGallery:
		return "gallery"
	case EntityTypeModernPage:
		return "modern-page"
	case EntityTypeClassicPage:
		return "classic-page"
	case EntityTypeArtist:
		return "artist"
	default:
		return ""
	}
}

type AnalyticsEntry struct {
	Id                     int         `db:"id,primarykey,autoincrement"`
	Route                  *string     `db:"route"`
	Timestamp              time.Time   `db:"timestamp"`
	UniqueVisit            bool        `db:"unique_visit"`
	OperatingSystem        *string     `db:"operating_system"`
	OperatingSystemVersion *string     `db:"operating_system_version"`
	Browser                *string     `db:"browser"`
	BrowserVersion         *string     `db:"browser_version"`
	DeviceType             *DeviceType `db:"device_type"`
	EntityId               *int        `db:"entity_id"`
	EntityType             *EntityType `db:"entity_type"`
	Language               *string     `db:"language"`
	Country                *string     `db:"country"`
	Device                 *string     `db:"device"`
	Brand                  *string     `db:"brand"`
	UserAgent              *string     `db:"user_agent"`
	Referer                *string     `db:"referer"`
	Status                 int         `db:"status"`
}

func (entry *AnalyticsEntry) MarshalJSON() ([]byte, error) {
	var deviceType *string
	if entry.DeviceType != nil {
		deviceType = new(entry.DeviceType.String())
	}
	var entityType *string
	if entry.EntityType != nil {
		entityType = new(entry.EntityType.String())
	}

	return json.Marshal(struct {
		Id                     int     `json:"id"`
		Route                  *string `json:"route"`
		Timestamp              string  `json:"timestamp"`
		UniqueVisit            bool    `json:"uniqueVisit"`
		OperatingSystem        *string `json:"operatingSystem"`
		OperatingSystemVersion *string `json:"operatingSystemVersion"`
		Browser                *string `json:"browser"`
		BrowserVersion         *string `json:"browserVersion"`
		DeviceType             *string `json:"deviceType"`
		EntityId               *int    `json:"entityId"`
		EntityType             *string `json:"entityType"`
		Language               *string `json:"language"`
		Country                *string `json:"country"`
		Device                 *string `json:"device"`
		Brand                  *string `json:"brand"`
		UserAgent              *string `json:"userAgent"`
	}{
		Id:                     entry.Id,
		Route:                  entry.Route,
		Timestamp:              entry.Timestamp.Format(time.RFC3339),
		UniqueVisit:            entry.UniqueVisit,
		OperatingSystem:        entry.OperatingSystem,
		OperatingSystemVersion: entry.OperatingSystemVersion,
		Browser:                entry.Browser,
		BrowserVersion:         entry.BrowserVersion,
		DeviceType:             deviceType,
		EntityId:               entry.EntityId,
		EntityType:             entityType,
		Language:               entry.Language,
		Country:                entry.Country,
		Device:                 entry.Device,
		Brand:                  entry.Brand,
		UserAgent:              entry.UserAgent,
	})
}

func FindAnalyticsById(id int) (*AnalyticsEntry, error) {
	return Get[AnalyticsEntry](id)
}

func FindAllAnalytics() ([]AnalyticsEntry, error) {
	return Select[AnalyticsEntry](`
select id,
       route,
       timestamp,
       unique_visit,
       operating_system,
       operating_system_version,
       browser,
       browser_version,
       device_type,
       entity_id,
       entity_type,
       language,
       country,
       device,
       brand,
       user_agent,
       referer,
       status
from analytics
`)
}

type AnalyticsPastIntervalResult struct {
	Visits int    `db:"visits"`
	Group  string `db:"group"`
}

func GetAnalyticsPastInterval(months *int, entityType *EntityType, id *int, uniqueOnly bool) ([]AnalyticsPastIntervalResult, error) {
	query := `
select count(*) as visits,
       timestamp as 'group'
from analytics
where status <> 404
`
	args := make([]any, 0)
	if months != nil {
		query += " and timestamp >= subdate(current_date, interval " + strconv.Itoa(*months) + " month)"
	}
	if uniqueOnly {
		query += " and unique_visit = true"
	}
	if entityType != nil && id != nil {
		query += " and entity_type = $1 and entity_id = $2"
		args = append(args, *entityType, *id)
	}
	query += " group by timestamp"

	return Select[AnalyticsPastIntervalResult](query, args...)
}

func GetAnalyticsPastIntervalGrouped(group string, months *int, entityType *EntityType, id *int, uniqueOnly bool) ([]AnalyticsPastIntervalResult, error) {
	groupColumns := make([]string, 0)
	switch group {
	case "os":
		groupColumns = append(groupColumns, "operating_system")
	case "os-version":
		groupColumns = append(groupColumns, "operating_system", "operating_system_version")
	case "browser-version":
		groupColumns = append(groupColumns, "browser", "browser_version")
	case "type":
		groupColumns = append(groupColumns, "device_type")
	case "date":
		groupColumns = append(groupColumns, "timestamp")
	default:
		groupColumns = append(groupColumns, group)
	}

	var groupColumnsForName string
	if len(groupColumns) > 1 {
		groupColumnsForName = "concat_ws(' ', " + groupColumns[0] + ", " + groupColumns[1] + ")"
	} else {
		groupColumnsForName = groupColumns[0]
	}

	query := `
select count(*) as visits, ` + groupColumnsForName + ` as 'group'
from analytics
where status <> 404`
	args := make([]any, 0)

	if months != nil {
		query += " and timestamp >= subdate(current_date, interval " + strconv.Itoa(*months) + " month)"
	}
	if uniqueOnly {
		query += " and unique_visit = true"
	}
	if entityType != nil && id != nil {
		query += " and entity_type = ? and entity_id = ?"
		args = append(args, *entityType, *id)
	}

	query += " group by "
	for i, col := range groupColumns {
		if i > 0 {
			query += ", "
		}
		query += col
	}

	return Select[AnalyticsPastIntervalResult](query, args...)
}

func GetTotalAnalyticsPastInterval(months *int) (int, error) {
	query := `
select count(*) as visits
from analytics
where status <> 404
  and unique_visit = true
`
	args := make([]any, 0)
	if months != nil {
		query += " and timestamp >= subdate(current_date, interval " + strconv.Itoa(*months) + " month)"
	}

	var result struct {
		Visits int `db:"visits"`
	}
	err := dbMap.SelectOne(&result, query, args...)
	return result.Visits, err
}

func GetAnalyticsPastNMonthsForEntity(months int, entityType EntityType, id int) ([]AnalyticsPastIntervalResult, error) {
	return GetAnalyticsPastInterval(&months, &entityType, &id, false)
}

func GetAnalyticsPastNMonthsForEntityGrouped(group string, months int, entityType EntityType, id int) ([]AnalyticsPastIntervalResult, error) {
	return GetAnalyticsPastIntervalGrouped(group, &months, &entityType, &id, false)
}

func (entry *AnalyticsEntry) GetRelatedMenuItem() (*MenuItem, error) {
	if entry.Route == nil {
		return nil, nil
	}
	return FindMenuItemByRoute(*entry.Route)
}
