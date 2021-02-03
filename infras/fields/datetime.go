package fields

import (
	"time"
)

type DateTime time.Time

func (j DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format("2006-01-02 15:04:05") + `"`), nil
}
