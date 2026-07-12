package mailing

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"jinya.de/goji/analytics"
	"jinya.de/goji/localization"
	"jinya.de/goji/locate"

	"github.com/wneessen/go-mail"
)

func sendEmailByRequest(r *http.Request, plainBody, htmlBody, subject, to string) error {
	slog.Debug("Send email", "to", to, "subject", subject)
	message := mail.NewMsg()
	_ = message.To(to)
	_ = message.FromFormat(localization.Localize(r, "mailer-from", os.Getenv("APP_NAME")), os.Getenv("MAILER_FROM"))

	return sendEmail(message, subject, plainBody, htmlBody)
}

func sendEmail(message *mail.Msg, subject string, plainBody string, htmlBody string) error {
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, plainBody)
	message.SetBodyString(mail.TypeTextHTML, htmlBody)

	senderHost := os.Getenv("MAILER_SERVER")
	senderUser := os.Getenv("MAILER_USER")
	senderPass := os.Getenv("MAILER_PASSWORD")
	useSsl := os.Getenv("MAILER_ENCRYPTION") == "ssl"
	useStartTls := os.Getenv("MAILER_ENCRYPTION") == "tls"
	port, err := strconv.Atoi(os.Getenv("MAILER_PORT"))
	if err != nil {
		if useSsl {
			port = mail.DefaultPortSSL
		} else if useStartTls {
			port = mail.DefaultPortTLS
		} else {
			port = mail.DefaultPort
		}
	}
	options := []mail.Option{
		mail.WithPort(port),
	}
	if senderUser != "" && senderPass != "" {
		options = append(options, mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(senderUser), mail.WithPassword(senderPass))
	}
	if useSsl {
		options = append(options, mail.WithSSL())
	}
	if useStartTls {
		options = append(options, mail.WithTLSPolicy(mail.DefaultTLSPolicy))
	}
	if !useSsl && !useStartTls {
		options = append(options, mail.WithTLSPolicy(mail.NoTLS))
	}
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		options = append(options, mail.WithDebugLog())
	}

	slog.Debug("Creating new client to smtp server", "host", senderHost, "port", port, "username", senderUser, "useSsl", useSsl, "useStartTls", useStartTls)
	client, err := mail.NewClient(senderHost, options...)
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	defer client.Close()

	slog.Debug("Sending email")
	err = client.DialAndSend(message)
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	return nil
}

func SendNewLoginMail(r *http.Request, artistName, artistEmail string) {
	ua := analytics.UserAgent.Parse(r.UserAgent())
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	location, _ := locate.LocateIp(ip)

	data := struct {
		ArtistName    string
		RemoteAddress string
		Browser       string
		Subject       string
		Location      locate.Location
	}{
		ArtistName:    artistName,
		RemoteAddress: r.RemoteAddr,
		Browser:       ua.Browser().String(),
		Subject:       localization.Localize(r, "new_login_subject"),
		Location:      location,
	}

	mjml, err := renderTemplate(r, "new-login", data)
	if err != nil {
		slog.Error("Failed to execute new login mjml body template", "error", err)
		return
	}

	plain, err := renderPlain(r, "new-login", data)
	if err != nil {
		slog.Error("Failed to execute new login plain body template", "error", err)
		return
	}

	_ = sendEmailByRequest(r, plain, renderMrmlTemplate(mjml), data.Subject, artistEmail)
}

func SendNewDeviceMail(r *http.Request, artistName, artistEmail string) {
	ua := analytics.UserAgent.Parse(r.UserAgent())
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	location, _ := locate.LocateIp(ip)

	data := struct {
		ArtistName    string
		RemoteAddress string
		Browser       string
		Subject       string
		Location      locate.Location
	}{
		ArtistName:    artistName,
		RemoteAddress: r.RemoteAddr,
		Browser:       ua.Browser().String(),
		Subject:       localization.Localize(r, "new_device_subject"),
		Location:      location,
	}

	mjml, err := renderTemplate(r, "new-saved-device", data)
	if err != nil {
		slog.Error("Failed to execute new saved device mjml body template", "error", err)
		return
	}

	plain, err := renderPlain(r, "new-saved-device", data)
	if err != nil {
		slog.Error("Failed to execute new saved device plain body template", "error", err)
		return
	}

	_ = sendEmailByRequest(r, plain, renderMrmlTemplate(mjml), data.Subject, artistEmail)
}

func SendTwoFactorCodeMail(r *http.Request, artistName, artistEmail, twoFactorCode string) {
	data := struct {
		ArtistName    string
		Subject       string
		TwoFactorCode string
	}{
		ArtistName:    artistName,
		Subject:       localization.Localize(r, "two_factor_code_subject"),
		TwoFactorCode: twoFactorCode,
	}

	mjml, err := renderTemplate(r, "two-factor-code", data)
	if err != nil {
		slog.Error("Failed to execute two factor code mjml body template", "error", err)
		return
	}

	plain, err := renderPlain(r, "two-factor-code", data)
	if err != nil {
		slog.Error("Failed to execute two factor code plain body template", "error", err)
		return
	}

	_ = sendEmailByRequest(r, plain, renderMrmlTemplate(mjml), data.Subject, artistEmail)
}
