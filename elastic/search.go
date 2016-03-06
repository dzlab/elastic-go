package elastic

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Search struct {
	url   string
	query map[string]string
}

type Query interface {
	Name() string
	KV() map[string]string
}

/*
 * General purpose query
 */
type General struct {
	name string
	kv   map[string]string
}

func (this *General) Name() string {
	return this.name
}
func (this *General) KV() map[string]string {
	return this.kv
}

func NewQuery(name string) *General {
	return &General{name: name, kv: make(map[string]string)}
}

func (this *Elasticsearch) Search(index, class string) *Search {
	var url string
	if index == "" {
		url = fmt.Sprintf("http://%s/_search", this.Addr)
	} else if class == "" {
		url = fmt.Sprintf("http://%s/%s/_search", this.Addr, index)
	} else {
		url = fmt.Sprintf("http://%s/%s/%s/_search", this.Addr, index, class)
	}
	return &Search{url: url, query: make(map[string]string)}
}

/*
* Add a query to this search request
 */
func (this *Search) AddQuery(query Query) *Search {
	this.query[query.Name()] = String(query.KV())
	return this
}

/*
 * Add a query argument/value, e.g. size, from, etc.
 */
func (this *Search) Add(argument, value string) *Search {
	this.query[argument] = value
	return this
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/:type/_search
 */
func (this *Search) Get() {
	body := ""
	if len(this.query) > 0 {
		body = String(this.query)
	}
	log.Println("query", body)
	var data io.Reader = nil
	if body != "" {
		data = bytes.NewReader([]byte(body))
	}
	reader, err := exec("GET", this.url, data)
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
func (this *General) Add(argument, value string) *General {
	this.kv[argument] = value
	return this
}

/*
 * Add a sub query (e.g. a field query)
 */
func (this *General) AddQuery(query Query) *General {
	this.kv[query.Name()] = String(query.KV())
	return this
}

/*
 * Return a string representation of the dictionary
 */
func String(dict map[string]string) string {
	if len(dict) == 0 {
		return "{}"
	}
	var value string = "{"
	for k, v := range dict {
		value += "\"" + k + "\"" + ":"
		if v[0] == byte('{') {
			value += v + ","
		} else {
			value += "\"" + v + "\"" + ","

		}
	}
	value = value[:len(value)-1] + "}"
	return value
}

/*
 * Boolean clause, it is a complex clause that allows to combine other clauses as 'must' match, 'must_not' match, 'should' match.
 */
type Bool struct {
	name string
	kv   map[string]string
}

func (this *Bool) Name() string {
	return this.name
}
func (this *Bool) KV() map[string]string {
	return this.kv
}

/*
 * Create a new 'bool' clause
 */
func NewBool() *Bool {
	kv := make(map[string]string)
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
	var prefix string
	if this.kv[key] == "" {
		prefix = "{"
	} else {
		prefix = this.kv[key][:len(this.kv[key])-1] + ", "
	}
	this.kv[key] = prefix + "\"" + query.Name() + "\"" + ":" + String(query.KV()) + "}"
}
