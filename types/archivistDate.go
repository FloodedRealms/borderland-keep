package types

import (
	"encoding/json"
	"strings"
	"time"
)

type ArcvhistDate time.Time

func (d *ArcvhistDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	s = s[0:strings.LastIndex(s, "T")]
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = ArcvhistDate(t)
	return nil
}

func (d ArcvhistDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d))
}

func (d ArcvhistDate) Format(s string) string {
	t := time.Time(d)
	return t.Format(s)
}

func (d ArcvhistDate) Date() time.Time {
	format := "2006-01-01"
	t := time.Time(d)
	time, _ := time.Parse(format, t.Format(format))
	return time
}
