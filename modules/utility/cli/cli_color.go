package cli

type CliColor string

const (
	RESET CliColor = "\u001b[0m"

	BLACK   CliColor = "\u001b[30m"
	RED     CliColor = "\u001b[31m"
	GREEN   CliColor = "\u001b[32m"
	YELLOW  CliColor = "\u001b[33m"
	BLUE    CliColor = "\u001b[34m"
	MAGENTA CliColor = "\u001b[35m"
	CYAN    CliColor = "\u001b[36m"
	WHITE   CliColor = "\u001b[37m"

	BRIGHT_BLACK   CliColor = "\u001b[30;1m"
	BRIGHT_RED     CliColor = "\u001b[31;1m"
	BRIGHT_GREEN   CliColor = "\u001b[32;1m"
	BRIGHT_YELLOW  CliColor = "\u001b[33;1m"
	BRIGHT_BLUE    CliColor = "\u001b[34;1m"
	BRIGHT_MAGENTA CliColor = "\u001b[35;1m"
	BRIGHT_CYAN    CliColor = "\u001b[36;1m"
	BRIGHT_WHITE   CliColor = "\u001b[37;1m"

	BACKGROUND_BLACK   CliColor = "\u001b[40m"
	BACKGROUND_RED     CliColor = "\u001b[41m"
	BACKGROUND_GREEN   CliColor = "\u001b[42m"
	BACKGROUND_YELLOW  CliColor = "\u001b[43m"
	BACKGROUND_BLUE    CliColor = "\u001b[44m"
	BACKGROUND_MAGENTA CliColor = "\u001b[45m"
	BACKGROUND_CYAN    CliColor = "\u001b[46m"
	BACKGROUND_WHITE   CliColor = "\u001b[47m"

	BRIGHT_BACKGROUND_BLACK   CliColor = "\u001b[40;1m"
	BRIGHT_BACKGROUND_RED     CliColor = "\u001b[41;1m"
	BRIGHT_BACKGROUND_GREEN   CliColor = "\u001b[42;1m"
	BRIGHT_BACKGROUND_YELLOW  CliColor = "\u001b[43;1m"
	BRIGHT_BACKGROUND_BLUE    CliColor = "\u001b[44;1m"
	BRIGHT_BACKGROUND_MAGENTA CliColor = "\u001b[45;1m"
	BRIGHT_BACKGROUND_CYAN    CliColor = "\u001b[46;1m"
	BRIGHT_BACKGROUND_WHITE   CliColor = "\u001b[47;1m"

	BOLD      CliColor = "\u001b[1m"
	UNDERLINE CliColor = "\u001b[4m"
	REVERSED  CliColor = "\u001b[7m"
)
