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
	}
	// expected result
	expected := []string{
		"/",
		"?k1=v1",
	}
	equals(t, actual, expected)
	// tests for cases where order doesn't matter
	str := urlString("/", Map{"k1": "", "k2": "v2"})
	if str != "/?k1&k2=v2" && str != "/?k2=v2&k1" {
		t.Errorf("%s should be equal to %s or %s", str, "/?k1&k2=v2", "/?k2=v2&k1")
	}
}
