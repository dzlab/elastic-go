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
	ShardsNumber    = "number_of_shards"
	ReplicasNumber  = "number_of_replicas"
	RefreshInterval = "refresh_interval"
	// analyzer params
	TOKENIZER      = "tokenizer"        //
	FILTER         = "filter"           //
	MinShingleSize = "min_shingle_size" //
	MaxShingleSize = "max_shingle_size" //
	OutputUnigrams = "output_unigrams"  //
)

type Index struct {
	url  string
	dict Dict
}

/*
 * Return a JSON representation of the body of this Index
 */
func (idx *Index) String() string {
	result, err := json.Marshal(idx.dict)
	if err != nil {
		log.Println(err)
	}
	return string(result)
}

func (client *Elasticsearch) Index(index string) *Index {
	url := fmt.Sprintf("http://%s/%s", client.Addr, index)
	return &Index{url: url, dict: make(Dict)}
}

/*
 * add a setting parameter
 */
func (idx *Index) Settings(settings Dict) {
	idx.dict[SETTINGS] = settings
}

/*
 * Set the mapping parameter
 */
func (idx *Index) Mappings(doctype string, mapping *Mapping) *Index {
	if idx.dict[MAPPINGS] == nil {
		idx.dict[MAPPINGS] = make(Dict)
	}
	idx.dict[MAPPINGS].(Dict)[doctype] = mapping.query
	return idx
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
func (idx *Index) SetAlias(alias string) *Index {
	idx.url += fmt.Sprintf("/%s/%s", ALIAS, alias)
	return idx
}

/*
 * Add a key-value settings
 */
func (idx *Index) AddSetting(name string, value interface{}) *Index {
	if idx.dict[SETTINGS] == nil {
		idx.dict[SETTINGS] = make(Dict)
	}
	idx.dict[SETTINGS].(Dict)[name] = value
	return idx
}

/*
 * Set the number of shards
 */
func (idx *Index) SetShardsNb(number int) *Index {
	idx.AddSetting(ShardsNumber, number)
	return idx
}

/*
 * Set the number of shards
 */
func (idx *Index) SetReplicasNb(number int) *Index {
	idx.AddSetting(ReplicasNumber, number)
	return idx
}

/*
 * Set the refresh interval
 */
func (idx *Index) SetRefreshInterval(interval string) *Index {
	idx.AddSetting(RefreshInterval, interval)
	return idx
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
func (analyzer *Analyzer) String() string {
	dict := make(Dict)
	dict[analyzer.name] = analyzer.kv
	return String(dict)
}

/*
 * Add an anlyzer to the index settings
 */
func (idx *Index) AddAnalyzer(analyzer *Analyzer) *Index {
	// if no "settings" create one
	if idx.dict[SETTINGS] == nil {
		idx.dict[SETTINGS] = make(Dict)
	}
	// if no "settings.analysis" create one
	if idx.dict[SETTINGS].(Dict)[ANALYSIS] == nil {
		idx.dict[SETTINGS].(Dict)[ANALYSIS] = make(Dict) //map[string]*Analyzer)
	}
	// insert the analyser ('name' and 'kv' attributes are taken separately)
	settings := idx.dict[SETTINGS].(Dict)
	analysis := settings[ANALYSIS].(Dict) //map[string]*Analyzer)
	analysis[analyzer.name] = analyzer.kv
	idx.dict[SETTINGS].(Dict)[ANALYSIS] = analysis
	return idx
}

/*
 * add an attribute to analyzer definition
 */
func (analyzer *Analyzer) Add1(key1, key2 string, value interface{}) *Analyzer {
	if len(analyzer.kv[key1]) == 0 {
		analyzer.kv[key1] = make(Dict)
	}
	analyzer.kv[key1][key2] = value
	return analyzer
}

/*
 * add a dictionary of attributes to analyzer definition
 */
func (analyzer *Analyzer) Add2(name string, value Dict) *Analyzer {
	if len(analyzer.kv[name]) == 0 {
		analyzer.kv[name] = make(Dict)
	}
	for k, v := range value {
		analyzer.kv[name][k] = v
	}
	return analyzer
}

/*
 * Pretify elasticsearch result
 */
func (analyzer *Index) Pretty() *Index {
	analyzer.url += "?pretty"
	return analyzer
}

/*
 * Create an index
 * PUT /:index
 */
func (analyzer *Index) Put() {
	query := String(analyzer.dict)
	log.Println("PUT", analyzer.url, query)
	reader, err := exec("PUT", analyzer.url, bytes.NewReader([]byte(query)))
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
func (idx *Index) Delete() {
	log.Println("DELETE", idx.url)
	reader, err := exec("DELETE", idx.url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	if data, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(data))
	}
}
