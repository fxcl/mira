package datetime

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Date
type Date struct {
	time.Time
}

// MarshalJSON encodes to a custom JSON format.
func (d Date) MarshalJSON() ([]byte, error) {
	// Return null if the time is zero.
	if d.IsZero() {
		return []byte("null"), nil
	}

	return []byte("\"" + d.Format("2006-01-02") + "\""), nil
}

// UnmarshalJSON decodes the JSON format.
func (d *Date) UnmarshalJSON(data []byte) error {
	var err error

	if len(data) == 2 || string(data) == "null" {
		return err
	}

	var now time.Time

	// Custom format parsing.
	if now, err = time.ParseInLocation("2006-01-02", string(data), time.Local); err == nil {
		*d = Date{now}
		return err
	}

	// Custom format parsing with quotes.
	if now, err = time.ParseInLocation("\"2006-01-02\"", string(data), time.Local); err == nil {
		*d = Date{now}
		return err
	}

	return err
}

// Value converts to a database value.
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}

	return d.Time, nil
}

// Scan converts a database value to Date.
func (d *Date) Scan(value interface{}) error {
	if val, ok := value.(time.Time); ok {
		*d = Date{Time: val}
		return nil
	}

	return errors.New("cannot convert value to timestamp")
}
