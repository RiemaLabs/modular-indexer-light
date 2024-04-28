package utils

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDurH_MarshalJSON(t *testing.T) {
	type (
		valueType   struct{ Timeout DurH }
		pointerType struct{ Timeout *DurH }
	)
	const defaultDur = 3 * time.Second

	{
		v := &valueType{Timeout: DurH{defaultDur}}
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		var d pointerType
		if err := json.Unmarshal(data, &d); err != nil {
			t.Fatal(err)
		}
		if d.Timeout.Duration != v.Timeout.Duration {
			t.Fatal(d.Timeout.Duration, v.Timeout.Duration)
		}
	}

	{
		v := &pointerType{Timeout: &DurH{defaultDur}}
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		var d valueType
		if err := json.Unmarshal(data, &d); err != nil {
			t.Fatal(err)
		}
		if d.Timeout.Duration != v.Timeout.Duration {
			t.Fatal(d.Timeout.Duration, v.Timeout.Duration)
		}
	}
}
