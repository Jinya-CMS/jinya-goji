package analytics

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/medama-io/go-useragent/agents"
	"jinya.de/goji/auth"
	"jinya.de/goji/database"
	"jinya.de/goji/http_helper"
	"jinya.de/goji/locate"
)

func TrackStandardRequest(r *http.Request) {
	entry, track, err := getBaseEntry(r)
	if !track {
		return
	}

	var entityId *int
	var entityType *database.EntityType
	route, err := database.FindMenuItemByRoute(strings.TrimPrefix(r.URL.Path, "/"))
	if err == nil {
		if route.ModernPageId != nil {
			entityId = route.ModernPageId
			entityType = new(database.EntityTypeModernPage)
		} else if route.ClassicPageId != nil {
			entityId = route.ClassicPageId
			entityType = new(database.EntityTypeClassicPage)
		} else if route.GalleryId != nil {
			entityId = route.GalleryId
			entityType = new(database.EntityTypeGallery)
		} else if route.CategoryId != nil {
			entityId = route.CategoryId
			entityType = new(database.EntityTypeBlogCategory)
		} else if route.ArtistId != nil {
			entityId = route.ArtistId
			entityType = new(database.EntityTypeArtist)
		}
	}

	entry.EntityId = entityId
	entry.EntityType = entityType

	err = entry.Create()
	if err != nil {
		slog.Error("Failed to track request", "error", err)
	}
}

func TrackBlogPostRequest(r *http.Request) {
	entry, track, err := getBaseEntry(r)
	if !track {
		return
	}

	vars := mux.Vars(r)
	slug := vars["slug"]

	blogPost, err := database.FindBlogPostBySlug(slug)
	if err != nil {
		return
	}

	entry.EntityId = &blogPost.Id
	entry.EntityType = new(database.EntityTypeBlogPost)

	err = entry.Create()
	if err != nil {
		slog.Error("Failed to track request", "error", err)
	}
}

func getBaseEntry(r *http.Request) (*database.AnalyticsEntry, bool, error) {
	parsedAgent := UserAgent.Parse(r.UserAgent())
	if parsedAgent.IsBot() || r.UserAgent() == "" {
		slog.Info("The request is made by a robot, ignore it")
		return nil, false, nil
	}
	if r.Header.Get("Accept-Language") == "" {
		slog.Info("The request is made by a browser without language, ignore it")
		return nil, false, nil
	}
	if len(r.CookiesNamed(auth.AuthenticationCookieName)) > 0 || len(r.CookiesNamed(auth.KnownDeviceCookieName)) > 0 {
		slog.Info("The request is made by an artist, ignore it")
		return nil, false, nil
	}

	ip := http_helper.GetClientIP(r)

	location, err := locate.LocalizeIp(ip)
	if err != nil {
		slog.Info("Failed to localize ip address", "error", err)
		location.Country = "ZZ"
	}

	refererHost := ""
	if r.Referer() != "" {
		referer, err := url.Parse(r.Referer())
		if err != nil {
			slog.Info("The referer is invalid", "error", err)
		} else {
			refererHost = referer.Host
		}
	}

	var deviceType database.DeviceType
	switch parsedAgent.Device() {
	case agents.DeviceDesktop:
		deviceType = database.DeviceTypeComputer
		break
	case agents.DeviceMobile:
		deviceType = database.DeviceTypeSmartphone
		break
	case agents.DeviceTablet:
		deviceType = database.DeviceTypeTablet
		break
	default:
		deviceType = database.DeviceTypeOther
		break
	}

	language := r.Header.Get("Accept-Language")
	if len(language) > 2 {
		language = language[:2]
	}

	entry := database.AnalyticsEntry{
		Route:           new(strings.TrimPrefix(r.URL.Path, "/")),
		Timestamp:       time.Now(),
		UniqueVisit:     refererHost == r.Host,
		OperatingSystem: new(parsedAgent.OS().String()),
		Browser:         new(parsedAgent.Browser().String()),
		BrowserVersion:  new(parsedAgent.BrowserVersion()),
		DeviceType:      &deviceType,
		Language:        &language,
		Country:         &location.Country,
		UserAgent:       new(r.UserAgent()),
		Referer:         &refererHost,
	}

	return &entry, true, err
}
