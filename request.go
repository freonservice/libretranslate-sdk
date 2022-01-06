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
	req, err := c.getRequest(ctx, c.languageURL)
	if err != nil {
		return nil, err
	}

	body, err := c.parseBody(ctx, req)
	if err != nil {
		return nil, err
	}

	var data []Language
	err = json.Unmarshal(body, &data)
	return data, err
}

func (c *client) GetFrontendSetting(ctx context.Context) (*FrontendSetting, error) {
	req, err := c.getRequest(ctx, c.frontendSettingURL)
	if err != nil {
		return nil, err
	}

	body, err := c.parseBody(ctx, req)
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

	req, err := c.postRequest(ctx, reqBody, c.translateURL)
	if err != nil {
		return "", err
	}

	body, err := c.parseBody(ctx, req)
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

func (c *client) getRequest(ctx context.Context, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
}

func (c *client) postRequest(ctx context.Context, body []byte, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
}

func (c *client) parseBody(ctx context.Context, req *http.Request) ([]byte, error) {
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")

	var rReq = new(retryablehttp.Request)
	rReq.Request = req
	resp, err := c.httpClient.Do(rReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return body, nil
	}

	var errorMsg ErrorMsg
	err = json.Unmarshal(body, &errorMsg)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(errorMsg.Error)
}
