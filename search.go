package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Dict map[string]interface{}

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

func NewQuery(name string) *Object {
	return &Object{name: name, kv: make(Dict)}
}

func (this *Elasticsearch) request(index, class string, id int64, request string) string {
	var url string
	if index == "" {
		url = fmt.Sprintf("http://%s/_%s", this.Addr, request)
	} else if class == "" {
		url = fmt.Sprintf("http://%s/%s/_%s", this.Addr, index, request)
	} else if id < 0 {
		url = fmt.Sprintf("http://%s/%s/%s/_%s", this.Addr, index, class, request)
	} else {
		url = fmt.Sprintf("http://%s/%s/%s/%d/_%s", this.Addr, index, class, id, request)
	}
	return url
}

/*
 * Create an Explain request, that will return explanation for why a document is returned by the query
 */
func (this *Elasticsearch) Explain(index, class string, id int64) *Search {
	var url string = this.request(index, class, id, "explain")
	return &Search{url: url, query: make(Dict)}
}

/*
 * Create a Validate request
 */
func (this *Elasticsearch) Validate(index, class string, explain bool) *Search {
	var url string = this.request(index, class, -1, "validate") + "/query"
	if explain {
		url += "?explain"
	}
	return &Search{url: url, query: make(Dict)}
}

/*
 * Create a Search request
 */
func (this *Elasticsearch) Search(index, class string) *Search {
	var url string = this.request(index, class, -1, "search")
	return &Search{url: url, params: make(map[string]string), query: make(Dict)}
}

/*
 * Add a url parameter/value, e.g. search_type (count, query_and_fetch, dfs_query_then_fetch/dfs_query_and_fetch, scan)
 */
func (this *Search) AddParam(name, value string) *Search {
	this.params[name] = value
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
 * Add a query argument/value, e.g. size, from, etc.
 */
func (this *Search) Add(argument string, value interface{}) *Search {
	this.query[argument] = value
	return this
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/:type/_search
 */
func (this *Search) Get() {
	// construct the url
	url := this.url
	if len(this.params) > 0 {
		url += "?"
		for name, value := range this.params {
			url += name
			if value != "" {
				url += "=" + value
			}
			url += url + "&"
		}
		url = url[:len(url)-1]
	}
	// construct the body
	body := ""
	if len(this.query) > 0 {
		body = String(this.query)
	}
	var data io.Reader = nil
	if body != "" {
		data = bytes.NewReader([]byte(body))
	}
	// submit the request
	log.Println("GET", url)
	log.Println("query", body)
	reader, err := exec("GET", url, data)
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
func (this *Object) AddMultiple(argument string, values ...string) *Object {
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
 * Return a string representation of the dictionary
 */
func String(dict Dict) string {
	marshaled, err := json.Marshal(dict)
	if err != nil {
		log.Println(err)
	}
	return string(marshaled)
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
 * add a clause
 */
func (this *Bool) add(key string, query Query) {
	dict := this.kv[key]
	if dict == nil {
		dict = make(Dict)
	}
	dict.(Dict)[query.Name()] = query.KV()
	this.kv[key] = dict
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
