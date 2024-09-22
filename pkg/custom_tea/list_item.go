package customtea

import golua "github.com/yuin/gopher-lua"

type ListItem struct {
	title       string
	description string
	filter      string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.description }
func (i ListItem) FilterValue() string { return i.filter }

func ListItemTableFrom(state *golua.LState, i ListItem) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("title", golua.LString(i.title))
	t.RawSetString("description", golua.LString(i.description))
	t.RawSetString("filter", golua.LString(i.filter))

	return t
}

func ListItemTable(state *golua.LState, title, description, filter string) *golua.LTable {
	t := state.NewTable()

	t.RawSetString("title", golua.LString(title))
	t.RawSetString("description", golua.LString(description))
	t.RawSetString("filter", golua.LString(filter))

	return t
}

func ListItemBuild(t *golua.LTable) ListItem {
	return ListItem{
		title:       string(t.RawGetString("title").(golua.LString)),
		description: string(t.RawGetString("description").(golua.LString)),
		filter:      string(t.RawGetString("filter").(golua.LString)),
	}
}
