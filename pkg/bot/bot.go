package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
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

func (b *Bot) MethodURI(method string) *url.URL {
	base := *b.endpointURI

	ref := url.URL{
		Path: path.Join("bot"+b.Cfg.Token, method),
	}

	return base.ResolveReference(&ref)
}

func (b *Bot) CallMethod(method string, body, result interface{}) error {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("cannot encode body data: %w", err)
	}

	uri := b.MethodURI(method)

	req, err := http.NewRequest("POST", uri.String(),
		bytes.NewReader(bodyData))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := b.Cfg.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	if err := DecodeResponse(data, result); err != nil {
		return err
	}

	return nil
}
