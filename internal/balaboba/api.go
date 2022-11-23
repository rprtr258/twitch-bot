package balaboba

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const apiurl = "https://yandex.ru/lab/api/yalm/"

// MinTimeout is a minimum time limit for api requests.
const MinTimeout = 20 * time.Second

// New makes new balaboba api client.
//
// If the timeout is not specified or it is less than MinTimeout
// it will be equal to MinTimeout.
// Anyway the request can be canceled via the context.
func New(lang Lang, httpClient http.Client, timeout ...time.Duration) *Client {
	return &Client{
		lang:       lang,
		httpClient: httpClient,
	}
}

// Client is Yandex Balaboba client.
type Client struct {
	httpClient http.Client
	lang       Lang
}

type responseBase struct {
	Error int `json:"error"`
}

func (r responseBase) err() int { return r.Error }

type errorable interface{ err() int }

func (c *Client) do(endpoint string, data interface{}, dst errorable) error {
	return c.doContext(context.Background(), endpoint, data, dst)
}

func (c *Client) doContext(ctx context.Context, endpoint string, data interface{}, dst errorable) error {
	err := c.request(ctx, apiurl+endpoint, data, dst)
	if err != nil {
		return err
	}
	if c := dst.err(); c != 0 {
		err = fmt.Errorf("balaboba: error code %d", c)
	}
	return err
}

func (c *Client) request(ctx context.Context, url string, data, dst interface{}) error {
	method := http.MethodGet
	var body io.Reader

	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
		method = http.MethodPost
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("balaboba: response status %s (%d)", resp.Status, resp.StatusCode)
	}

	if dst == nil {
		return nil
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(dst); err != nil {
		raw, _ := io.ReadAll(io.MultiReader(dec.Buffered(), resp.Body))
		err = fmt.Errorf("balaboba: %s\nresponse: %s", err.Error(), string(raw))
	}
	return err
}

// Lang represents balaboba language.
type Lang uint8

// available languages.
const (
	Rus Lang = iota
	Eng
)
