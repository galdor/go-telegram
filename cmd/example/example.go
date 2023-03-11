package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/exograd/go-program"
	"github.com/galdor/go-netrc"
	telegrambot "github.com/galdor/go-telegram/pkg/bot"
)

const (
	DefaultUsername = "go_telegram_example_bot"
)

func main() {
	p := program.NewProgram("example",
		"an example program for the go-telegram library")

	p.AddCommand("bot-info", "print information about the bot user",
		cmdBotInfo)

	p.ParseCommandLine()
	p.Run()
}

func cmdBotInfo(p *program.Program) {
	bot := createBot(p)

	var user telegrambot.User
	if err := bot.CallMethod("getMe", nil, &user); err != nil {
		p.Fatal("cannot fetch user information: %v", err)
	}

	t := NewTable()

	t.AddRow("Id", user.Id)
	t.AddRow("Bot", user.IsBot)
	t.AddRow("Premium", user.IsPremium)
	t.AddRow("First name", user.FirstName)
	t.AddRow("Last name", user.LastName)
	t.AddRow("Username", user.Username)
	t.AddRow("Language code", user.LanguageCode)

	t.Write()
}

func createBot(p *program.Program) *telegrambot.Bot {
	token, err := netrcLookup("api.telegram.org", DefaultUsername)
	if err != nil {
		p.Fatal("cannot lookup token in netrc: %v", err)
	}

	botCfg := telegrambot.BotCfg{
		Token: token,
	}

	bot, err := telegrambot.NewBot(botCfg)
	if err != nil {
		p.Fatal("cannot create bot: %v", err)
	}

	return bot
}

func runBot(p *program.Program, bot *telegrambot.Bot) {
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
