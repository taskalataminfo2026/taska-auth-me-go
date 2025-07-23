package rusty

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/taskalataminfo2026/taska-auth-me-go/cmd/api/config"
	"github.com/taskalataminfo2026/taska-auth-me-go/cmd/api/utils"
	"github.com/taskalataminfo2026/tool-kit-lib-go/pkg/rusty"
	"github.com/taskalataminfo2026/tool-kit-lib-go/pkg/transport/http_client"

	"net/http"
	"net/url"
	"strings"
)

//go:generate mockgen -destination=../../testutils/mocks/rusty_client_mock.go -package=mocks -source=./rusty_client.go

type IRustyClient interface {
	Get(ctx context.Context, url string, headers map[string]string, queryParams map[string]string, tags []string) RustyResponse
	Post(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse
	Patch(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse
	Put(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse
	Delete(ctx context.Context, url string, headers map[string]string, params map[string]interface{}, tags []string) RustyResponse
}

type RustyClient struct {
}

type RustyResponse struct {
	Body       []byte
	StatusCode int
	Error      error
}

func (client *RustyClient) getEndpointOptions(headers map[string]string) []rusty.EndpointOption {
	var endpointOptions []rusty.EndpointOption
	for key, value := range headers {
		endpointOptions = append(endpointOptions, rusty.WithHeader(key, value))
	}
	return endpointOptions
}

func (client *RustyClient) getRequestOptions(params map[string]interface{}) []rusty.RequestOption {
	var requestOptions []rusty.RequestOption
	for key, value := range params {
		requestOptions = append(requestOptions, rusty.WithParam(key, value))
	}
	return requestOptions
}

func (client *RustyClient) getQueryParamsOptions(queryParams map[string]string) url.Values {
	query := url.Values{}
	for key, value := range queryParams {
		query.Add(key, value)
	}
	return query
}

func (client *RustyClient) generateResponse(ctx context.Context, response *rusty.Response, err error, tags []string) RustyResponse {
	rustyResponse := RustyResponse{
		Error: err,
	}
	if err != nil {
		log.Error(ctx, fmt.Sprintf("API ERROR¡: %v, error: %v", err.Error(), err))
	}
	if response != nil {
		rustyResponse.Body = response.Body
		rustyResponse.StatusCode = response.StatusCode
		tags = utils.Merge(tags, fmt.Sprintf("status code: %v, resp:%v", response.StatusCode, strings.TrimSpace(string(response.Body))))
		log.Info(ctx, "response")
	}
	if rustyResponse.StatusCode == 0 {
		rustyResponse.StatusCode = http.StatusFailedDependency
	}

	return rustyResponse
}

func (client *RustyClient) Get(ctx context.Context, url string, headers map[string]string, queryParams map[string]string, tags []string) RustyResponse {
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), "method:GET", fmt.Sprintf("headers:%v", headers), fmt.Sprintf("queryParams:%v", queryParams))
	requester := http_client.NewRetryable(
		config.RustyConfig.RetryCount,
		http_client.WithTimeout(config.RustyConfig.DefaultTimeOut),
	)

	log.Info(ctx, fmt.Sprintf("GET - url: %v, headers: %v", url, headers))
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), "method:GET")

	endpoint, err := rusty.NewEndpoint(requester, url, client.getEndpointOptions(headers)...)
	if err != nil {
		return RustyResponse{
			Error: err,
		}
	}

	response, err := endpoint.Get(ctx, rusty.WithQuery(client.getQueryParamsOptions(queryParams)))
	return client.generateResponse(ctx, response, err, tags)
}

func (client *RustyClient) Post(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse {
	headers["Content-type"] = "application/json"
	bodyJSON, _ := json.Marshal(&body)
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), "method:POST", fmt.Sprintf("headers:%v", headers), fmt.Sprintf("body:%v", string(bodyJSON)))
	requester := http_client.NewRetryable(
		0,
		http_client.WithTimeout(config.RustyConfig.DefaultTimeOut),
	)

	log.Info(ctx, fmt.Sprintf("POST - url: %v, headers: %v, body %v", url, headers, string(bodyJSON)))
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), fmt.Sprintf("body:%v", string(bodyJSON)), "method:POST")

	endpoint, err := rusty.NewEndpoint(requester, url, client.getEndpointOptions(headers)...)
	if err != nil {
		return RustyResponse{
			Error: err,
		}
	}

	response, err := endpoint.Post(ctx, rusty.WithBody(body))
	return client.generateResponse(ctx, response, err, tags)
}

func (client *RustyClient) Patch(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse {
	headers["Content-type"] = "application/json"
	bodyJSON, _ := json.Marshal(&body)
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), "method:PATCH", fmt.Sprintf("headers:%v", headers), fmt.Sprintf("body:%v", string(bodyJSON)))
	requester := http_client.NewRetryable(
		0,
		http_client.WithTimeout(config.RustyConfig.DefaultTimeOut),
	)

	log.Info(ctx, fmt.Sprintf("PATCH - url: %v, headers: %v, body %v", url, headers, string(bodyJSON)))
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), fmt.Sprintf("body:%v", string(bodyJSON)), "method:PATCH")

	endpoint, err := rusty.NewEndpoint(requester, url, client.getEndpointOptions(headers)...)
	if err != nil {
		return RustyResponse{
			Error: err,
		}
	}

	response, err := endpoint.Patch(ctx, rusty.WithBody(body))
	return client.generateResponse(ctx, response, err, tags)
}

func (client *RustyClient) Put(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse {
	headers["Content-type"] = "application/json"
	bodyJSON, _ := json.Marshal(&body)
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), "method:PUT", fmt.Sprintf("headers:%v", headers), fmt.Sprintf("body:%v", string(bodyJSON)))
	requester := http_client.NewRetryable(
		0,
		http_client.WithTimeout(config.RustyConfig.DefaultTimeOut),
	)

	log.Info(ctx, fmt.Sprintf("PUT - url: %v, headers: %v, body %v", url, headers, string(bodyJSON)))
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), fmt.Sprintf("body:%v", string(bodyJSON)), "method:PUT")

	endpoint, err := rusty.NewEndpoint(requester, url, client.getEndpointOptions(headers)...)
	if err != nil {
		return RustyResponse{
			Error: err,
		}
	}

	response, err := endpoint.Put(ctx, rusty.WithBody(body))
	return client.generateResponse(ctx, response, err, tags)
}

func (client *RustyClient) Delete(ctx context.Context, url string, headers map[string]string, params map[string]interface{}, tags []string) RustyResponse {
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), "method:DELETE", fmt.Sprintf("headers:%v", headers), fmt.Sprintf("params:%v", params))
	requester := http_client.NewRetryable(
		0,
		http_client.WithTimeout(config.RustyConfig.DefaultTimeOut),
	)

	log.Info(ctx, fmt.Sprintf("DELETE - url: %v, headers: %v", url, headers))
	tags = utils.Merge(tags, fmt.Sprintf("url:%v", url), "method:DELETE")

	endpoint, err := rusty.NewEndpoint(requester, url, client.getEndpointOptions(headers)...)
	if err != nil {
		return RustyResponse{
			Error: err,
		}
	}

	response, err := endpoint.Delete(ctx, client.getRequestOptions(params)...)
	return client.generateResponse(ctx, response, err, tags)
}
