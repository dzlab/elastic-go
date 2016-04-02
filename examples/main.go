package main

import (
	e "github.com/dzlab/elastic-go"
)

func main() {
	client := &e.Elasticsearch{Addr: "localhost:9200"}
	client.Search("", "").Add("from", 30).Add("size", 10).Get()
	client.Search("", "").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match").Add("tweet", "elasticsearch"))).Get()
	client.Search("index_2014*", "type1,type2").Get()
	// reindex in batch
	client.Search("old_index", "").AddParam("search_type", "scan").AddParam("scroll", "1m").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("range").AddQuery(e.NewQuery("data").Add("gte", "2014-01-01").Add("lt", "2014-02-01")))).Add("size", 1000).Get()
	client.Index("gb").Delete()
	// Bulk
	client.Index("my_store").Delete()
	client.Index("my_store").Mappings("products", e.NewMapping("").AddProperty("productID", "string", "not_analyzed")).Put()
	client.Bulk("my_store", "products").AddOperation(e.NewOperation(1).Add("price", 10).Add("productID", "XHDK-A-1293-#fJ3")).AddOperation(e.NewOperation(2).Add("price", 20).Add("productID", "KDKE-B-9947-#kL5")).AddOperation(e.NewOperation(3).Add("price", 30).Add("productID", "JODL-X-1937-#pV7")).AddOperation(e.NewOperation(4).Add("price", 30).Add("productID", "QQPX-R-3956-#aD8")).Post()
	client.Search("my_store", "products").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("filtered").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match_all"))).AddQuery(e.NewQuery("filter").AddQuery(e.NewQuery("term").Add("price", 30))))).Get()
	client.Analyze("my_store").Field("productID").Get("XHDK-A-1293-#fJ3")
	client.Search("my_store", "products").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("filtered").AddQuery(e.NewQuery("filter").AddQuery(e.NewBool().AddShould(e.NewQuery("term").Add("price", 20)).AddShould(e.NewQuery("term").Add("productID", "XHDK-A-1293-#fJ3")).AddMustNot(e.NewTerm().Add("price", 30)))))).Get()
	//"{mappings: {tweet: {properties: {tweet:{type: \"string\", analyzer: \"english\"}, date: {type: \"date\"}, name: {type: \"string\"}, user_id: {type: \"long\"}}}}}"
	client.Index("gb_v1").Put()
	client.Index("gb_v1").SetAlias("gb").Put()
	// create Aliases operations
	client.Alias().AddAction("remove", "my_index_v1", "my_index").AddAction("add", "my_index_v2", "my_index").Post()
	// create an index example
	client.Index("my_index").Delete()
	cf := e.NewAnalyzer("char_filter").Add1("&_to_and", "type", "mapping").Add2("&_to_and", map[string]interface{}{"mappings": []string{"&=> and "}})
	f := e.NewAnalyzer("filter").Add2("my_stopwords", map[string]interface{}{"type": "stop", "stopwords": []string{"the", "a"}})
	a := e.NewAnalyzer("analyzer").Add2("my_analyzer", e.Dict{"type": "custom", "char_filter": []string{"html_strip", "&_to_and"}, "tokenizer": "standard", "filter": []string{"lowercase", "my_stopwords"}})
	client.Index("my_index").AddAnalyzer(cf).AddAnalyzer(f).AddAnalyzer(a).Put()
	client.Analyze("my_index").Analyzer("my_analyzer").Get("The quick & brown fox")
	// create mapping
	client.Mapping("gb", "tweet").AddProperty("tag", "string", "not_analyzed").Put()
	client.Mapping("gb", "tweet").Get()
	client.Analyze("gb").Get("tweet")
	client.Analyze("gb").Get("tag")
	client.Validate("gb", "tweet", true).AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("tweet").Add("match", "really powerful"))).Get()
	// sort by date
	client.Search("", "").AddQuery(
		e.NewQuery("query").AddQuery(e.NewQuery("filtered").AddQuery(e.NewQuery("filter").AddQuery(e.NewQuery("term").Add("user_id", "1")))),
	).AddQuery(
		e.NewQuery("sort").AddQuery(e.NewQuery("date").Add("order", "desc")),
	).Get()
}
