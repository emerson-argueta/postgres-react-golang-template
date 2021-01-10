package goal

import "time"

const (
	// TimestampFormat for Timestamp type
	TimestampFormat = "2006-01-02 15:04:05"
)

// Timestamp for a goal's message. The format of the timestamp is YYYY-MM-DD HH:MM:SS
type Timestamp string

// NewTimestamp created. If timestamp is nil then current timestamp returned.
// Error returned when timestamp is invalid.
func NewTimestamp(timestamp *string) (res Timestamp, e error) {
	if timestamp == nil {
		newTimestamp := time.Now().Format(TimestampFormat)
		return Timestamp(newTimestamp), nil
	}
	parsed, e := time.Parse(TimestampFormat, *timestamp)
	validatedTimestamp := parsed.Format(TimestampFormat)

	if e != nil {
		return "", e
	}

	return Timestamp(validatedTimestamp), nil
}

// ToString from Timestamp type
func (t *Timestamp) ToString() string {
	return string(*t)
}
