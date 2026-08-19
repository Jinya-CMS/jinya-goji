package http_helper

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		clientIP := strings.TrimSpace(ips[0])
		if clientIP != "" {
			return clientIP
		}
	}

	ip := r.RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	return ip
}
