package main

/*
Black        0;30     Dark Gray     1;30
Red          0;31     Light Red     1;31
Green        0;32     Light Green   1;32
Brown/Orange 0;33     Yellow        1;33
Blue         0;34     Light Blue    1;34
Purple       0;35     Light Purple  1;35
Cyan         0;36     Light Cyan    1;36
Light Gray   0;37     White         1;37
*/

func color(c string, s string) string {
	return "\033[" + c + "m" + s + "\033[0m"
}

func red(s string) string {
	return color("1;31", s)
}

func green(s string) string {
	return color("0;32", s)
}

func yellow(s string) string {
	return color("0;33", s)
}

func blue(s string) string {
	return color("1;34", s)
}

func cyan(s string) string {
	return color("0;36", s)
}

func (p *fiestaPrompt) setPrompt(s string) {
	var newPrompt string

	if len(s) > 0 {
		newPrompt = "fiesta (" + red(s) + ")> "
	} else {
		newPrompt = "fiesta > "
	}

	p.I.SetPrompt(newPrompt)
}
