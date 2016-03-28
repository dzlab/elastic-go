package elastic

import (
	"testing"
)

// test for Alias actions
func TestActions(t *testing.T) {
	// given input
	actual := []string{
		newAlias().AddAction("remove", "my_index_v1", "my_index").AddAction("add", "my_index_v2", "my_index").String(),
	}
	// expected result
	expected := []string{
		`{"actions":[{"remove":{"alias":"my_index","index":"my_index_v1"}},{"add":{"alias":"my_index","index":"my_index_v2"}}]}`,
	}
	equals(t, actual, expected)
}
