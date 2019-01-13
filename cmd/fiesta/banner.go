package main

import "fmt"

func (p *fiestaPrompt) Banner() {
	shotBanner()
}

func cutreshotBanner() {
	shots := "                                                  _____\n" +
		"    |~~~~~|        |     |        |     |        |     |\n" +
		"    |     |        |~~~~~|        |     |        |     |\n" +
		"    |     |        |     |        |~~~~~|        |     |\n" +
		"    |_____|        |_____|        |_____|        |     |\n\n"

	var s string
	for _, c := range shots {
		cc := string(c)
		if cc == "~" {
			s = s + yellow(cc)
		} else {
			s = s + cyan(cc)
		}
	}

	fmt.Print(s)
}

func shotBanner() {
	shots := "                                                           _____    \n" +
		"      \\~~~~~~~/        \\       /        \\       /           ) (     \n" +
		"       \\     /          \\~~~~~/          \\     /            )_(     \n" +
		"        \\ _ /            \\ _ /            \\~~~/            /   \\    \n" +
		"         ) (              ) (              ) (            /     \\   \n" +
		"        _)_(_            _)_(_            _)_(_          /       \\  \n\n "

	var s string
	for _, c := range shots {
		cc := string(c)
		if cc == "~" {
			s = s + yellow(cc)
		} else {
			s = s + cyan(cc)
		}
	}

	fmt.Print(s)
}
