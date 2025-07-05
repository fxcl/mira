package datetime

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Datetime
type Datetime struct {
	time.Time
}

// MarshalJSON encodes to a custom JSON format.
func (d Datetime) MarshalJSON() ([]byte, error) {
	// Return null if the time is zero.
	if d.IsZero() {
		return []byte("null"), nil
	}

	return []byte("\"" + d.Format("2006-01-02 15:04:05") + "\""), nil
}

// UnmarshalJSON decodes the JSON format.
func (d *Datetime) UnmarshalJSON(data []byte) error {
	var err error

	if len(data) == 2 || string(data) == "null" {
		return err
	}

	var now time.Time

	// Custom format parsing.
	if now, err = time.ParseInLocation("2006-01-02 15:04:05", string(data), time.Local); err == nil {
		*d = Datetime{now}
		return err
	}

	// Custom format parsing with quotes.
	if now, err = time.ParseInLocation("\"2006-01-02 15:04:05\"", string(data), time.Local); err == nil {
		*d = Datetime{now}
		return err
	}

	// Default format parsing.
	if now, err = time.ParseInLocation(time.RFC3339, string(data), time.Local); err == nil {
		*d = Datetime{now}
		return err
	}

	if now, err = time.ParseInLocation("\""+time.RFC3339+"\"", string(data), time.Local); err == nil {
		*d = Datetime{now}
		return err
	}

	return err
}

// Value converts to a database value.
func (d Datetime) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}

	return d.Time, nil
}

// Scan converts a database value to Datetime.
func (d *Datetime) Scan(value interface{}) error {
	if val, ok := value.(time.Time); ok {
		*d = Datetime{Time: val}
		return nil
	}

	return errors.New("cannot convert value to timestamp")
}
