// Package domain provides domain monitoring and WHOIS lookup functionality.
package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/integration"
	"gabe565.com/domain-watch/internal/message"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/r3labs/diff/v3"
)

// New creates a new Domain instance for the given domain name.
func New(conf *config.Config, name string) Domain {
	return Domain{conf: conf, Name: name}
}

// Domain represents a domain being monitored for changes and expiration.
type Domain struct {
	conf *config.Config

	Name               string
	CurrWhois          whoisparser.WhoisInfo
	PrevWhois          *whoisparser.WhoisInfo
	ExpiresAt          time.Time
	TimeLeft           time.Duration
	TriggeredThreshold int
}

// Whois fetches and parses WHOIS information for the domain.
func (d Domain) Whois() (whoisparser.WhoisInfo, error) {
	raw, err := whois.Whois(d.Name)
	if err != nil && raw == "" {
		return whoisparser.WhoisInfo{}, err
	}

	return whoisparser.Parse(raw)
}

// Log returns a logger with the domain name attached.
func (d Domain) Log() *slog.Logger {
	return slog.With("domain", d.Name)
}

// ErrNoExpiration is returned when the expiration date cannot be found.
var ErrNoExpiration = errors.New("could not find expiration date")

// Run fetches current WHOIS data and checks for notifications.
func (d *Domain) Run(ctx context.Context, integrations integration.Integrations) error {
	var err error
	d.CurrWhois, err = d.Whois()
	if err != nil {
		return fmt.Errorf("failed to fetch whois: %w", err)
	}
	defer func() {
		d.PrevWhois = &d.CurrWhois
	}()

	l := d.Log()

	if d.CurrWhois.Domain.ExpirationDateInTime == nil {
		return ErrNoExpiration
	}

	d.ExpiresAt = d.CurrWhois.Domain.ExpirationDateInTime.Local()
	d.TimeLeft = time.Until(d.ExpiresAt).Truncate(24 * time.Hour)

	l.Info("Fetched whois", "expires", d.ExpiresAt, "days_left", d.TimeLeft.Hours()/24.0)

	if err := d.CheckNotifications(ctx, integrations); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// NotifyThreshold checks if the domain has reached a configured expiration threshold and sends a notification.
func (d *Domain) NotifyThreshold(ctx context.Context, integrations integration.Integrations) error {
	if d.TimeLeft != 0 {
		daysLeft := int(d.TimeLeft.Hours() / 24)
		for _, threshold := range d.conf.Threshold {
			if d.TriggeredThreshold != threshold && daysLeft <= threshold {
				msg := message.NewThresholdMessage(d.Name, daysLeft)
				integrations.Send(ctx, msg)
				d.TriggeredThreshold = threshold
				break
			}
		}
	}
	return nil
}

// NotifyStatusChange compares the current and previous WHOIS status and sends a notification if changed.
func (d *Domain) NotifyStatusChange(ctx context.Context, integrations integration.Integrations) error {
	if d.PrevWhois != nil {
		changes, err := diff.Diff(d.PrevWhois.Domain.Status, d.CurrWhois.Domain.Status)
		if err != nil {
			return err
		}

		if len(changes) > 0 {
			msg := message.NewStatusChangedMessage(d.Name, changes)
			integrations.Send(ctx, msg)
		}
	}
	return nil
}

// CheckNotifications runs all notification checks for the domain.
func (d *Domain) CheckNotifications(ctx context.Context, i integration.Integrations) error {
	if err := d.NotifyThreshold(ctx, i); err != nil {
		return err
	}
	if err := d.NotifyStatusChange(ctx, i); err != nil {
		return err
	}
	return nil
}
