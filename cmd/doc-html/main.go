package main

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/ArtificialLegacy/imgscal/pkg/doc"
)

const LIB_DIR = "./pkg/lua/lib"
const DOC_DIR = "./docs"

func main() {
	_, err := os.Stat(DOC_DIR)
	if err != nil {
		os.Mkdir(DOC_DIR, 0o777)
	}

	fs, err := os.ReadDir(LIB_DIR)
	if err != nil {
		panic(fmt.Sprintf("cannot open lua lib dir with err: %s", err))
	}

	docs := []doc.Lib{}

	for _, f := range fs {
		bs, err := os.ReadFile(path.Join(LIB_DIR, f.Name()))
		if err != nil {
			panic(fmt.Sprintf("failed file read on file %s", f.Name()))
		}

		docs = append(docs, doc.Parse(f.Name(), bs))
	}

	for _, lib := range docs {
		tmpl, err := template.New("doc.html").ParseFiles("./pkg/doc/doc.html")
		if err != nil {
			panic(fmt.Sprintf("failed to create tmpl: %s", err))
		}

		var f *os.File
		f, err = os.OpenFile(path.Join(DOC_DIR, strings.TrimSuffix(lib.File, ".go")+".html"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(f, lib)
		if err != nil {
			panic(err)
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}

}
