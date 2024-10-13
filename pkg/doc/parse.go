package doc

import (
	"fmt"
	"strings"
)

func Parse(filename string, file []byte) Lib {
	name := strings.TrimSuffix(filename, ".go")
	docs := Lib{File: filename, Name: name, Display: name}
	docs.FileClean = strings.TrimSuffix(docs.File, ".go")
	if docs.FileClean == "imgscal" {
		docs.FileClean = "index"
	}

	lines := strings.Split(string(file), "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, TAG_LIB) {
			docs.Display = strings.TrimPrefix(line, TAG_LIB)

			i++
			if strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_IMPORT) {
				docs.Name = strings.TrimPrefix(strings.TrimSpace(lines[i]), TAG_IMPORT)
				i++
			}

			if strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_DESC) {
				i++
				ln := strings.TrimSpace(lines[i])
				for strings.HasPrefix(ln, TAG_EMPTY) && !strings.HasPrefix(ln, TAG_EXISTS) {
					docs.Desc = append(docs.Desc, strings.TrimPrefix(ln, TAG_EMPTY))
					i++
					ln = strings.TrimSpace(lines[i])
				}
			}

			for strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_SECTION) {
				i++
				sec := []string{}
				ln := strings.TrimSpace(lines[i])
				for strings.HasPrefix(ln, TAG_EMPTY) && !strings.HasPrefix(ln, TAG_EXISTS) {
					sec = append(sec, strings.TrimPrefix(ln, TAG_EMPTY))
					i++
					ln = strings.TrimSpace(lines[i])
				}
				docs.Scs = append(docs.Scs, sec)
			}
		}

		line = strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, TAG_FUNC) {
			doc := Fn{Block: false}
			doc.Fn = strings.TrimPrefix(line, TAG_FUNC)
			doc.Name = strings.Split(doc.Fn, "(")[0]

			i++
			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_ARG); i++ {
				line := strings.TrimSpace(lines[i])
				d := Arg{}

				if strings.HasPrefix(line, TAG_ARG_REQ) {
					d.Str = strings.TrimPrefix(line, TAG_ARG_REQ)
					d.Opt = false
				} else if strings.HasPrefix(line, TAG_ARG_OPT) {
					d.Str = strings.TrimPrefix(line, TAG_ARG_OPT)
					d.Opt = true
				}

				t := strings.FieldsFunc(d.Str, func(r rune) bool {
					return r == '{' || r == '}'
				})

				d.Str = t[0]
				d.Type = t[1]
				if len(t) > 2 {
					d.Desc = " " + strings.Join(t[2:], "")
				}

				doc.Args = append(doc.Args, d)
			}

			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_RETURNS); i++ {
				line := strings.TrimSpace(lines[i])
				t := strings.FieldsFunc(strings.TrimPrefix(line, TAG_RETURNS), func(r rune) bool {
					return r == '{' || r == '}'
				})
				d := Return{}
				d.Type = t[0]
				if len(t) > 1 {
					d.Str = t[1]
				}

				doc.Returns = append(doc.Returns, d)
			}

			if strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_BLOCK) {
				doc.Block = true
				i++
			}

			if strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_DESC) {
				i++
				ln := strings.TrimSpace(lines[i])
				for strings.HasPrefix(ln, TAG_EMPTY) && !strings.HasPrefix(ln, TAG_EXISTS) {
					doc.Desc = append(doc.Desc, strings.TrimPrefix(ln, TAG_EMPTY))
					i++
					ln = strings.TrimSpace(lines[i])
				}
			}

			docs.Fns = append(docs.Fns, doc)
		} else if strings.HasPrefix(line, TAG_CONSTANTS) {
			doc := Const{}
			group := strings.FieldsFunc(strings.TrimPrefix(line, TAG_CONSTANTS), func(r rune) bool {
				return r == '{' || r == '}'
			})
			doc.Group = group[0]
			doc.Type = group[1]

			i++
			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_CONST); i++ {
				line := strings.TrimSpace(lines[i])
				doc.Consts = append(doc.Consts, strings.TrimPrefix(line, TAG_CONST))
			}

			docs.Cns = append(docs.Cns, doc)
		} else if strings.HasPrefix(line, TAG_STRUCT) {
			doc := Struct{}
			doc.Struct = strings.TrimPrefix(line, TAG_STRUCT)

			i++
			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_PROP); i++ {
				line := strings.TrimSpace(lines[i])
				d := Prop{}
				t := strings.FieldsFunc(strings.TrimPrefix(line, TAG_PROP), func(r rune) bool {
					return r == '{' || r == '}'
				})

				d.Str = t[0]
				d.Type = t[1]
				if len(t) > 2 {
					d.Desc = t[2]
				}

				doc.Props = append(doc.Props, d)
			}

			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_METHOD); i++ {
				line := strings.TrimSpace(lines[i])
				d := Method{}
				t := strings.Split(strings.TrimPrefix(line, TAG_METHOD), " - ")

				d.Name = strings.Split(t[0], "(")[0]
				d.Type = t[0]
				if len(t) > 1 {
					d.Desc = strings.Join(t[1:], " - ")
				}

				doc.Methods = append(doc.Methods, d)
			}

			if strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_DESC) {
				i++
				ln := strings.TrimSpace(lines[i])
				for strings.HasPrefix(ln, TAG_EMPTY) && !strings.HasPrefix(ln, TAG_EXISTS) {
					doc.Desc = append(doc.Desc, strings.TrimPrefix(ln, TAG_EMPTY))
					i++
					ln = strings.TrimSpace(lines[i])
				}
			}

			docs.Sts = append(docs.Sts, doc)
		} else if strings.HasPrefix(line, TAG_INTERFACE) {
			doc := Interface{}
			doc.Interface = strings.TrimPrefix(line, TAG_INTERFACE)

			i++
			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_PROP); i++ {
				line := strings.TrimSpace(lines[i])
				d := Prop{}
				t := strings.FieldsFunc(strings.TrimPrefix(line, TAG_PROP), func(r rune) bool {
					return r == '{' || r == '}'
				})

				d.Str = t[0]
				d.Type = t[1]
				if len(t) > 2 {
					d.Desc = t[2]
				}

				doc.Props = append(doc.Props, d)
			}

			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_METHOD); i++ {
				line := strings.TrimSpace(lines[i])
				d := Method{}
				t := strings.Split(strings.TrimPrefix(line, TAG_METHOD), " - ")

				d.Name = strings.Split(t[0], "(")[0]
				d.Type = t[0]
				if len(t) > 1 {
					d.Desc = strings.Join(t[1:], " - ")
				}

				doc.Methods = append(doc.Methods, d)
			}

			if strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_DESC) {
				i++
				ln := strings.TrimSpace(lines[i])
				for strings.HasPrefix(ln, TAG_EMPTY) && !strings.HasPrefix(ln, TAG_EXISTS) {
					doc.Desc = append(doc.Desc, strings.TrimPrefix(ln, TAG_EMPTY))
					i++
					ln = strings.TrimSpace(lines[i])
				}
			}

			docs.Its = append(docs.Its, doc)
		} else if strings.HasPrefix(line, TAG_EXISTS) || strings.HasPrefix(line, TAG_EMPTY) {
			fmt.Printf("Unknown doc tag: %s in %s\n", line, filename)
		} else if strings.HasPrefix(line, TAG_INCORRECT) {
			fmt.Printf("Possible invalid doc tag: %s in %s\n", line, filename)
		}

	}

	return docs
}
