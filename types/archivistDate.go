package types

import (
	"encoding/json"
	"strings"
	"time"
)

type ArcvhistDate time.Time

const formatString = "2006-01-02"

func (d *ArcvhistDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(formatString, s)
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
	t := time.Time(d)
	time, _ := time.Parse(formatString, t.Format(formatString))
	return time
}

func (d ArcvhistDate) String() string {
	return time.Time(d).Format(formatString)
}
