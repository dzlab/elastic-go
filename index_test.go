package elastic

import (
	"testing"
)

// test for queries with settings
func TestSettings(t *testing.T) {
	// given input
	actual := []string{
		newIndex().SetShardsNb(1).String(),
		newIndex().AddSetting("number_of_replicas", 0).String(),
		newIndex().AddAnalyzer(NewAnalyzer("char_filter").Add1("my_stopwords", "type", "stop").Add1("my_stopwords", "stopwords", []string{"the", "a"})).String(),
	}
	// expected result
	expected := []string{
		`{"settings":{"number_of_shards":1}}`,
		`{"settings":{"number_of_replicas":0}}`,
		`{"settings":{"analysis":{"char_filter":{"my_stopwords":{"stopwords":["the","a"],"type":"stop"}}}}}`,
	}
	// check if OK
	equals(t, actual, expected)
}

// test for analyzers
func TestAnalyzer(t *testing.T) {
	actual := []string{
		NewAnalyzer("char_filter").Add1("my_stopwords", "type", "stop").Add1("my_stopwords", "stopwords", []string{"the", "a"}).String(),
	}
	expected := []string{
		`{"char_filter":{"my_stopwords":{"stopwords":["the","a"],"type":"stop"}}}`,
	}
	// check if OK
	equals(t, actual, expected)
}
