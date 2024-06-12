package doc

import (
	"fmt"
	"io"
)

func Format(out io.StringWriter, lib Lib) {
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
}
