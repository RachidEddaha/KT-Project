package entitiescustom

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ReleaseDate wraps time.Time making me able to format the date with a non standard format
type ReleaseDate struct{ time.Time }

const ReleaseDateFormat = "2006-01-02"

// Value returns a driver value
func (rd ReleaseDate) Value() (driver.Value, error) {
	b, err := json.Marshal(&rd)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan returns a parsed campaign
func (rd *ReleaseDate) Scan(src interface{}) error {
	srcTime := src.(time.Time)
	// fix format
	date, err := time.Parse(ReleaseDateFormat,
		srcTime.Format(ReleaseDateFormat))
	if err != nil {
		return err
	}
	rd.Time = date
	return nil
}

func (rd ReleaseDate) MarshalJSON() ([]byte, error) {
	date := rd.Time.Format(ReleaseDateFormat)
	date = fmt.Sprintf(`"%s"`, date)
	return []byte(date), nil
}

func (rd *ReleaseDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	date, err := time.Parse(ReleaseDateFormat, s)
	if err != nil {
		return err
	}
	rd.Time = date
	return
}
