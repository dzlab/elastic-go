package elastic

import (
	"testing"
)

type Map map[string]string

// test for urlString
func TestUrlString(t *testing.T) {
	// given input
	actual := []string{
		urlString("/", Map{}),
		urlString("?", Map{"k1": "v1"}),
		urlString("/", Map{"k1": "", "k2": "v2"}),
	}
	// expected result
	expected := []string{
		"/",
		"?k1=v1",
		"/?k1&k2=v2",
	}
	equals(t, actual, expected)
}
