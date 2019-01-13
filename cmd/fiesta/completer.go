package main

import (
	"../../pkg/core"
	"github.com/chzyer/readline"
)

func (p *fiestaPrompt) listAttacksFunc() func(string) []string {
	return func(line string) []string {
		return p.listAttacks()
	}
}

func (p *fiestaPrompt) listAttacks() []string {
	return core.Attacks.DumpCmd()
}

func (p *fiestaPrompt) listOptionsFunc() func(string) []string {
	return func(line string) []string {
		return p.listOptions()
	}
}

func (p *fiestaPrompt) listOptions() []string {
	var res []string
	for _, v := range p.c.Options() {
		res = append(res, v.Param)
	}
	return res
}

func (p *fiestaPrompt) PrepareCompleter() {
	p.CompleterRoot = readline.NewPrefixCompleter(
		readline.PcItem("use", readline.PcItemDynamic(p.listAttacksFunc())),
		readline.PcItem("show", readline.PcItem("attacks")),
		readline.PcItem("clear"),
		readline.PcItem("exit"),
	)

	p.CompleterUse = readline.NewPrefixCompleter(
		readline.PcItem("show", readline.PcItem("options")),
		readline.PcItem("set", readline.PcItemDynamic(p.listOptionsFunc())),
		readline.PcItem("unset", readline.PcItemDynamic(p.listOptionsFunc())),
		readline.PcItem("exploit"),
		readline.PcItem("clear"),
		readline.PcItem("back"),
		readline.PcItem("exit"),
	)
}
