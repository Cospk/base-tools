package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Cospk/base-tools/errs"
	"io"
	"net/http"
	"time"
)

// ClientConfig 定义 HTTP 客户端的配置。
type ClientConfig struct {
	Timeout         time.Duration
	MaxConnsPerHost int
}

// NewClientConfig 创建默认的客户端配置。
func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:         15 * time.Second,
		MaxConnsPerHost: 100,
	}
}

// HTTPClient 封装 http.Client 并包含额外的配置。
type HTTPClient struct {
	client *http.Client
	config *ClientConfig
}

// NewHTTPClient 使用提供的配置创建新的 HTTPClient。
func NewHTTPClient(config *ClientConfig) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				MaxConnsPerHost: config.MaxConnsPerHost,
			},
		},
		config: config,
	}
}

// NewHTTPClientWithClient 使用提供的配置和 http.Client 创建新的 HTTPClient。
func NewHTTPClientWithClient(client *http.Client, config *ClientConfig) *HTTPClient {
	return &HTTPClient{
		client: client,
		config: config,
	}
}

// Get 执行 HTTP GET 请求并返回响应体。
func (c *HTTPClient) Get(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, errs.WrapMsg(err, "GET request failed", "url", url)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.WrapMsg(err, "failed to read response body", "url", url)
	}
	return body, nil
}

// Post 发送 JSON 编码的 POST 请求并返回响应体。
func (c *HTTPClient) Post(ctx context.Context, url string, headers map[string]string, data any, timeout int) ([]byte, error) {
	if timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}
	body := bytes.NewBuffer(nil)
	if data != nil {
		if err := json.NewEncoder(body).Encode(data); err != nil {
			return nil, errs.WrapMsg(err, "JSON encode failed", "data", data)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, errs.WrapMsg(err, "NewRequestWithContext failed", "url", url, "method", http.MethodPost)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errs.WrapMsg(err, "HTTP request failed")
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.WrapMsg(err, "failed to read response body")
	}

	return result, nil
}

// PostReturn sends a JSON-encoded POST request and decodes the JSON response into  output parameter.
func (c *HTTPClient) PostReturn(ctx context.Context, url string, headers map[string]string, input, output any, timeout int) error {
	responseBytes, err := c.Post(ctx, url, headers, input, timeout)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(responseBytes, output); err != nil {
		return errs.WrapMsg(err, "JSON unmarshal failed")
	}
	return nil
}
