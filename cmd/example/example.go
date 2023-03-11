package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/exograd/go-program"
	"github.com/galdor/go-telegram/pkg/bot"
)

func main() {
	p := program.NewProgram("example",
		"an example program for the go-telegram library")

	p.SetMain(exampleMain)

	p.ParseCommandLine()
	p.Run()
}

func exampleMain(p *program.Program) {
	botCfg := bot.BotCfg{
		Token: "", // TODO
	}

	bot, err := bot.NewBot(botCfg)
	if err != nil {
		p.Fatal("cannot create bot: %v", err)
	}

	if err := bot.Start(); err != nil {
		p.Fatal("cannot start bot: %v", err)
	}

	defer bot.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case signo := <-sigChan:
		p.Info("\nreceived signal %d (%v)", signo, signo)
	}
}
