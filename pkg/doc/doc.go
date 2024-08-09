package doc

const (
	TAG_EMPTY     = "/// "
	TAG_EXISTS    = "/// @"
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
)

type Lib struct {
	Name string
	Fns  []Fn
	Cns  []Const
	Sts  []Struct
}

type Arg struct {
	Str string
	Opt bool
}

type Fn struct {
	Fn      string
	Args    []Arg
	Returns []string
	Block   bool
	Desc    []string
}

type Const struct {
	Group  string
	Consts []string
}

type Struct struct {
	Struct  string
	Props   []string
	Methods []string
	Desc    []string
}
