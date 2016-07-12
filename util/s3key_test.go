package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyNameStartsWithTarws(t *testing.T) {
	key := S3KeyName("path", time.Now())

	assert.Regexp(t, "^tarws", key)
}

func TestKeyContainsDate(t *testing.T) {
	date := time.Date(2016, time.January, 31, 4, 44, 44, 44, time.UTC)
	key := S3KeyName("path", date)

	assert.Contains(t, key, "20160131-0444")
}

func TestKeyContainsPath(t *testing.T) {
	key := S3KeyName(".", time.Now())

	assert.Contains(t, key, "util")
	assert.NotRegexp(t, "\\/", key, "key must not contain path separators")
}
