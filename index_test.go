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
		NewAnalyzer("filter").Add1("my_stopwords", "type", "stop").Add1("my_stopwords", "stopwords", []string{"the", "a"}).String(),
		NewAnalyzer("char_filter").Add1("&_to_and", "type", "mapping").Add2("&_to_and", Dict{"mappings": []string{"&=> and "}}).String(),
		NewAnalyzer("analyzer").Add2("my_analyzer", Dict{"char_filter": []string{"html_strip", "&_to_and"}, "filter": []string{"lowercase", "my_stopwords"}, "tokenizer": "standard", "type": "custom"}).String(),
	}
	expected := []string{
		`{"filter":{"my_stopwords":{"stopwords":["the","a"],"type":"stop"}}}`,
		`{"char_filter":{"&_to_and":{"mappings":["&=> and "],"type":"mapping"}}}`,
		`{"analyzer":{"my_analyzer":{"char_filter":["html_strip","&_to_and"],"filter":["lowercase","my_stopwords"],"tokenizer":"standard","type":"custom"}}}`,
	}
	// check if OK
	equals(t, actual, expected)
}
