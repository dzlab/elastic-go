elastic-go
==============
[![Build Status](https://travis-ci.org/dzlab/elastic-go.png)](https://travis-ci.org/dzlab/elastic-go)

elastic-go is a golang client library that wraps Elasticsearch REST API. It currently has support for:
* Search
* Index
* Mapping
* Analyze
* ... more to come

### Installation
```go get github.com/dzlab/elastic-go```

### Documentation
http://godoc.org/github.com/dzlab/elastic-go

### Usage
```
import e "github.com/dzlab/elastic-go"
...
client := &e.Elasticsearch{Addr: "localhost:9200"}
client.Search("", "").Add("from", 30).Add("size", 10).Get()
// create an index example
client.Index("my_index").Delete()
cf := e.NewAnalyzer("char_filter").Add1("&_to_and", "type", "mapping").Add2("&_to_and", map[string]interface{}{"mappings": []string{"&=> and "}})
f := e.NewAnalyzer("filter").Add2("my_stopwords", map[string]interface{}{"type": "stop", "stopwords": []string{"the", "a"}})
a := e.NewAnalyzer("analyzer").Add2("my_analyzer", e.Dict{"type": "custom", "char_filter": []string{"html_strip", "&_to_and"}, "tokenizer": "standard", "filter": []string{"lowercase", "my_stopwords"}})
client.Index("my_index").AddAnalyzer(cf).AddAnalyzer(f).AddAnalyzer(a).Put()
```

### Contribute
This library is still under very active development. Any contribution is welcome.

Some planned features:

* A REPL to interact easily with Elasticsearch
* ...
