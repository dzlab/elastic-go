package elastic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	MAPPING           = "mapping"
	MAPPINGS          = "mappings"
	TYPE              = "type"
	INDEX             = "index"
	PROPERTIES        = "properties"
	MATCH             = "match"
	MatchMappingType  = "match_mapping_type"
	DynamicTemplates  = "dynamic_templates"
	DEFAULT           = "_default_"
	PositionOffsetGap = "position_offset_gap"
)

/*
 * mappings between the json fields and how Elasticsearch store them
 */
type Mapping struct {
	url   string
	query Dict
}

/*
 * Create a new mapping query
 */
func NewMapping(url string) *Mapping {
	return &Mapping{url: url, query: make(Dict)}
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/:type/_mapping
 */
func (client *Elasticsearch) Mapping(index, doctype string) *Mapping {
	var url string = client.request(index, doctype, -1, MAPPING)
	return NewMapping(url)
}

/*
 * Get a string representation of this mapping API
 */
func (mapping *Mapping) String() string {
	return String(mapping.query)
}

/*
 * Add a mapping for a type's property (e.g. type, index, analyzer, etc.)
 */
func (mapping *Mapping) AddProperty(fieldname, propertyname string, propertyvalue interface{}) *Mapping {
	if mapping.query[PROPERTIES] == nil {
		mapping.query[PROPERTIES] = make(Dict)
	}
	property := mapping.query[PROPERTIES].(Dict)[fieldname]
	if property == nil {
		property = make(Dict)
	}
	property.(Dict)[propertyname] = propertyvalue
	mapping.query[PROPERTIES].(Dict)[fieldname] = property
	return mapping
}

/*
 * Add a mapping for a type of objects
 */
func (mapping *Mapping) AddDocumentType(class *DocType) *Mapping {
	if mapping.query[MAPPINGS] == nil {
		mapping.query[MAPPINGS] = Dict{}
	}
	mapping.query[MAPPINGS].(Dict)[class.name] = class.dict
	return mapping
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/_mapping/:type
 */
func (mapping *Mapping) Get() {
	log.Println("GET", mapping.url)
	reader, err := exec("GET", mapping.url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}

/*
 * Update a mappings between the json fields and how Elasticsearch store them
 * PUT /:index/_mapping/:type
 */
func (mapping *Mapping) Put() {
	query := mapping.String()
	body := bytes.NewReader([]byte(query))
	log.Println("PUT", mapping.url, query)
	reader, err := exec("PUT", mapping.url, body)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}

/*
 * A document type
 */
type DocType struct {
	name string
	dict Dict
}

/*
 * Create a '_default_' type that encapsulates shared/default settings
 * e.g. specify index wide dynamic templates
 */
func NewDefaultType() *DocType {
	return NewDocType(DEFAULT)
}

/*
 * Create a new mapping template
 */
func NewDocType(name string) *DocType {
	return &DocType{name: name, dict: make(Dict)}
}

/*
 * Add property to this document type
 */
func (doctype *DocType) AddProperty(name string, value interface{}) *DocType {
	doctype.dict[name] = value
	return doctype
}

/*
 * Add a template to this document type
 */
func (doctype *DocType) AddTemplate(tmpl *Template) *DocType {
	doctype.dict[tmpl.name] = tmpl.dict
	return doctype
}

/*
 * Add a dynamic template to this mapping
 */
func (doctype *DocType) AddDynamicTemplate(tmpl *Template) *DocType {
	if doctype.dict[DynamicTemplates] == nil {
		doctype.dict[DynamicTemplates] = []Dict{}
	}
	dict := make(Dict)
	dict[tmpl.name] = tmpl.dict
	doctype.dict[DynamicTemplates] = append(doctype.dict[DynamicTemplates].([]Dict), dict)
	return doctype
}

/*
 * Get a string representation of this document type
 */
func (doctype *DocType) String() string {
	dict := make(Dict)
	dict[doctype.name] = doctype.dict
	return String(dict)
}

/*
 * A mapping template
 */
type Template struct {
	name string
	dict Dict
}

func NewAllTemplate() *Template {
	return NewTemplate(ALL)
}

/*
 * Create a new mapping template
 */
func NewTemplate(name string) *Template {
	return &Template{name: name, dict: make(Dict)}
}

/*
 * Add a match string (e.g. '*', '_es')
 */
func (template *Template) AddMatch(match string) *Template {
	template.dict[MATCH] = match
	return template
}

/*
 * Add a property to this template
 */
func (template *Template) AddProperty(name string, value interface{}) *Template {
	template.dict[name] = value
	return template
}

/*
 * Add a property to the `mapping` object
 */
func (template *Template) AddMappingProperty(name string, value interface{}) *Template {
	if template.dict[MAPPING] == nil {
		template.dict[MAPPING] = make(Dict)
	}
	template.dict[MAPPING].(Dict)[name] = value
	return template
}

/*
 * Get a string representation of this template
 */
func (template *Template) String() string {
	dict := make(Dict)
	dict[template.name] = template.dict
	return String(dict)
}
