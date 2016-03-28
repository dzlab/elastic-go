package elastic

import (
	"testing"
)

// test for mappings attributes
func TestMappings(t *testing.T) {
	// given input
	actual := []string{
		NewMapping("").AddProperty("tag", "string", "not_analyzed").String(),
		NewMapping("").AddDocumentType(NewDefaultType().AddTemplate(NewAllTemplate().AddProperty("enabled", false))).String(),
		NewMapping("").AddDocumentType(NewDocType("my_type").AddDynamicTemplate(NewTemplate("es").AddMatch("_es").AddProperty(MATCH_MAPPING_TYPE, "string").AddMappingProperty("type", "string").AddMappingProperty("analyzer", "spanish"))).String(),
	}
	// expected result
	expected := []string{
		`{"properties":{"tag":{"index":"not_analyzed","type":"string"}}}`,
		`{"mappings":{"_default_":{"_all":{"enabled":false}}}}`,
		`{"mappings":{"my_type":{"dynamic_templates":[{"es":{"mapping":{"analyzer":"spanish","type":"string"},"match":"_es","match_mapping_type":"string"}}]}}}`,
	}
	// check if OK
	equals(t, actual, expected)
}

// test for mappings templates
func TestDocType(t *testing.T) {
	actual := []string{
		NewDefaultType().AddTemplate(NewAllTemplate().AddProperty("enabled", false)).String(),
		NewDocType("my_type").AddDynamicTemplate(NewTemplate("es").AddMatch("_es").AddProperty(MATCH_MAPPING_TYPE, "string").AddMappingProperty("type", "string").AddMappingProperty("analyzer", "spanish")).String(),
		NewDocType("my_type").AddProperty("date_detection", false).String(),
	}
	expected := []string{
		`{"_default_":{"_all":{"enabled":false}}}`, // disable '_all' field in all types
		`{"my_type":{"dynamic_templates":[{"es":{"mapping":{"analyzer":"spanish","type":"string"},"match":"_es","match_mapping_type":"string"}}]}}`,
		`{"my_type":{"date_detection":false}}`,
	}
	equals(t, actual, expected)
}

// test for mappings templates
func TestTemplates(t *testing.T) {
	actual := []string{
		NewAllTemplate().AddProperty("enabled", false).String(),
		NewTemplate("es").AddMatch("_es").AddProperty(MATCH_MAPPING_TYPE, "string").AddMappingProperty("type", "string").AddMappingProperty("analyzer", "spanish").String(),
	}
	expected := []string{
		`{"_all":{"enabled":false}}`, // disable '_all' field
		`{"es":{"mapping":{"analyzer":"spanish","type":"string"},"match":"_es","match_mapping_type":"string"}}`,
	}
	equals(t, actual, expected)
}
