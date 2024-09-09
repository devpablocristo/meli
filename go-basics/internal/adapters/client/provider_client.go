package client

import (
	"net/http"
	"time"

	"github.com/osalomon89/go-basics/internal/core/domain"
	"github.com/osalomon89/go-basics/internal/core/ports"
)

type providerClient struct {
	httpClient *http.Client
}

func NewProviderClient() ports.ProviderClient {
	t := http.Transport{
		IdleConnTimeout:     5 * time.Second,
		MaxConnsPerHost:     100,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}

	httpClient := &http.Client{
		Timeout:   3 * time.Second,
		Transport: &t,
	}

	return &providerClient{
		httpClient: httpClient,
	}
}

func (c *providerClient) GetProvider(id int) (domain.Provider, error) {
	var result domain.Provider

	// endpoint, err := rusty.NewEndpoint(c.httpClient, "http://localhost:8000/items/{id}")
	// if err != nil {
	// 	return result, errors.New("error")
	// }

	// resp, err := endpoint.Get(context.Background(), rusty.WithParam("id", id))
	// if err != nil {
	// 	return result, errors.New("error")
	// }

	// if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
	// 	errorMessage := fmt.Sprintf("unexpected response status %v", resp.StatusCode)
	// 	errorDetail := fmt.Sprintf("error getting items from API, %s", errorMessage)

	// 	fmt.Println(errorDetail)

	// 	return result, errors.New("error")
	// }

	// err = json.Unmarshal(resp.Body, &result)
	// if err != nil {
	// 	return result, errors.New("error")
	// }

	return result, nil
}
