package log

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

func DEBUG(formating string, args ...interface{}) {
	LOG("DEBUG", formating, args...)
}
func INFO(formating string, args ...interface{}) {
	LOG("INFO", formating, args...)
}
func ERROR(formating string, args ...interface{}) {
	LOG("ERROR", formating, args...)
}

func LOG(level string, formating string, args ...interface{}) {
	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(2)
	// fmt.Println(reflect.TypeOf(pc), reflect.ValueOf(pc))
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo

		filename = filepath.Base(filename) // /full/path/basename.go => basename.go
	}

	log.Printf("[%5s] %10s:%-3d : %20s : %s\n", level, filename, line, funcname, fmt.Sprintf(formating, args...))
}
