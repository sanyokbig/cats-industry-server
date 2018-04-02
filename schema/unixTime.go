package schema

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Type to parse from eve api
type UnixTime uint64

func (t *UnixTime) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(time.RFC3339, string(data[1:len(data)-1]))
	if err != nil {
		log.Debugf("failed to parse data: ", err)
		*t = 0
		return nil
	}
	*t = UnixTime(parsed.Unix())
	return nil
}
