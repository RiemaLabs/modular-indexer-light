package utils

import (
	"encoding/json"
	"time"
)

// DurH is the JSON wrapper of time.Duration. The `H` here means *humanized*.
type DurH struct {
	time.Duration
}

func (d *DurH) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	if d != nil {
		d.Duration = dur
	}
	return nil
}

func (d *DurH) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}
