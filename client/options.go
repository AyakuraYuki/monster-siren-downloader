package client

import "time"

type Option func(*Client)

func WithRetryCount(retryCount int) Option {
	return func(client *Client) {
		if retryCount >= 0 {
			client.cli.SetRetryCount(retryCount)
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(client *Client) {
		if timeout >= 0 {
			client.cli.SetTimeout(timeout)
		}
	}
}

func WithUserAgent(userAgent string) Option {
	return func(client *Client) {
		client.cli.SetHeader("User-Agent", userAgent)
	}
}
