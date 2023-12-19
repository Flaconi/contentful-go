package cma

import (
	"context"
	"fmt"
	"github.com/flaconi/contentful-go/internal/cma/api_keys"
	"github.com/flaconi/contentful-go/internal/cma/preview_api_keys"
	"github.com/flaconi/contentful-go/service/cma"
	"github.com/flaconi/contentful-go/service/common"
	"io"
	"net/http"
	"net/url"
)

var _ common.SpaceIdClient = &SpaceIdClient{}

type SpaceIdClient struct {
	client  *Client
	spaceId string
}

func (c *SpaceIdClient) Get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Get(ctx, fmt.Sprintf("/spaces/%s%s", c.spaceId, path), queryParams, headers)
}

func (c *SpaceIdClient) Post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Post(ctx, fmt.Sprintf("/spaces/%s%s", c.spaceId, path), queryParams, headers, body)
}

func (c *SpaceIdClient) Put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.client.Put(ctx, fmt.Sprintf("/spaces/%s%s", c.spaceId, path), queryParams, headers, body)
}

func (c *SpaceIdClient) Delete(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.client.Delete(ctx, fmt.Sprintf("/spaces/%s%s", c.spaceId, path), queryParams, headers)
}

func (c *SpaceIdClient) WithEnvironment(environment string) common.EnvironmentClient {
	return &EnvironmentClient{
		Client:      c.client,
		Environment: environment,
	}
}

func (c *SpaceIdClient) ApiKeys() cma.ApiKeys {
	return api_keys.NewApiKeysService(c)
}

func (c *SpaceIdClient) PreviewApiKeys() cma.PreviewApiKeys {
	return preview_api_keys.NewPreviewApiKeysService(c)
}
