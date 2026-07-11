package database

import "time"

type DeviceType int

const (
	DeviceTypeSmartphone DeviceType = iota
	DeviceTypeComputer
	DeviceTypeTablet
	DeviceTypeOther
)

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

type AnalyticsEntry struct {
	Id                     int         `db:"id,primarykey,autoincrement" json:"id"`
	Route                  *string     `db:"route" json:"route"`
	Timestamp              time.Time   `db:"timestamp" json:"timestamp"`
	UniqueVisit            bool        `db:"unique_visit" json:"uniqueVisit"`
	OperatingSystem        *string     `db:"operating_system" json:"operatingSystem"`
	OperatingSystemVersion *string     `db:"operating_system_version" json:"operatingSystemVersion"`
	Browser                *string     `db:"browser" json:"browser"`
	BrowserVersion         *string     `db:"browser_version" json:"browserVersion"`
	DeviceType             *DeviceType `db:"device_type" json:"deviceType"`
	EntityId               *int        `db:"entity_id" json:"entityId"`
	EntityType             *EntityType `db:"entity_type" json:"entityType"`
	Language               *string     `db:"language" json:"language"`
	Country                *string     `db:"country" json:"country"`
	Device                 *string     `db:"device" json:"device"`
	Brand                  *string     `db:"brand" json:"brand"`
	UserAgent              *string     `db:"user_agent" json:"userAgent"`
	Referer                *string     `db:"referer" json:"-"`
	Status                 int         `db:"status" json:"-"`
}
