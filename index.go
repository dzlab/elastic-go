package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	ANALYSIS = "analysis"
	SETTINGS = "settings"
	ALIAS    = "_alias"
	// settings params
	SHARDS_NB        = "number_of_shards"
	REPLICAS_NB      = "number_of_replicas"
	REFRESH_INTERVAL = "refresh_interval"
	// analyzer params
	TOKENIZER        = "tokenizer"        //
	FILTER           = "filter"           //
	MIN_SHINGLE_SIZE = "min_shingle_size" //
	MAX_SHINGLE_SIZE = "max_shingle_size" //
	OUTPUT_UNIGRAMS  = "output_unigrams"  //
)

type Index struct {
	url  string
	dict Dict
}

/*
 * Return a JSON representation of the body of this Index
 */
func (this *Index) String() string {
	result, err := json.Marshal(this.dict)
	if err != nil {
		log.Println(err)
	}
	return string(result)
}

func (this *Elasticsearch) Index(index string) *Index {
	url := fmt.Sprintf("http://%s/%s", this.Addr, index)
	return &Index{url: url, dict: make(Dict)}
}

/*
 * add a setting parameter
 */
func (this *Index) Settings(settings Dict) {
	this.dict[SETTINGS] = settings
}

/*
 * Set the mapping parameter
 */
func (this *Index) Mappings(doctype string, mapping *Mapping) *Index {
	if this.dict[MAPPINGS] == nil {
		this.dict[MAPPINGS] = make(Dict)
	}
	this.dict[MAPPINGS].(Dict)[doctype] = mapping.query
	return this
}

/*
 * Create new index settings
 */
func newIndex() *Index {
	return &Index{dict: make(Dict)}
}

/*
 * Define an alias for this index
 */
func (this *Index) SetAlias(alias string) *Index {
	this.url += fmt.Sprintf("/%s/%s", ALIAS, alias)
	return this
}

/*
 * Add a key-value settings
 */
func (this *Index) AddSetting(name string, value interface{}) *Index {
	if this.dict[SETTINGS] == nil {
		this.dict[SETTINGS] = make(Dict)
	}
	this.dict[SETTINGS].(Dict)[name] = value
	return this
}

/*
 * Set the number of shards
 */
func (this *Index) SetShardsNb(number int) *Index {
	this.AddSetting(SHARDS_NB, number)
	return this
}

/*
 * Set the number of shards
 */
func (this *Index) SetReplicasNb(number int) *Index {
	this.AddSetting(REPLICAS_NB, number)
	return this
}

/*
 * Set the refresh interval
 */
func (this *Index) SetRefreshInterval(interval string) *Index {
	this.AddSetting(REFRESH_INTERVAL, interval)
	return this
}

/*
 * Analyzer/Filter
 */
type Analyzer struct {
	name string
	kv   map[string]Dict
}

/*
 * Create a new analyzer
 */
func NewAnalyzer(name string) *Analyzer {
	return &Analyzer{name: name, kv: make(map[string]Dict)}
}

/*
 * Return a JSON string representation of this analyzer
 */
func (this *Analyzer) String() string {
	dict := make(Dict)
	dict[this.name] = this.kv
	return String(dict)
}

/*
 * Add an anlyzer to the index settings
 */
func (this *Index) AddAnalyzer(analyzer *Analyzer) *Index {
	// if no "settings" create one
	if this.dict[SETTINGS] == nil {
		this.dict[SETTINGS] = make(Dict)
	}
	// if no "settings.analysis" create one
	if this.dict[SETTINGS].(Dict)[ANALYSIS] == nil {
		this.dict[SETTINGS].(Dict)[ANALYSIS] = make(Dict) //map[string]*Analyzer)
	}
	// insert the analyser ('name' and 'kv' attributes are taken separately)
	settings := this.dict[SETTINGS].(Dict)
	analysis := settings[ANALYSIS].(Dict) //map[string]*Analyzer)
	analysis[analyzer.name] = analyzer.kv
	this.dict[SETTINGS].(Dict)[ANALYSIS] = analysis
	return this
}

/*
 * add an attribute to analyzer definition
 */
func (this *Analyzer) Add1(key1, key2 string, value interface{}) *Analyzer {
	if len(this.kv[key1]) == 0 {
		this.kv[key1] = make(Dict)
	}
	this.kv[key1][key2] = value
	return this
}

/*
 * add a dictionary of attributes to analyzer definition
 */
func (this *Analyzer) Add2(name string, value Dict) *Analyzer {
	if len(this.kv[name]) == 0 {
		this.kv[name] = make(Dict)
	}
	for k, v := range value {
		this.kv[name][k] = v
	}
	return this
}

/*
 * Pretify elasticsearch result
 */
func (this *Index) Pretty() *Index {
	this.url += "?pretty"
	return this
}

/*
 * Create an index
 * PUT /:index
 */
func (this *Index) Put() {
	query := String(this.dict)
	log.Println("PUT", this.url, query)
	reader, err := exec("PUT", this.url, bytes.NewReader([]byte(query)))
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}

}

/*
 * delete an index
 * DELETE /:index
 */
func (this *Index) Delete() {
	log.Println("DELETE", this.url)
	reader, err := exec("DELETE", this.url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}
