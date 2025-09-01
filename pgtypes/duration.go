package pgtypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) Scan(src interface{}) error {
	if src == nil {
		d.Duration = 0
		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("cannot scan type %T into Duration", src)
	}

	parsed, err := parsePostgresInterval(s)
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

func (d Duration) Value() (driver.Value, error) {
	// Stored in Go duration string format, e.g., '1h2m3s'
	return d.String(), nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Duration) UnmarshalText(data []byte) error {
	parsed, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

func (Duration) GormDataType() string {
	return "interval"
}

func (Duration) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	if db.Dialector.Name() == "postgres" {
		return "interval"
	}
	return ""
}

func FromDuration(d time.Duration) Duration {
	return Duration{d}
}

func (d Duration) AsDuration() time.Duration {
	return d.Duration
}

func (d Duration) String() string {
	return d.Duration.String()
}

func (d Duration) Equals(other Duration) bool {
	return d.Duration == other.Duration
}

// parsePostgresInterval attempts to convert a PostgreSQL interval string into time.Duration.
// It only supports the time portion (ignores months/years).
func parsePostgresInterval(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	// Try to parse using Go's time.ParseDuration first
	if dur, err := time.ParseDuration(s); err == nil {
		return dur, nil
	}

	// Try parsing Postgres-style HH:MM:SS or D days HH:MM:SS
	// This won't support months/years, which are ambiguous
	parts := strings.Fields(s)
	var total time.Duration
	var offset int

	if len(parts) == 3 && parts[1] == "days" {
		days, err := time.ParseDuration(parts[0] + "24h")
		if err != nil {
			return 0, err
		}
		total += days
		offset = 2
	}

	t, err := time.Parse("15:04:05", parts[offset])
	if err != nil {
		return 0, errors.New("unsupported interval format: " + s)
	}

	total += time.Duration(t.Hour()) * time.Hour
	total += time.Duration(t.Minute()) * time.Minute
	total += time.Duration(t.Second()) * time.Second
	return total, nil
}
