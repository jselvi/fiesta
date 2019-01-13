package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
)

func (p *fiestaPrompt) Executor(line string) {

	line = strings.TrimSpace(line)
	line = strings.ToLower(line)
	vLine := strings.Split(line, " ")

	switch {

	// EXIT
	case line == "exit" || line == "quit":
		p.exit = true

	// USE
	case strings.HasPrefix(line, "use "):
		if len(p.Attack) > 0 { // ignore when attack was already selected
			return
		}

		attack := line[4:]
		for _, v := range p.listAttacks() {
			if v == attack {
				// Prompt
				p.Attack = attack
				p.setPrompt(p.Attack)
				p.I.Config.AutoComplete = p.CompleterUse
				// Core
				p.c.SetAttack(attack)
				return
			}
		}
		fmt.Println(p.Usage("use"))

	// BACK
	case line == "back":
		if len(p.Attack) == 0 { // ignore when no attack selected
			return
		}

		p.Attack = ""
		p.setPrompt(p.Attack)
		p.I.Config.AutoComplete = p.CompleterRoot

	// SHOW
	case line == "show attacks":
		attacks := p.listAttacks()
		for _, v := range attacks {
			fmt.Println(v)
		}

	case line == "show options":
		if len(p.Attack) == 0 { // not shown when no attack selected
			return
		}
		fmt.Printf("\nModule options (%s):\n\n", p.Attack)
		fmt.Println("\tName         Current Setting   Description")
		fmt.Println("\t----         ---------------   -----------")
		for _, v := range p.c.Options() {
			fmt.Printf("\t%-10s   %-15s   %s\n", v.Param, v.Value, v.Descr)
		}
		fmt.Printf("\n")

	case strings.HasPrefix(line, "show "):
		fmt.Println(p.Usage("show"))

	// SET / UNSET
	case strings.HasPrefix(line, "set "):
		if len(p.Attack) == 0 { // ignore when no attack selected
			return
		}

		if len(vLine) < 3 {
			fmt.Println(p.Usage("set"))
			return
		}
		option := strings.ToUpper(vLine[1])
		n := len(vLine[0]) + len(vLine[1]) + 2
		value := line[n:]
		p.c.SetOption(option, value)
		fmt.Printf("%s = %s", option, value)

	case strings.HasPrefix(line, "unset "):
		if len(p.Attack) == 0 { // ignore when no attack selected
			return
		}

		if len(vLine) < 2 {
			fmt.Println(p.Usage("unset"))
			return
		}
		option := strings.ToUpper(vLine[1])
		p.c.UnsetOption(option)

	// EXPLOIT
	case line == "exploit" || line == "run":
		if len(p.Attack) == 0 { // ignore when no attack selected
			return
		}

		statusCh, msgCh, resCh, errorCh := p.c.Exploit()

		// close channel if CTRL+C
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go p.stopIfCtrlC(sigs)

		// Exit if we receive an error
		go func() {
			e, ok := <-errorCh
			p.c.Break()
			if ok {
				p.magic.NewLine()
				p.magic.PrintError(e)
			}
		}()

		// Print status if received
		go func() {
			for {
				s, ok := <-statusCh
				if !ok {
					break
				}
				p.magic.PrintStatus(s)
			}
		}()

		// Print result if received
		go func() {
			for {
				s, ok := <-resCh
				if !ok {
					break
				}
				p.magic.PrintResult(s)
			}
		}()

		// print results until CTRL+C or finish
		for m, ok := <-msgCh; ok; m, ok = <-msgCh {
			p.magic.Print(m)
		}
		p.magic.NewLine()

	// NO COMMAND
	case line == "":
		return

	// USAGE
	case len(p.Usage(vLine[0])) > 0:
		fmt.Println(p.Usage(vLine[0]))

		/*
			default: // execute OS command
				out, _ := exec.Command("sh", "-c", line).Output()
				if len(out) > 9 {
					fmt.Println(string(out))
				}
		*/

	}

}

func (p *fiestaPrompt) stopIfCtrlC(signal chan os.Signal) {
	<-signal
	p.c.Break()
}

func (p *fiestaPrompt) Run() {
	// set default Prompt
	p.setPrompt("")

	log.SetOutput(p.I.Stderr())
	p.exit = false
	for !p.exit {
		line, e := p.I.Readline()
		if e == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if e == io.EOF {
			break
		}

		p.Executor(line)
	}
}

func (p *fiestaPrompt) Usage(cmd string) string {
	var usage string

	switch cmd {
	case "use":
		usage = "Usage: use [attack]"
	case "show":
		usage = "Usage: show [attacks|options]"
	case "set":
		usage = "Usage: set [option] [value]"
	case "unset":
		usage = "Usage: unset [option]"
	default:
		usage = ""
	}

	return usage
}
