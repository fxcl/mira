package datetime

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Time
type Time struct {
	time.Time
}

// MarshalJSON encodes to a custom JSON format.
func (t Time) MarshalJSON() ([]byte, error) {

	// Return null if the time is zero.
	if t.IsZero() {
		return []byte("null"), nil
	}

	return []byte("\"" + t.Format("15:04:05") + "\""), nil
}

// UnmarshalJSON decodes the JSON format.
func (t *Time) UnmarshalJSON(data []byte) error {

	var err error

	if len(data) == 2 || string(data) == "null" {
		return err
	}

	var now time.Time

	// Custom format parsing.
	if now, err = time.ParseInLocation("15:04:05", string(data), time.Local); err == nil {
		*t = Time{now}
		return err
	}

	// Custom format parsing with quotes.
	if now, err = time.ParseInLocation("\"15:04:05\"", string(data), time.Local); err == nil {
		*t = Time{now}
		return err
	}

	return err
}

// Value converts to a database value.
func (t Time) Value() (driver.Value, error) {

	if t.IsZero() {
		return nil, nil
	}

	return t.Time, nil
}

// Scan converts a database value to Time.
func (t *Time) Scan(value interface{}) error {

	if val, ok := value.(time.Time); ok {
		*t = Time{Time: val}
		return nil
	}

	return errors.New("cannot convert value to timestamp")
}
