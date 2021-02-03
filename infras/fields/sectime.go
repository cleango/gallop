package fields

import (
	"database/sql/driver"
	"fmt"
	"github.com/spf13/cast"
	"time"
)

// SecTime format json time field by myself
type SecTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t SecTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t SecTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time.Unix(), nil
}

// Scan valueof time.Time
func (t *SecTime) Scan(v interface{}) error {
	value := cast.ToInt64(cast.ToString(v))
	*t = SecTime{Time: time.Unix(value, 0)}
	return nil
}
