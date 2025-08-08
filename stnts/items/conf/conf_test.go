package conf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockConf() *Conf {
	return New("Name", "Desc", "Local")
}

func TestNew(t *testing.T) {
	// success
	conf := mockConf()
	assert.Equal(t, "Name", conf.Name)
	assert.Equal(t, "Desc", conf.Desc)
	assert.Equal(t, "Local", conf.TimeZone)
}

func TestLocation(t *testing.T) {
	// setup
	conf := mockConf()

	// success
	loca := conf.Location()
	assert.Equal(t, time.Local, loca)
}
