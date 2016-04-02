package elastic

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Dict map[string]interface{}

// fields of a Search API call
const (
	// APIs
	EXPLAIN  = "explain"
	VALIDATE = "validate"
	SEARCH   = "search"
	// Query elements
	ALL     = "_all"
	INCLUDE = "include_in_all"
	SOURCE  = "_source"
	// url params
	SEARCH_TYPE = "search_type"
	SCROLL      = "scroll"
)

/*
 * a request representing a search
 */
type Search struct {
	url    string
	params map[string]string
	query  Dict
}

type Query interface {
	Name() string
	KV() Dict
}

/*
 * General purpose query
 */
type Object struct {
	name string
	kv   Dict
}

func (this *Object) Name() string {
	return this.name
}

func (this *Object) KV() Dict {
	return this.kv
}

/*
 * Create a new query object
 */
func NewQuery(name string) *Object {
	return &Object{name: name, kv: make(Dict)}
}

/*
 * Get a string representation of this object
 */
func (this *Object) String() string {
	return String(this.KV())
}

/*
 * Create an Explain request, that will return explanation for why a document is returned by the query
 */
func (this *Elasticsearch) Explain(index, class string, id int64) *Search {
	var url string = this.request(index, class, id, EXPLAIN)
	return &Search{url: url, query: make(Dict)}
}

/*
 * Create a Validate request
 */
func (this *Elasticsearch) Validate(index, class string, explain bool) *Search {
	var url string = this.request(index, class, -1, VALIDATE) + "/query"
	if explain {
		url += "?" + EXPLAIN
	}
	return &Search{url: url, query: make(Dict)}
}

/*
 * Create a Search request
 */
func (this *Elasticsearch) Search(index, class string) *Search {
	var url string = this.request(index, class, -1, SEARCH)
	return &Search{url: url, params: make(map[string]string), query: make(Dict)}
}

/*
 * Create a new Search API call
 */
func newSearch() *Search {
	return &Search{url: "", params: make(map[string]string), query: make(Dict)}
}

/*
 * Add a url parameter/value, e.g. search_type (count, query_and_fetch, dfs_query_then_fetch/dfs_query_and_fetch, scan)
 */
func (this *Search) AddParam(name, value string) *Search {
	this.params[name] = value
	return this
}

/*
 * Pretiffy the response result
 */
func (this *Search) Pretty() *Search {
	this.AddParam("pretty", "")
	return this
}

/*
* Add a query to this search request
 */
func (this *Search) AddQuery(query Query) *Search {
	this.query[query.Name()] = query.KV()
	return this
}

/*
 * Add to _source (i.e. specify another field that should be extracted)
 */
func (this *Search) AddSource(source string) *Search {
	var sources []string
	if this.query[SOURCE] == nil {
		sources = []string{}
	} else {
		sources = this.query[SOURCE].([]string)
	}
	sources = append(sources, source)
	this.query[SOURCE] = sources
	return this
}

/*
 * Add a query argument/value, e.g. size, from, etc.
 */
func (this *Search) Add(argument string, value interface{}) *Search {
	this.query[argument] = value
	return this
}

/*
 * Get a string representation of this Search API call
 */
func (this *Search) String() string {
	body := ""
	if len(this.query) > 0 {
		body = String(this.query)
	}
	return body
}

/*
 * Construct the url of this Search API call
 */
func (this *Search) urlString() string {
	return urlString(this.url, this.params)
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/:type/_search
 */
func (this *Search) Get() {
	// construct the url
	url := this.urlString()
	// construct the body
	query := this.String()
	var body io.Reader
	if query != "" {
		body = bytes.NewReader([]byte(query))
	}
	// submit the request
	log.Println("GET", url, query)
	reader, err := exec("GET", url, body)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}

/*
 * Add a query argument/value
 */
func (this *Object) Add(argument string, value interface{}) *Object {
	this.kv[argument] = value
	return this
}

/*
 * specify multiple values to match
 */
func (this *Object) AddMultiple(argument string, values ...interface{}) *Object {
	this.kv[argument] = values
	return this
}

/*
 * Add a sub query (e.g. a field query)
 */
func (this *Object) AddQuery(query Query) *Object {
	this.kv[query.Name()] = query.KV()
	return this
}

/*
 * Boolean clause, it is a complex clause that allows to combine other clauses as 'must' match, 'must_not' match, 'should' match.
 */
type Bool struct {
	name string
	kv   Dict
}

func (this *Bool) Name() string {
	return this.name
}
func (this *Bool) KV() Dict {
	return this.kv
}

/*
 * Create a new 'bool' clause
 */
func NewBool() *Bool {
	kv := make(Dict)
	return &Bool{name: "bool", kv: kv}
}

/*
 * Add a 'must' clause to this 'bool' clause
 */
func (this *Bool) AddMust(query Query) *Bool {
	this.add("must", query)
	return this
}

/*
 * Add a 'must_not' clause to this 'bool' clause
 */
func (this *Bool) AddMustNot(query Query) *Bool {
	this.add("must_not", query)
	return this
}

/*
 * Add a 'should' clause to this 'bool' clause
 */
func (this *Bool) AddShould(query Query) *Bool {
	this.add("should", query)
	return this
}

/*
 * Add a parameter to this `bool` query
 */
func (this *Bool) Add(name string, value interface{}) *Bool {
	this.kv[name] = value
	return this
}

/*
 * add a clause
 */
func (this *Bool) add(key string, query Query) {
	collection := this.kv[key]
	// check if query.Name exists, otherwise transform the map to array
	if collection == nil {
		// at first the collection is a map
		collection = make(Dict)
		collection.(Dict)[query.Name()] = query.KV()
	} else {
		// when more items are added, then it becomes an array
		dict := make(Dict)
		dict[query.Name()] = query.KV()
		// check if it is a map
		if _, ok := collection.(Dict); ok {
			array := []Dict{} // transform previous map into array
			for k, v := range collection.(Dict) {
				d := make(Dict)
				d[k] = v
				array = append(array, d)
			}
			collection = array
		}
		collection = append(collection.([]Dict), dict)
	}
	this.kv[key] = collection
}

/*
 * Create a new 'terms' filter, it is like 'term' but can match multiple values
 */
func NewTerms() *Object {
	return NewQuery("terms")
}

/*
 * Create a new 'term' filter
 */
func NewTerm() *Object {
	return NewQuery("term")
}

/*
 * Create a new `exists` filter.
 */
func NewExists() *Object {
	return NewQuery("exists")
}

/*
 * Create a new `missing` filter (the inverse of `exists`)
 */
func NewMissing() *Object {
	return NewQuery("missing")
}
