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

const (
	TAG_FUNC      = "/// @func "
	TAG_ARG       = "/// @arg "
	TAG_RETURNS   = "/// @returns "
	TAG_CONSTANTS = "/// @constants "
	TAG_CONST     = "/// @const "
)

func main() {
	_, err := os.Stat(DOC_DIR)
	if err != nil {
		os.Mkdir(DOC_DIR, 0o666)
	}

	fs, err := os.ReadDir(LIB_DIR)
	if err != nil {
		panic(fmt.Sprintf("cannot open lua lib dir with err: %s", err))
	}

	docs := []Lib{}

	for _, f := range fs {
		bs, err := os.ReadFile(path.Join(LIB_DIR, f.Name()))
		if err != nil {
			panic(fmt.Sprintf("failed file read on file %s", f.Name()))
		}

		docs = append(docs, parseFile(f.Name(), bs))
	}

	for _, lib := range docs {
		if len(lib.Fns) == 0 && len(lib.Cns) == 0 {
			continue
		}

		out := bytes.Buffer{}

		out.WriteString(fmt.Sprintf("# %s\n", lib.Name))

		if len(lib.Fns) != 0 {
			out.WriteString("\n## Functions\n")
		}

		for _, fn := range lib.Fns {
			out.WriteString(fmt.Sprintf("\n### %s\n\n", fn.Fn))

			if len(fn.Args) > 0 {
				out.WriteString(fmt.Sprintf("#### Args [%s]\n\n", fn.Fn))
				for _, arg := range fn.Args {
					out.WriteString(fmt.Sprintf("* %s\n", arg))
				}
			}

			if len(fn.Returns) > 0 {
				out.WriteString("\n")
				out.WriteString(fmt.Sprintf("#### Returns [%s]\n\n", fn.Fn))
				for _, arg := range fn.Returns {
					out.WriteString(fmt.Sprintf("* %s\n", arg))
				}
			}
		}

		if len(lib.Cns) != 0 {
			out.WriteString("\n## Constants\n")
		}

		for _, cn := range lib.Cns {
			out.WriteString(fmt.Sprintf("\n### %s\n\n", cn.Group))

			for _, con := range cn.Consts {
				out.WriteString(fmt.Sprintf("* %s\n", con))
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

type Lib struct {
	Name string
	Fns  []Fn
	Cns  []Const
}

type Fn struct {
	Fn      string
	Args    []string
	Returns []string
}

type Const struct {
	Group  string
	Consts []string
}

func parseFile(name string, file []byte) Lib {
	name = strings.TrimSuffix(name, ".go")
	docs := Lib{Name: name}

	lines := strings.Split(string(file), "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, TAG_FUNC) {
			doc := Fn{}
			doc.Fn = strings.TrimPrefix(line, TAG_FUNC)

			i++
			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_ARG); i++ {
				line := strings.TrimSpace(lines[i])
				doc.Args = append(doc.Args, strings.TrimPrefix(line, TAG_ARG))
			}

			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_RETURNS); i++ {
				line := strings.TrimSpace(lines[i])
				doc.Returns = append(doc.Returns, strings.TrimPrefix(line, TAG_RETURNS))
			}

			docs.Fns = append(docs.Fns, doc)
		} else if strings.HasPrefix(line, TAG_CONSTANTS) {
			doc := Const{}
			doc.Group = strings.TrimPrefix(line, TAG_CONSTANTS)

			i++
			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_CONST); i++ {
				line := strings.TrimSpace(lines[i])
				doc.Consts = append(doc.Consts, strings.TrimPrefix(line, TAG_CONST))
			}

			docs.Cns = append(docs.Cns, doc)
		}

	}

	return docs
}
