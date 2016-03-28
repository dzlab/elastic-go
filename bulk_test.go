package elastic

import (
	"testing"
)

// test for bulk
func TestBulk(t *testing.T) {
	actual := []string{
		newBulk().AddOperation(NewOperation(1).Add("price", 10).Add("productID", "XHDK-A-1293-#fJ3")).AddOperation(NewOperation(2).Add("price", 20).Add("productID", "KDKE-B-9947-#kL5")).String(),
	}
	expected := []string{
		`{"index":{"_id":1}}
{"price":10,"productID":"XHDK-A-1293-#fJ3"}
{"index":{"_id":2}}
{"price":20,"productID":"KDKE-B-9947-#kL5"}`,
	}
	equals(t, actual, expected)
}

// test for operation
func TestOperation(t *testing.T) {
	actual := []string{
		NewOperation(1).Add("other_field", "some data").String(),
		NewOperation(1).AddMultiple("tags", "search", "open_source").String(),
		NewOperation(1).AddMultiple("tags", "search", nil).String(),
		NewOperation(1).AddMultiple("tags").String(),
	}
	expected := []string{
		`{"other_field":"some data"}`,
		`{"tags":["search","open_source"]}`,
		`{"tags":["search",null]}`,
		`{"tags":null}`,
	}
	equals(t, actual, expected)
}
