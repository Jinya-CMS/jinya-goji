package database

import (
	"time"
)

type TotpMode int

const (
	TotpModeEmail TotpMode = iota
	TotpModeOtphp
)

func (t TotpMode) String() string {
	switch t {
	case TotpModeEmail:
		return "email"
	case TotpModeOtphp:
		return "app"
	default:
		return "email"
	}
}

type Artist struct {
	Id                   int        `db:"id,primarykey,autoincrement" json:"id"`
	Email                string     `db:"email" json:"email"`
	Enabled              bool       `db:"enabled" json:"enabled"`
	TwoFactorToken       *string    `db:"two_factor_token" json:"-"`
	Roles                []string   `db:"roles" json:"roles"`
	ArtistName           string     `db:"artist_name" json:"artistName"`
	ProfilePicture       *string    `db:"profile_picture" json:"profilePicture"`
	AboutMe              *string    `db:"about_me" json:"aboutMe"`
	FailedLoginAttempts  *int       `db:"failed_login_attempts" json:"-"`
	LoginBlockedUntil    *time.Time `db:"login_blocked_until" json:"-"`
	PrefersColorScheme   *bool      `db:"prefers_color_scheme" json:"-"`
	Password             string     `db:"password" json:"-"`
	TotpMode             TotpMode   `db:"totp_mode" json:"totpMode"`
	TotpSecret           *string    `db:"totp_secret" json:"-"`
	LoginMailEnabled     bool       `db:"login_mail_enabled" json:"loginMailEnabled"`
	NewDeviceMailEnabled bool       `db:"new_device_mail_enabled" json:"newDeviceMailEnabled"`
}
