package jsonfield

import "time"

type Date time.Time

func (j Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format("2006-01-02") + `"`), nil
}
