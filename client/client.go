package client

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"

	cjson "github.com/AyakuraYuki/monster-siren-downloader/component/json"
)

func init() {
	cjson.RegisterFuzzyDecoders()
}

const (
	apiHost = `https://monster-siren.hypergryph.com`
)

type Client struct {
	cli    *resty.Client
	closed atomic.Bool
}

func New(opts ...Option) *Client {
	client := &Client{
		cli: resty.New().
			SetBaseURL(apiHost).
			SetTimeout(30 * time.Second).
			SetRetryCount(2).
			SetHeaders(map[string]string{
				"Accept":          "*/*",
				"Accept-Language": "zh-CN,zh;q=0.9,ja;q=0.8,en;q=0.7,en-GB;q=0.6,en-US;q=0.5",
				"Referer":         "https://monster-siren.hypergryph.com/",
				"User-Agent":      fmt.Sprintf("Go/%s monster-siren-downloader/%s", runtime.Version(), Version()),
			}),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}
