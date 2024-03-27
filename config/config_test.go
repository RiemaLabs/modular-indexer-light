package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteReadFile(t *testing.T) {
	before := LoadBlacklist()
	assert.Equal(t, []*Blacklist{}, before)

	data1 := &Blacklist{
		Evidence: &Evidence{
			Height: 1,
		},
	}

	data2 := &Blacklist{
		Evidence: &Evidence{
			Height: 2,
		},
	}
	AppendBlacklist(data1)
	AppendBlacklist(data2)

	after := LoadBlacklist()
	assert.Equal(t, []*Blacklist{data1, data2}, after)
	os.Remove(BlacklistFile)
}
