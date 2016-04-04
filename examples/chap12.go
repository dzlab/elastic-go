package main

import (
	e "github.com/dzlab/elastic-go"
)

/*
 * Examples of queries based on Elasticsearch Definitive Guide, chapter 12
 * Structured search examples
 */
func main() {
	client := &e.Elasticsearch{Addr: "localhost:9200"}
	// index
	client.Index("my_store").Delete()
	client.Index("my_store").Mappings("products", e.NewMapping("").AddProperty("productID", "string", "not_analyzed")).Put()
	// Bulk
	client.Bulk("my_store", "products").AddOperation(e.NewOperation(1).Add("price", 10).Add("productID", "XHDK-A-1293-#fJ3")).AddOperation(e.NewOperation(2).Add("price", 20).Add("productID", "KDKE-B-9947-#kL5")).AddOperation(e.NewOperation(3).Add("price", 30).Add("productID", "JODL-X-1937-#pV7")).AddOperation(e.NewOperation(4).Add("price", 30).Add("productID", "QQPX-R-3956-#aD8")).Post()
	// Search
	client.Search("my_store", "products").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("filtered").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("match_all"))).AddQuery(e.NewQuery("filter").AddQuery(e.NewQuery("term").Add("price", 30))))).Get()
	// analyze
	client.Analyze("my_store").Field("productID").Get("XHDK-A-1293-#fJ3")
	// search
	client.Search("my_store", "products").AddQuery(e.NewQuery("query").AddQuery(e.NewQuery("filtered").AddQuery(e.NewQuery("filter").AddQuery(e.NewBool().AddShould(e.NewQuery("term").Add("price", 20)).AddShould(e.NewQuery("term").Add("productID", "XHDK-A-1293-#fJ3")).AddMustNot(e.NewTerm().Add("price", 30)))))).Get()
}
