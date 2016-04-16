package main

import (
	e "github.com/dzlab/elastic-go"
	"time"
)

/*
 * Examples of queries based on Elasticsearch Definitive Guide, chapter 14
 * Multi-field search examples
 */
func chap14() {
	c := &e.Elasticsearch{Addr: "localhost:9200"}

	// write a condition for each field then gather them into a `bool` search query
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewQuery("match").Add("title", "War and Peace")).AddShould(e.NewQuery("match").Add("author", "Leo Tolstoy")))).Get()

	// the `bool` query is the mainstay for multi-clause queries
	// we can add a preference for the book version, each clause at the same level has same weight so use a separate clause to reduce the weight of the book version preference
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewQuery("match").Add("title", "War and Peace")).AddShould(e.NewQuery("match").Add("author", "Leo Tolstoy")).AddShould(e.NewBool().AddShould(e.NewQuery("match").Add("translator", "Constance Garnett")).AddShould(e.NewQuery("match").Add("translator", "Louise Maude"))))).Get()

	// we can also set an explicite weight for a clause via `boost` parameter
	// a reasonable range of `boost` value is between 1 and 10, upto 15
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewQuery("match").AddQuery(e.NewQuery("title").Add("query", "War and Peace").Add("boost", 2))).AddShould(e.NewQuery("match").AddQuery(e.NewQuery("author").Add("query", "Leo Tolstoy").Add("boost", 2))).AddShould(e.NewBool().AddShould(e.NewQuery("match").Add("translator", "Constance Garnett")).AddShould(e.NewQuery("match").Add("translator", "Louise Maude"))))).Get()

	// there is common search strategies: best fields, most fields, cross fields
	// Best fields search strategy
	c.Index("my_index").Delete()
	c.Insert("my_index", "my_type").Document(1, map[string]string{"title": "Quick brown rabbits", "body": "Brown rabbits are commonly seen."}).Put()
	c.Insert("my_index", "my_type").Document(2, map[string]string{"title": "Keeping pets healthy", "body": "My quick brown fox eats rabbits on a regular basis"}).Put()

	// wait the documents to be searchable
	time.Sleep(1 * time.Second)

	// searching in title and body
	// doc 1 will have higher score as `bool` query sums score of subqueries, multiply by the number of matching clauses, then divide by total number of clauses
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewQuery("match").Add("title", "Brown fox")).AddShould(e.NewQuery("match").Add("body", "Brown fox")))).Get()

	// insead of `bool` query use `dis_max` (Disjunction, i.e. `or`, Max Query) to return documents that match any of the given query
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("dis_max").AddQueries("queries", e.NewQuery("match").Add("title", "Brown fox"), e.NewQuery("match").Add("body", "Brown fox")))).Get()

	// tuning best fields queries
	// `dis_max` query simply uses `_score` from best matches,
	//it's possible to take into account `_score` from other matching clauses via `tie_breaker` param which will be multiplied by the matching clauses
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("dis_max").AddQueries("queries", e.NewQuery("match").Add("title", "Brown pets"), e.NewQuery("match").Add("body", "Brown pets")).Add("tie_breaker", 0.3))).Get()

	// `multi_match` query as a concise rewrite of previous queries
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("multi_match").Add("query", "Quick brown fox").Add("type", "best_fields").AddMultiple("fields", "title", "body").Add("tie_breaker", 0.3).Add("minimum_should_match", "30%"))).Get()

	// use of wildcards in field names
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("multi_match").Add("query", "Quick brown fox").Add("fields", "*_title"))).Get()

	// boosting individual fields by adding the ^boost after field name
	c.Search("my_index", "my_type").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("multi_match").Add("query", "Quick brown fox").AddMultiple("fields", "*_title", "chapter_title^2"))).Get()

	// Multi-field mapping, combing stemming analyzer (e.g. english) with standard analyzer
	c.Index("my_index").Delete()
	c.Index("my_index").SetShardsNb(1).Mappings("my_type", e.NewMapping().AddProperty("title", "type", "string").AddProperty("title", "analyzer", "english").AddProperty("title", "fields", e.Dict{"std": e.Dict{"type": "string", "analyzer": "standard"}})).Put()
	c.Insert("my_index", "my_type").Document(1, e.Dict{"title": "My rabbit jumps"}).Put()
	c.Insert("my_index", "my_type").Document(2, e.Dict{"title": "Jumping jack rabbits"}).Put()
	time.Sleep(1 * time.Second)
	c.Search("my_index", "").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("title", "jumping rabbits"))).Get()
	// query using title.std field, only document 2 will match
	c.Search("my_index", "").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewMatch().Add("title.std", "jumping rabbits"))).Get()
	// query both fields and combine their scores with `bool` query
	c.Search("my_index", "").Pretty().AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "jumping rabbits").Add("type", "most_fields").AddMultiple("fields", "title", "title.std"))).Get()

	// cross-fields entity search
	// naive approach: a `bool` query to sum up the score for each matched field
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewBool().AddShould(e.NewMatch().Add("street", "Poland Street W1V")).AddShould(e.NewMatch().Add("city", "Poland Street W1V")).AddShould(e.NewMatch().Add("country", "Poland Street W1V")).AddShould(e.NewMatch().Add("postcode", "Poland Street W1V")))).Pretty().Get()

	// or use a multi_match to avoid repeating the query
	c.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "Poland Street W1V").Add("type", "most_fields").AddMultiple("fields", "street", "city", "country", "postcode"))).Pretty().Get()
	// elasticsearch is generating a match query for each field, we can check this
	c.Validate("", "", true).AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "Poland Street W1V").Add("type", "most_fields").AddMultiple("fields", "street", "city", "country", "postcode"))).Pretty().Get()

	// custom _all fields to copy to it values of a combination of fields and search on them
	c.Index("my_index").Delete()
	c.Index("my_index").Mappings("person", e.NewMapping().AddProperty("first_name", "type", "string").AddProperty("first_name", "copy_to", "full_name").AddProperty("last_name", "type", "string").AddProperty("last_name", "copy_to", "full_name").AddProperty("full_name", "type", "string")).Pretty().Put()

	// multi_match query with type `cross_fields` to search on combined fields instead of modifying mapping (as previous) for each possible combination
	// check difference of how the musti_match search is broken for most_fields then cross_fields
	c.Validate("", "", true).AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "peter smith").Add("type", "most_fields").Add("operator", "and").AddMultiple("fields", "first_name", "last_name"))).Pretty().Get()
	c.Validate("", "", true).AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "peter smith").Add("type", "cross_fields").Add("operator", "and").AddMultiple("fields", "first_name", "last_name"))).Pretty().Get()
	// bossting a field (e.g. title) over less relevant fields (e.g. description)
	c.Search("book", "").AddQuery(e.NewQuery("query").AddQuery(e.NewMultiMatch().Add("query", "peter smith").Add("type", "cross_fields").AddMultiple("fields", "title^2", "description"))).Pretty().Get()

	// exact value fields (i.e. having 'not_analyzed' analyzer mappting) should not be used with `multi_match` queries as it will search for query field as a single term
}
