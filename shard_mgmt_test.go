package elastic

import "testing"

// test for shard management queries
func TestShardMgmtOp(t *testing.T) {
	// given input
	actual := []string{
		newShardMgmtOp(REFRESH).urlString(),
		newShardMgmtOp(FLUSH).AddParam("wait_for_ongoing", "").urlString(),
		newShardMgmtOp(OPTIMIZE).AddParam("max_num_segment", "1").urlString(),
	}
	// expected result
	expected := []string{
		`refresh`,
		`flush?wait_for_ongoing`,
		`optimize?max_num_segment=1`,
	}
	// check if OK
	equals(t, actual, expected)
}
