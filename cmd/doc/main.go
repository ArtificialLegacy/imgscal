package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
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

	docs := []LibStruct{}

	for _, f := range fs {
		bs, err := os.ReadFile(path.Join(LIB_DIR, f.Name()))
		if err != nil {
			panic(fmt.Sprintf("failed file read on file %s", f.Name()))
		}

		docs = append(docs, parseFile(f.Name(), bs))
	}

	for _, lib := range docs {
		if len(lib.Docs) == 0 {
			continue
		}

		out := bytes.Buffer{}

		out.WriteString(fmt.Sprintf("# %s\n", lib.Name))

		for _, fn := range lib.Docs {
			out.WriteString(fmt.Sprintf("\n## %s\n\n", fn.Fn))

			if len(fn.Args) > 0 {
				out.WriteString(fmt.Sprintf("### Args [%s]\n\n", fn.Fn))
				for _, arg := range fn.Args {
					out.WriteString(fmt.Sprintf("* %s\n", arg))
				}
			}

			if len(fn.Returns) > 0 {
				out.WriteString("\n")
				out.WriteString(fmt.Sprintf("### Returns [%s]\n\n", fn.Fn))
				for _, arg := range fn.Returns {
					out.WriteString(fmt.Sprintf("* %s\n", arg))
				}
			}
		}

		outFile, err := os.OpenFile(path.Join("./docs", lib.Name+".md"), os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			panic(fmt.Sprintf("failed to open file to save docs: %s", path.Join("./docs", lib.Name+".md")))
		}
		defer outFile.Close()

		outFile.Write(out.Bytes())
	}

}

type LibStruct struct {
	Name string
	Docs []DocStruct
}

type DocStruct struct {
	Fn      string
	Args    []string
	Returns []string
}

func parseFile(name string, file []byte) LibStruct {
	name = strings.TrimSuffix(name, ".go")
	docs := LibStruct{Name: name}

	lines := strings.Split(string(file), "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if !strings.HasPrefix(line, "/// @func ") {
			continue
		}

		doc := DocStruct{}
		doc.Fn = strings.TrimPrefix(line, "/// @func ")

		i++
		for ; strings.HasPrefix(strings.TrimSpace(lines[i]), "/// @arg "); i++ {
			line := strings.TrimSpace(lines[i])
			doc.Args = append(doc.Args, strings.TrimPrefix(line, "/// @arg "))
		}

		i++
		for ; strings.HasPrefix(strings.TrimSpace(lines[i]), "/// @returns "); i++ {
			line := strings.TrimSpace(lines[i])
			doc.Args = append(doc.Args, strings.TrimPrefix(line, "/// @returns "))
		}

		docs.Docs = append(docs.Docs, doc)
	}

	return docs
}
