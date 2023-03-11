package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/exograd/go-program"
	"github.com/galdor/go-netrc"
	"github.com/galdor/go-telegram/pkg/bot"
)

func main() {
	p := program.NewProgram("example",
		"an example program for the go-telegram library")

	p.AddArgument("username", "the telegram username of the bot")

	p.SetMain(exampleMain)

	p.ParseCommandLine()
	p.Run()
}

func exampleMain(p *program.Program) {
	username := p.ArgumentValue("username")

	token, err := netrcLookup("api.telegram.org", username)
	if err != nil {
		p.Fatal("cannot lookup token in netrc: %v", err)
	}

	botCfg := bot.BotCfg{
		Token: token,
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

func netrcLookup(host, login string) (string, error) {
	var es netrc.Entries
	if err := es.Load(netrc.DefaultPath()); err != nil {
		return "", fmt.Errorf("cannot load netrc entries: %w", err)
	}

	s := netrc.Search{
		Machine: host,
		Login:   login,
	}

	matches := es.Search(s)
	if len(matches) == 0 {
		return "", nil
	}

	return matches[0].Password, nil
}
