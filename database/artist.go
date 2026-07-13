package database

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/DerKnerd/gorp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type TotpMode int

const (
	TotpModeEmail TotpMode = iota
	TotpModeApp
)

func (t TotpMode) String() string {
	switch t {
	case TotpModeEmail:
		return "email"
	case TotpModeApp:
		return "app"
	default:
		return "email"
	}
}

type Artist struct {
	Id                   int        `db:"id,primarykey,autoincrement"`
	Email                string     `db:"email,unique"`
	Enabled              bool       `db:"enabled"`
	TwoFactorToken       *string    `db:"two_factor_token"`
	RolesPlain           string     `db:"roles"`
	Roles                []string   `db:"-"`
	ArtistName           string     `db:"artist_name"`
	ProfilePicture       *string    `db:"profile_picture"`
	AboutMe              *string    `db:"about_me"`
	FailedLoginAttempts  *int       `db:"failed_login_attempts"`
	LoginBlockedUntil    *time.Time `db:"login_blocked_until"`
	PrefersColorScheme   *bool      `db:"prefers_color_scheme"`
	Password             string     `db:"password"`
	TotpMode             TotpMode   `db:"totp_mode"`
	TotpSecret           *string    `db:"totp_secret"`
	LoginMailEnabled     bool       `db:"login_mail_enabled"`
	NewDeviceMailEnabled bool       `db:"new_device_mail_enabled"`
}

func (artist *Artist) PostGet(s gorp.SqlExecutor) error {
	return json.Unmarshal([]byte(artist.RolesPlain), &artist.Roles)
}

func (artist *Artist) PreInsert(s gorp.SqlExecutor) error {
	roles, err := json.Marshal(artist.Roles)
	if err != nil {
		return err
	}

	artist.RolesPlain = string(roles)

	return nil
}

func (artist *Artist) PreUpdate(s gorp.SqlExecutor) error {
	roles, err := json.Marshal(artist.Roles)
	if err != nil {
		return err
	}

	artist.RolesPlain = string(roles)

	return nil
}

func (artist *Artist) MarshalJSON() ([]byte, error) {
	colorScheme := "auto"
	if artist.PrefersColorScheme != nil {
		if *artist.PrefersColorScheme {
			colorScheme = "dark"
		} else {
			colorScheme = "light"
		}
	}

	return json.Marshal(struct {
		Id                   int      `json:"id"`
		Email                string   `json:"email"`
		Enabled              bool     `json:"enabled"`
		Roles                []string `json:"roles"`
		ArtistName           string   `json:"artistName"`
		ProfilePicture       *string  `json:"profilePicture"`
		AboutMe              *string  `json:"aboutMe"`
		TotpMode             string   `json:"totpMode"`
		ColorScheme          string   `json:"colorScheme"`
		LoginMailEnabled     bool     `json:"loginMailEnabled"`
		NewDeviceMailEnabled bool     `json:"newDeviceMailEnabled"`
	}{
		Id:                   artist.Id,
		Email:                artist.Email,
		Enabled:              artist.Enabled,
		Roles:                artist.Roles,
		ArtistName:           artist.ArtistName,
		ProfilePicture:       artist.ProfilePicture,
		AboutMe:              artist.AboutMe,
		TotpMode:             artist.TotpMode.String(),
		ColorScheme:          colorScheme,
		LoginMailEnabled:     artist.LoginMailEnabled,
		NewDeviceMailEnabled: artist.NewDeviceMailEnabled,
	})
}

func (artist *Artist) Create() error {
	return Insert(artist)
}

func (artist *Artist) Update() error {
	return Update(artist)
}

func (artist *Artist) Delete() error {
	return Delete(artist)
}

func FindArtistById(id int) (*Artist, error) {
	return Get[Artist](id)
}

func FindArtistByEmail(email string) (*Artist, error) {
	return SelectOne[Artist](`
select id,
       email,
       enabled,
       two_factor_token,
       roles,
       artist_name,
       profile_picture,
       about_me,
       failed_login_attempts,
       login_blocked_until,
       prefers_color_scheme,
       password,
       totp_mode,
       totp_secret,
       login_mail_enabled,
       new_device_mail_enabled
from users
where email = $1
`, email)
}

func FindAllArtists() ([]Artist, error) {
	return Select[Artist](`
select id,
       email,
       enabled,
       two_factor_token,
       roles,
       artist_name,
       profile_picture,
       about_me,
       failed_login_attempts,
       login_blocked_until,
       prefers_color_scheme,
       password,
       totp_mode,
       totp_secret,
       login_mail_enabled,
       new_device_mail_enabled
from users
`)
}

func (artist *Artist) SetTwoFactorCode() error {
	token := ""
	for i := 0; i < 6; i++ {
		val, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return err
		}
		token += fmt.Sprintf("%d", val.Int64())
	}
	artist.TwoFactorToken = &token
	return nil
}

func (artist *Artist) ValidateDevice(knownDeviceCode string) bool {
	count, _ := dbMap.SelectInt(`
select count(*)
from users u
         join known_device kd on kd.user_id = u.id
where u.id = $1`, artist.Id)

	return count > 0
}

func (artist *Artist) ChangePassword(oldPassword, password string) error {
	if artist.ValidatePassword(oldPassword) {
		err := artist.SetPassword(password)
		if err != nil {
			return err
		}
		err = artist.Update()
		if err != nil {
			return err
		}

		_, _ = dbMap.Exec(`
delete
from api_key
where user_id = $1`, artist.Id)

		return nil
	}

	return nil
}

func (artist *Artist) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(artist.Password), []byte(password))
	return err == nil && artist.Enabled
}

func (artist *Artist) SetPassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	artist.Password = string(hash)
	return nil
}

func CountAdmins() (int, error) {
	artists, err := FindAllArtists()
	if err != nil {
		return 0, err
	}

	admins := 0
	for _, artist := range artists {
		if artist.Enabled && slices.Contains(artist.Roles, "ROLE_ADMIN") {
			admins++
		}
	}

	return admins, nil
}

func (artist *Artist) RegisterFailedLogin() error {
	attempts := 0
	if artist.FailedLoginAttempts != nil {
		attempts = *artist.FailedLoginAttempts
	}
	attempts++
	artist.FailedLoginAttempts = &attempts
	if attempts >= 5 {
		artist.LoginBlockedUntil = new(time.Now().Add(10 * time.Minute))
	}
	return artist.Update()
}

func (artist *Artist) UnlockAccount() error {
	artist.FailedLoginAttempts = nil
	artist.LoginBlockedUntil = nil
	return artist.Update()
}

func (artist *Artist) VerifyTotpCode(code string) bool {
	if artist.TotpSecret != nil && artist.TotpMode == TotpModeApp {
		return totp.Validate(code, *artist.TotpSecret)
	}

	return artist.TwoFactorToken != nil && *artist.TwoFactorToken == code
}
