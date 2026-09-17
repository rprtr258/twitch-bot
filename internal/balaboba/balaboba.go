package balaboba

import (
	"context"
	"net/http"
)

type Style string

const (
	Standart            = ""
	UserManual          = ""
	Recipes             = ""
	ShortStories        = ""
	WikipediaSipmlified = ""
	MovieSynopses       = ""
	FolkWisdom          = ""
	Rus                 = ""
)

var StylesByID = map[int]Style(nil)

type Client struct{}
type ClientConfig struct {
	Lang string
	HTTP *http.Client
}
type Response struct {
	BadQuery bool
	Text     string
}

func New(ClientConfig) *Client                                            { return nil }
func (*Client) Generate(context.Context, string, Style) (Response, error) { return Response{}, nil }
