package main

import (
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/doc"
)

const LIB_DIR = "./pkg/lua/lib"
const DOC_DIR = "./docs"

func main() {
	_, err := os.Stat(DOC_DIR)
	if err != nil {
		os.Mkdir(DOC_DIR, 0o666)
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
		if len(lib.Fns) == 0 && len(lib.Cns) == 0 {
			continue
		}

		outFile, err := os.OpenFile(path.Join("./docs", lib.Name+".md"), os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			panic(fmt.Sprintf("failed to open file to save docs: %s", path.Join("./docs", lib.Name+".md")))
		}
		defer outFile.Close()

		doc.Format(outFile, lib)
	}

}
