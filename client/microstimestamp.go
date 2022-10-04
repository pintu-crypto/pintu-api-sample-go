package client

import (
	"strconv"
	"time"
)

// MicrosTimestamp is a timestamp that marshals
// to JSON as an RFC-3339 UTC time with microsecond precision.
type MicrosTimestamp time.Time

// MicrosTimestampFormatString is similar to time.RFC3339Nano but with reduced precision
// to the microsecond level and fixed to 6dp
const microsTimestampFormatString = "2006-01-02T15:04:05.000000Z07:00"

func (t MicrosTimestamp) String() string {
	return time.Time(t).UTC().Format(microsTimestampFormatString)
}

// MarshalJSON converts the timestamp to a JSON string.
func (t MicrosTimestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.String())), nil
}

// UnmarshalJSON parses a quoted JSON time string.
func (t *MicrosTimestamp) UnmarshalJSON(data []byte) (err error) {
	var value string
	value, err = strconv.Unquote(string(data))
	if err != nil {
		return
	}
	var tm time.Time
	tm, err = time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return err
	}
	*t = MicrosTimestamp(tm)
	return
}
