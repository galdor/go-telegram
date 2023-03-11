package bot

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

const (
	DefaultEndpoint = "https://api.telegram.org"
)

type BotCfg struct {
	Endpoint string `json:"endpoint,omitempty"`
	Token    string `json:"token"`

	HTTPClient *http.Client `json:"-"`
}

type Bot struct {
	Cfg BotCfg

	endpointURI *url.URL

	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewBot(cfg BotCfg) (*Bot, error) {
	if cfg.Endpoint == "" {
		cfg.Endpoint = DefaultEndpoint
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("missing or empty token")
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	uri, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint uri: %w", err)
	}

	bot := Bot{
		Cfg: cfg,

		endpointURI: uri,

		stopChan: make(chan struct{}),
	}

	return &bot, nil
}

func (b *Bot) Start() error {
	b.wg.Add(1)
	go b.main()

	return nil
}

func (b *Bot) Stop() {
	close(b.stopChan)
	b.wg.Wait()
}

func (b *Bot) main() {
	defer b.wg.Done()

	for {
		select {
		case <-b.stopChan:
			return
		}
	}
}
