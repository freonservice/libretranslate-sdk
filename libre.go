package libretranslate

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/powerman/structlog"
)

const (
	defaultRetryMax    = 5
	defaultConnTimeout = 15 * time.Second
)

type Client interface {
	SetRetryMax(retryMax int) Client
	SetConnTimeout(connTimeout time.Duration) Client
	SetURL(apiURL string) Client
	SetKey(apiKey string) Client
	GetLanguages(ctx context.Context) ([]Language, error)
	GetFrontendSetting(ctx context.Context) (*FrontendSetting, error)
	Translate(ctx context.Context, q, source, target string) (string, error)
}

type client struct {
	retryMax int
	apiURL   string
	apiKey   string

	languageURL        string
	frontendSettingURL string
	translateURL       string

	connTimeout time.Duration

	httpClient *retryablehttp.Client
	logger     *structlog.Logger
}

func NewLibreTranslate(apiURL string) Client {
	c := &client{
		apiURL:             apiURL,
		retryMax:           defaultRetryMax,
		connTimeout:        defaultConnTimeout,
		languageURL:        fmt.Sprintf("%s/languages", apiURL),
		frontendSettingURL: fmt.Sprintf("%s/frontend/settings", apiURL),
		translateURL:       fmt.Sprintf("%s/translate", apiURL),
		logger:             structlog.New(),
	}
	c.initClient()
	return c
}

func (c *client) SetRetryMax(retryMax int) Client {
	c.retryMax = retryMax
	return c
}

func (c *client) SetConnTimeout(connTimeout time.Duration) Client {
	c.connTimeout = connTimeout
	return c
}

func (c *client) SetURL(apiURL string) Client {
	c.apiURL = apiURL
	return c
}

func (c *client) SetKey(apiKey string) Client {
	c.apiKey = apiKey
	return c
}

func (c *client) initClient() {
	client := &http.Client{
		Timeout: c.connTimeout,
		Transport: &http.Transport{
			DialContext:         (&net.Dialer{}).DialContext,
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
	}

	httpClient := &retryablehttp.Client{
		RetryMax:   c.retryMax,
		Logger:     c.logger,
		HTTPClient: client,
		Backoff:    retryablehttp.LinearJitterBackoff,
		CheckRetry: retryablehttp.DefaultRetryPolicy,
		RequestLogHook: func(logger retryablehttp.Logger, request *http.Request, i int) {
			if i > 0 {
				logger.Printf("retry url %s attempt %d", request.URL.Path, i)
			}
		},
	}
	c.httpClient = httpClient
}
