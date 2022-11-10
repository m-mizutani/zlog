package zlog_test

import (
	"testing"

	"github.com/m-mizutani/zlog"
	"github.com/stretchr/testify/assert"
)

func TestLog_OrderedKeys(t *testing.T) {
	log := zlog.Log{
		Values: make(map[string]any),
	}
	log.Values["a"] = 1
	log.Values["b"] = 1
	log.Values["c"] = 1
	log.Values["d"] = 1
	log.Values["e"] = 1

	ordered := log.OrderedKeys()
	assert.Equal(t, "a", ordered[0])
	assert.Equal(t, "b", ordered[1])
	assert.Equal(t, "c", ordered[2])
	assert.Equal(t, "d", ordered[3])
	assert.Equal(t, "e", ordered[4])
}
