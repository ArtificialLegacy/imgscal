package cli

type CliColor string

const (
	COLOR_RESET CliColor = "\u001b[0m"

	COLOR_BLACK   CliColor = "\u001b[30m"
	COLOR_RED     CliColor = "\u001b[31m"
	COLOR_GREEN   CliColor = "\u001b[32m"
	COLOR_YELLOW  CliColor = "\u001b[33m"
	COLOR_BLUE    CliColor = "\u001b[34m"
	COLOR_MAGENTA CliColor = "\u001b[35m"
	COLOR_CYAN    CliColor = "\u001b[36m"
	COLOR_WHITE   CliColor = "\u001b[37m"

	COLOR_BRIGHT_BLACK   CliColor = "\u001b[30;1m"
	COLOR_BRIGHT_RED     CliColor = "\u001b[31;1m"
	COLOR_BRIGHT_GREEN   CliColor = "\u001b[32;1m"
	COLOR_BRIGHT_YELLOW  CliColor = "\u001b[33;1m"
	COLOR_BRIGHT_BLUE    CliColor = "\u001b[34;1m"
	COLOR_BRIGHT_MAGENTA CliColor = "\u001b[35;1m"
	COLOR_BRIGHT_CYAN    CliColor = "\u001b[36;1m"
	COLOR_BRIGHT_WHITE   CliColor = "\u001b[37;1m"

	COLOR_BACKGROUND_BLACK   CliColor = "\u001b[40m"
	COLOR_BACKGROUND_RED     CliColor = "\u001b[41m"
	COLOR_BACKGROUND_GREEN   CliColor = "\u001b[42m"
	COLOR_BACKGROUND_YELLOW  CliColor = "\u001b[43m"
	COLOR_BACKGROUND_BLUE    CliColor = "\u001b[44m"
	COLOR_BACKGROUND_MAGENTA CliColor = "\u001b[45m"
	COLOR_BACKGROUND_CYAN    CliColor = "\u001b[46m"
	COLOR_BACKGROUND_WHITE   CliColor = "\u001b[47m"

	COLOR_BRIGHT_BACKGROUND_BLACK   CliColor = "\u001b[40;1m"
	COLOR_BRIGHT_BACKGROUND_RED     CliColor = "\u001b[41;1m"
	COLOR_BRIGHT_BACKGROUND_GREEN   CliColor = "\u001b[42;1m"
	COLOR_BRIGHT_BACKGROUND_YELLOW  CliColor = "\u001b[43;1m"
	COLOR_BRIGHT_BACKGROUND_BLUE    CliColor = "\u001b[44;1m"
	COLOR_BRIGHT_BACKGROUND_MAGENTA CliColor = "\u001b[45;1m"
	COLOR_BRIGHT_BACKGROUND_CYAN    CliColor = "\u001b[46;1m"
	COLOR_BRIGHT_BACKGROUND_WHITE   CliColor = "\u001b[47;1m"

	COLOR_BOLD      CliColor = "\u001b[1m"
	COLOR_UNDERLINE CliColor = "\u001b[4m"
	COLOR_REVERSED  CliColor = "\u001b[7m"

	COLOR_CURSOR_HOME   CliColor = "\u001b[H"
	COLOR_CURSOR_LINEUP CliColor = "\u001b M"
	COLOR_CURSOR_SAVE   CliColor = "\u001b 7"
	COLOR_CURSOR_LOAD   CliColor = "\u001b 8"

	COLOR_ERASE_DOWN   CliColor = "\u001b[0J"
	COLOR_ERASE_UP     CliColor = "\u001b[1J"
	COLOR_ERASE_SCREEN CliColor = "\u001b[2J"
	COLOR_ERASE_SAVED  CliColor = "\u001b[3J"

	COLOR_ERASE_LINE_END   CliColor = "\u001b[0K"
	COLOR_ERASE_LINE_START CliColor = "\u001b[1K"
	COLOR_ERASE_LINE       CliColor = "\u001b[2K"

	COLOR_CURSOR_INVISIBLE CliColor = "\u001b[?25l"
	COLOR_CURSOR_VISIBLE   CliColor = "\u001b[?25h"
)
