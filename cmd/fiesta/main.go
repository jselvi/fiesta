package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/jselvi/fiesta/pkg/core"
)

type fiestaPrompt struct {
	I             *readline.Instance
	CompleterRoot *readline.PrefixCompleter
	CompleterUse  *readline.PrefixCompleter

	magic  printMagic
	Attack string
	exit   bool

	c core.Core
}

func main() {
	var p fiestaPrompt
	var e error

	// Parse params
	var rcFile string
	flag.StringVar(&rcFile, "r", "", "Execute resource file (- for stdin)")
	flag.Parse()

	// Read commands from resource file (if any)
	var rcCommands []string
	if len(rcFile) > 0 {
		file, errFile := os.Open(rcFile)
		if errFile != nil {
			fmt.Printf("File %s does not exist or could not be opened\n", rcFile)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			rcCommands = append(rcCommands, scanner.Text())
		}
	}

	// Show banner
	p.Banner()
	p.magic.PrintStatus("Are you ready for the FIESTA? ;)")
	p.magic.NewLine()

	// Prepare Prompt
	p.PrepareCompleter()
	p.I, e = readline.NewEx(&readline.Config{
		AutoComplete:      p.CompleterRoot,
		HistorySearchFold: true,
	})
	if e != nil {
		panic(e)
	}
	defer p.I.Close()

	// Insert RC commands into ReadLine
	for _, s := range rcCommands {
		p.I.WriteStdin([]byte(s + "\n"))
	}

	// Run interactive prompt
	p.Run()
}
