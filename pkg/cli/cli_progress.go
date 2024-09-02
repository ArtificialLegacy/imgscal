package cli

import "fmt"

func Progress(current, total, width int, title, bar, empty string, reset bool) {
	percent := float64(current%total) / float64(total)
	if current > 0 && percent == 0 {
		percent = 1
	}

	count := int(percent * float64(width))

	color := COLOR_WHITE
	if percent == 1 {
		color = COLOR_BRIGHT_GREEN + COLOR_BOLD
	} else if percent > 0.9 {
		color = COLOR_GREEN
	} else if percent > 0.5 {
		color = COLOR_CYAN
	} else if percent > 0.25 {
		color = COLOR_YELLOW
	} else if percent > 0 {
		color = COLOR_RED
	} else {
		color = COLOR_RED + COLOR_BOLD
	}

	progress := fmt.Sprintf("  %s [%s", title, color)
	for i := 0; i < count; i++ {
		progress += bar
	}

	for i := 0; i < width-count; i++ {
		progress += empty
	}

	progress += fmt.Sprintf("%s] %d/%d (%s%d%%%s)", COLOR_RESET, current, total, color, int(percent*100), COLOR_RESET)

	if reset {
		if count == width {
			progress += "\n"
		} else {
			progress += fmt.Sprintf("\u001b[%dD", len(progress))
		}
	} else {
		progress += "\n"
	}

	fmt.Print(progress)
}
