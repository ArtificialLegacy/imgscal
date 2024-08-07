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
		out.WriteString(fmt.Sprintf("\n### %s\n", fn.Fn))

		if fn.Block {
			out.WriteString("\n")
			out.WriteString("**❗ Note: This function is blocking, it will interupt concurrent execution.**\n")
		}

		if len(fn.Desc) > 0 {
			out.WriteString("\n")
			for _, d := range fn.Desc {
				out.WriteString(fmt.Sprintf("%s\n", d))
			}
		}

		if len(fn.Args) > 0 {
			out.WriteString("\n")
			out.WriteString(fmt.Sprintf("#### Args [%s]\n\n", fn.Fn))
			for _, arg := range fn.Args {
				if arg.Opt {
					out.WriteString(fmt.Sprintf("* *\\*%s*\n", arg.Str))
				} else {
					out.WriteString(fmt.Sprintf("* %s\n", arg.Str))
				}
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

	if len(lib.Sts) != 0 {
		out.WriteString("\n## Structs\n")
	}

	for _, st := range lib.Sts {
		out.WriteString(fmt.Sprintf("\n### %s\n", st.Struct))

		if len(st.Desc) > 0 {
			out.WriteString("\n")
			for _, d := range st.Desc {
				out.WriteString(fmt.Sprintf("%s\n", d))
			}
		}

		if len(st.Props) > 0 {
			out.WriteString("\n")
			out.WriteString(fmt.Sprintf("#### Props [%s]\n\n", st.Struct))
			for _, prop := range st.Props {
				out.WriteString(fmt.Sprintf("* %s\n", prop))
			}
		}

		if len(st.Methods) > 0 {
			out.WriteString("\n")
			out.WriteString(fmt.Sprintf("#### Methods [%s]\n\n", st.Struct))
			for _, m := range st.Methods {
				out.WriteString(fmt.Sprintf("* %s\n", m))
			}
		}
	}
}
