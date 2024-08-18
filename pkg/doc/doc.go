package doc

const (
	TAG_EMPTY     = "/// "
	TAG_EXISTS    = "/// @"
	TAG_LIB       = "/// @lib "
	TAG_IMPORT    = "/// @import "
	TAG_FUNC      = "/// @func "
	TAG_ARG       = "/// @arg"
	TAG_ARG_REQ   = "/// @arg "
	TAG_ARG_OPT   = "/// @arg? "
	TAG_RETURNS   = "/// @returns "
	TAG_CONSTANTS = "/// @constants "
	TAG_CONST     = "/// @const "
	TAG_BLOCK     = "/// @blocking"
	TAG_DESC      = "/// @desc"
	TAG_STRUCT    = "/// @struct "
	TAG_PROP      = "/// @prop "
	TAG_METHOD    = "/// @method "
	TAG_INCORRECT = "// @"
	TAG_SECTION   = "/// @section"
)

type Lib struct {
	File      string
	FileClean string
	Name      string
	Display   string
	Desc      []string
	Fns       []Fn
	Scs       []string
	Cns       []Const
	Sts       []Struct
	Friends   []*Lib
}

type Arg struct {
	Str  string
	Opt  bool
	Type string
	Desc string
}

type Return struct {
	Str  string
	Type string
}

type Fn struct {
	Name    string
	Fn      string
	Args    []Arg
	Returns []Return
	Block   bool
	Desc    []string
}

type Const struct {
	Group  string
	Consts []string
}

type Prop struct {
	Str  string
	Type string
	Desc string
}

type Method struct {
	Name string
	Type string
	Desc string
}

type Struct struct {
	Struct  string
	Props   []Prop
	Methods []Method
	Desc    []string
}
