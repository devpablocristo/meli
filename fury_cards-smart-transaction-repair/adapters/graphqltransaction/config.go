package graphqltransaction

import "net/url"

type Config struct {
	BaseURL url.URL
}

func NewConfig(baseURL *url.URL) *Config {
	return &Config{
		BaseURL: *baseURL,
	}
}
