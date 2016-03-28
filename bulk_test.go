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
