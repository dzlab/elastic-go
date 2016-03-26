package elastic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	MAPPING            = "mapping"
	MAPPINGS           = "mappings"
	CLASS              = "type"
	INDEX              = "index"
	PROPERTIES         = "properties"
	MATCH              = "match"
	MATCH_MAPPING_TYPE = "match_mapping_type"
	DYNAMIC_TEMPLATES  = "dynamic_templates"
	DEFAULT            = "_default_"
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
func newMapping(url string) *Mapping {
	return &Mapping{url: url, query: make(Dict)}
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/_mapping/:type
 */
func (this *Elasticsearch) Mapping(index, class string) *Mapping {
	var url string
	if class == "" {
		url = fmt.Sprintf("http://%s/%s", this.Addr, index)
	} else {
		url = fmt.Sprintf("http://%s/%s/%s/_mapping", this.Addr, index, class)
	}
	return newMapping(url)
}

/*
 * Get a string representation of this mapping API
 */
func (this *Mapping) String() string {
	return String(this.query)
}

/*
 * Add a mapping for a type's proeprty
 */
func (this *Mapping) AddProperty(name, class, index string) *Mapping {
	if this.query[PROPERTIES] == nil {
		this.query[PROPERTIES] = make(Dict)
	}
	property := make(Dict)
	property[CLASS] = class
	property[INDEX] = index
	this.query[PROPERTIES].(Dict)[name] = property
	return this
}

/*
 * Add a mapping for a type of objects
 */
func (this *Mapping) AddDocumentType(class *DocType) *Mapping {
	if this.query[MAPPINGS] == nil {
		this.query[MAPPINGS] = Dict{}
	}
	this.query[MAPPINGS].(Dict)[class.name] = class.dict
	return this
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/_mapping/:type
 */
func (this *Mapping) Get() {
	log.Println("GET", this.url)
	reader, err := exec("GET", this.url, nil)
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
func (this *Mapping) Put() {
	body := this.String()
	data := bytes.NewReader([]byte(body))
	log.Println("PUT", this.url)
	reader, err := exec("PUT", this.url, data)
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
func (this *DocType) AddProperty(name string, value interface{}) *DocType {
	this.dict[name] = value
	return this
}

/*
 * Add a template to this document type
 */
func (this *DocType) AddTemplate(tmpl *Template) *DocType {
	this.dict[tmpl.name] = tmpl.dict
	return this
}

/*
 * Add a dynamic template to this mapping
 */
func (this *DocType) AddDynamicTemplate(tmpl *Template) *DocType {
	if this.dict[DYNAMIC_TEMPLATES] == nil {
		this.dict[DYNAMIC_TEMPLATES] = []Dict{}
	}
	dict := make(Dict)
	dict[tmpl.name] = tmpl.dict
	this.dict[DYNAMIC_TEMPLATES] = append(this.dict[DYNAMIC_TEMPLATES].([]Dict), dict)
	return this
}

/*
 * Get a string representation of this document type
 */
func (this *DocType) String() string {
	dict := make(Dict)
	dict[this.name] = this.dict
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
func (this *Template) AddMatch(match string) *Template {
	this.dict[MATCH] = match
	return this
}

/*
 * Add a property to this template
 */
func (this *Template) AddProperty(name string, value interface{}) *Template {
	this.dict[name] = value
	return this
}

/*
 * Add a property to the `mapping` object
 */
func (this *Template) AddMappingProperty(name string, value interface{}) *Template {
	if this.dict[MAPPING] == nil {
		this.dict[MAPPING] = make(Dict)
	}
	this.dict[MAPPING].(Dict)[name] = value
	return this
}

/*
 * Get a string representation of this template
 */
func (this *Template) String() string {
	dict := make(Dict)
	dict[this.name] = this.dict
	return String(dict)
}
