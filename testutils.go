package elastic

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

/*
 * Identify the caller method
 */
func caller() string {
	calldepth := 2
	// get the stack trace
	pc := make([]uintptr, calldepth) // at least 1 entry needed
	runtime.Callers(calldepth, pc)
	f := runtime.FuncForPC(pc[1])   // 0 wil be the logging function (e.g. Debugf, Info, etc.)
	file, line := f.FileLine(pc[1]) // see previous comment
	shortfile := file[strings.LastIndex(file, "/")+1:]
	method := f.Name()[strings.LastIndex(f.Name(), "/")+1:]
	// remove (*) when there is pointers
	method = strings.Replace(method, "(", "", -1)
	method = strings.Replace(method, "*", "", -1)
	method = strings.Replace(method, ")", "", -1)
	return fmt.Sprintf("%s:%d %s", shortfile, line, method)
}

/*
 * assert if all entries of arrays are equals
 */
func equals(t *testing.T, actual, expected []string) {
	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			from := "(" + caller() + ")"
			t.Errorf("%s Should be equal\n%s\n%s", from, actual[i], expected[i])
		}
	}
}