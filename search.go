package elastic

import ()

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
	SearchType = "search_type"
	SCROLL     = "scroll"
	// query names
	DisMax            = "dis_max"
	MultiMatch        = "multi_match"
	MatchPhrase       = "match_phrase" // 'phrase' search query
	MatchPhrasePrefix = "match_phrase_prefix"
	Prefix            = "prefix"   // search terms with given prefix
	Wildcard          = "wildcard" // search terms with widcard
	RegExp            = "regexp"   // filter terms application to regular expression
	RESCORE           = "rescore"  // rescore result of previous query
	RescoreQuery      = "rescroe_query"
	// query params
	MinimumShouldMatch = "minimum_should_match"
	SLOP               = "slop"           // in 'phrase' queries to describe proximity/word ordering
	MaxExpansions      = "max_expansions" // controls how many terms the prefix is allowed to match
	WindowSize         = "window_size"    // number of document from each shard
)

/*
 * a request representing a search
 */
type Search struct {
	client *Elasticsearch
	parser *SearchResultParser
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

/*
 * Name() returns the name of this query object
 */
func (obj *Object) Name() string {
	return obj.name
}

/*
 * KV() returns the key-value store representing the body of this query
 */
func (obj *Object) KV() Dict {
	return obj.kv
}

/*
 * NewQuery Create a new query object
 */
func NewQuery(name string) *Object {
	return &Object{name: name, kv: make(Dict)}
}

/*
 * NewMatch Create a new match query
 */
func NewMatch() *Object {
	return NewQuery(MATCH)
}

/*
 * NewMultiMatch Create a new multi_match query
 */
func NewMultiMatch() *Object {
	return NewQuery(MultiMatch)
}

/*
 * NewMatchPhrase Create a `match_phrase` query to find words that are near each other
 */
func NewMatchPhrase() *Object {
	return NewQuery(MatchPhrase)
}

/*
 * NewRescore Create a `rescore` query
 */
func NewRescore() *Object {
	return NewQuery(RESCORE)
}

/*
 * NewRescoreQuery Create a `rescore` query algorithm
 */
func NewRescoreQuery() *Object {
	return NewQuery(RescoreQuery)
}

/*
 * newQuery used for test purpose
 */
func newQuery() *Object {
	return &Object{name: "", kv: make(Dict)}
}

/*
 * String get a string representation of this object
 */
func (obj *Object) String() string {
	return String(obj.KV())
}

/*
 * Explain create an Explaination request, that will return explanation for why a document is returned by the query
 */
func (client *Elasticsearch) Explain(index, class string, id int64) *Search {
	url := client.request(index, class, id, EXPLAIN)
	return newSearch(client, url)
}

/*
 * Validate create a Validation request
 */
func (client *Elasticsearch) Validate(index, class string, explain bool) *Search {
	url := client.request(index, class, -1, VALIDATE) + "/query"
	if explain {
		url += "?" + EXPLAIN
	}
	return newSearch(client, url)
}

/*
 * Create a Search request
 */
func (client *Elasticsearch) Search(index, class string) *Search {
	url := client.request(index, class, -1, SEARCH)
	return newSearch(client, url)
}

/*
 * Create a new Search API call
 */
func newSearch(client *Elasticsearch, url string) *Search {
	return &Search{
		client: client,
		parser: &SearchResultParser{},
		url:    url,
		params: make(map[string]string),
		query:  make(Dict),
	}
}

/*
 * Add a url parameter/value, e.g. search_type (count, query_and_fetch, dfs_query_then_fetch/dfs_query_and_fetch, scan)
 */
func (search *Search) AddParam(name, value string) *Search {
	search.params[name] = value
	return search
}

/*
 * Pretiffy the response result
 */
func (search *Search) Pretty() *Search {
	search.AddParam("pretty", "")
	return search
}

/*
* Add a query to this search request
 */
func (search *Search) AddQuery(query Query) *Search {
	search.query[query.Name()] = query.KV()
	return search
}

/*
 * Add to _source (i.e. specify another field that should be extracted)
 */
func (search *Search) AddSource(source string) *Search {
	var sources []string
	if search.query[SOURCE] == nil {
		sources = []string{}
	} else {
		sources = search.query[SOURCE].([]string)
	}
	sources = append(sources, source)
	search.query[SOURCE] = sources
	return search
}

/*
 * Add a query argument/value, e.g. size, from, etc.
 */
func (search *Search) Add(argument string, value interface{}) *Search {
	search.query[argument] = value
	return search
}

/*
 * Get a string representation of this Search API call
 */
func (search *Search) String() string {
	body := ""
	if len(search.query) > 0 {
		body = String(search.query)
	}
	return body
}

/*
 * Construct the url of this Search API call
 */
func (search *Search) urlString() string {
	return urlString(search.url, search.params)
}

/*
 * request mappings between the json fields and how Elasticsearch store them
 * GET /:index/:type/_search
 */
func (search *Search) Get() {
	// construct the url
	url := search.urlString()
	// construct the body
	query := search.String()

	search.client.Execute("GET", url, query, search.parser)
}

/*
 * Add a query argument/value
 */
func (obj *Object) Add(argument string, value interface{}) *Object {
	obj.kv[argument] = value
	return obj
}

/*
 * specify multiple values to match
 */
func (obj *Object) AddMultiple(argument string, values ...interface{}) *Object {
	obj.kv[argument] = values
	return obj
}

/*
 * Add multiple queries, under given `name`
 */
func (obj *Object) AddQueries(name string, queries ...Query) *Object {
	for _, q := range queries {
		parent := NewQuery(name)
		parent.AddQuery(q)
		obj.AddQuery(parent)
	}
	return obj
}

/*
 * Add a sub query (e.g. a field query)
 */
func (obj *Object) AddQuery(query Query) *Object {
	collection := obj.kv[query.Name()]
	// check if query.Name exists, otherwise transform the map to array
	if collection == nil {
		// at first the collection is a map
		collection = query.KV()
	} else {
		// when more items are added, then it becomes an array
		dict := query.KV()
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
	obj.kv[query.Name()] = collection
	return obj
}

/*
 * Boolean clause, it is a complex clause that allows to combine other clauses as 'must' match, 'must_not' match, 'should' match.
 */
type Bool struct {
	name string
	kv   Dict
}

/*
 * Name() returns the name of this 'bool' query
 */
func (b *Bool) Name() string {
	return b.name
}

/*
 * KV() returns the key-value store representing the body of this 'bool' query
 */
func (b *Bool) KV() Dict {
	return b.kv
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
func (b *Bool) AddMust(query Query) *Bool {
	b.add("must", query)
	return b
}

/*
 * Add a 'must_not' clause to this 'bool' clause
 */
func (b *Bool) AddMustNot(query Query) *Bool {
	b.add("must_not", query)
	return b
}

/*
 * Add a 'should' clause to this 'bool' clause
 */
func (b *Bool) AddShould(query Query) *Bool {
	b.add("should", query)
	return b
}

/*
 * Add a parameter to this `bool` query
 */
func (b *Bool) Add(name string, value interface{}) *Bool {
	b.kv[name] = value
	return b
}

/*
 * add a clause
 */
func (b *Bool) add(key string, query Query) {
	collection := b.kv[key]
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
	b.kv[key] = collection
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
