package libretranslate

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

const (
	FormatText = "text"
	FormatHTML = "html"
)

func (c *client) GetLanguages(ctx context.Context) ([]Language, error) {
	body, err := c.get(ctx, c.languageURL)
	if err != nil {
		return nil, err
	}
	var data []Language
	err = json.Unmarshal(body, &data)
	return data, err
}

func (c *client) GetFrontendSetting(ctx context.Context) (*FrontendSetting, error) {
	body, err := c.get(ctx, c.frontendSettingURL)
	if err != nil {
		return nil, err
	}
	var data FrontendSetting
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *client) Translate(ctx context.Context, q, source, target string) (string, error) {
	params := TranslateRequest{
		Q:      q,
		Source: source,
		Target: target,
		Format: FormatText,
		Key:    c.apiKey,
	}

	reqBody, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	body, err := c.post(ctx, reqBody, c.translateURL)
	if err != nil {
		return "", err
	}

	var translated Translated
	err = json.Unmarshal(body, &translated)
	if err != nil {
		return "", err
	}
	return translated.Text, nil
}

func (c *client) get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, errors.Wrapf(err, "problem with req url %s", url)
	}
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")

	var rReq = new(retryablehttp.Request)
	rReq.Request = req
	resp, err := c.httpClient.Do(rReq)
	if err != nil {
		return nil, errors.Wrapf(err, "problem with resp url %s", url)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *client) post(ctx context.Context, reqBody []byte, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")

	var rReq = new(retryablehttp.Request)
	rReq.Request = req
	resp, err := c.httpClient.Do(rReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
