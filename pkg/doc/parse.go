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
				doc.Args = append(doc.Args, strings.TrimPrefix(line, TAG_ARG))
			}

			for ; strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_RETURNS); i++ {
				line := strings.TrimSpace(lines[i])
				doc.Returns = append(doc.Returns, strings.TrimPrefix(line, TAG_RETURNS))
			}

			if strings.HasPrefix(strings.TrimSpace(lines[i]), TAG_BLOCK) {
				doc.Block = true
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
