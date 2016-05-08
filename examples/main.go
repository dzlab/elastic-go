package main

import (
	"flag"
)

var (
	name  = flag.String("run", "", "name of the samples to run")
	funcs = map[string]func(){
		"chap12": chap12,
		"chap13": chap13,
		"chap14": chap14,
		"chap15": chap15,
		"chap16": chap16,
		"chap17": chap17,
		"chap18": chap18,
		"chap19": chap19,
		"chap20": chap20,
		"chap21": chap21,
		"chap22": chap22,
		"chap23": chap23,
		"chap24": chap24,
		"chap26": chap26,
		"chap28": chap28,
		"chap29": chap29,
	}
)

func main() {
	flag.Parse()
	if funcs[*name] == nil {
		panic("Cannot run example: " + *name)
	}
	funcs[*name]()
}
