package doc

const (
	TAG_FUNC      = "/// @func "
	TAG_ARG       = "/// @arg "
	TAG_RETURNS   = "/// @returns "
	TAG_CONSTANTS = "/// @constants "
	TAG_CONST     = "/// @const "
	TAG_BLOCK     = "/// @blocking"
)

type Lib struct {
	Name string
	Fns  []Fn
	Cns  []Const
}

type Fn struct {
	Fn      string
	Args    []string
	Returns []string
	Block   bool
}

type Const struct {
	Group  string
	Consts []string
}
