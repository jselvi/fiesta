package main

import (
	"fmt"
	"strings"
)

type printMagic struct {
	prev string
}

func (p *printMagic) isNextLetter(str string) bool {
	if str == p.prev {
		return true
	}

	if len(str)-len(p.prev) != 1 {
		return false
	}
	if strings.Contains(str, p.prev) {
		return true
	}
	return false
}

func (p *printMagic) NewLine() {
	p.prev = ""
	fmt.Printf("\n")
}

func (p *printMagic) Print(str string) {
	if p.isNextLetter(str) {
		p.prev = str
		str = "\r" + str
	} else {
		p.prev = str
		str = "\n" + str
	}
	fmt.Printf(str)
}

func (p *printMagic) Printf(str string, param ...string) {
	if p.isNextLetter(str) {
		str = "\r" + str
	} else {
		p.prev = str
		str = "\n" + str
	}
	fmt.Printf(str, param)
}

func (p *printMagic) PrintStatus(msg string) {
	s := blue("[*]") + " " + msg
	fmt.Println(s)
}

func (p *printMagic) PrintResult(msg string) {
	s := "\n" + green("[+]") + " " + msg
	fmt.Println(s)
}

func (p *printMagic) PrintError(e error) {
	s := red("[-]") + " " + e.Error()
	fmt.Println(s)
}
