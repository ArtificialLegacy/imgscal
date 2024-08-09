package doc

import "strings"

func Parse(name string, file []byte) Lib {
	name = strings.TrimSuffix(name, ".go")
	docs := Lib{Name: name}

	lines := strings.Split(string(file), "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, TAG_FUNC) {
			doc := Fn{Block: false}
			doc.Fn = strings.TrimPrefix(line, TAG_FUNC)

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

				doc.Args = append(doc.Args, d)
			}

			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_RETURNS); i++ {
				line := strings.TrimSpace(lines[i])
				doc.Returns = append(doc.Returns, strings.TrimPrefix(line, TAG_RETURNS))
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
			doc.Group = strings.TrimPrefix(line, TAG_CONSTANTS)

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
				doc.Props = append(doc.Props, strings.TrimPrefix(line, TAG_PROP))
			}

			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_METHOD); i++ {
				line := strings.TrimSpace(lines[i])
				doc.Methods = append(doc.Methods, strings.TrimPrefix(line, TAG_METHOD))
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
		}

	}

	return docs
}
