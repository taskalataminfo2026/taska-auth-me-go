package clients

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"net/url"
//	"strings"
//
//	"github.com/labstack/gommon/log"
//)
//
//// IRustyClient defines the interface for making HTTP requests with Rusty.
//type IRustyClient interface {
//	Get(ctx context.Context, url string, headers map[string]string, queryParams map[string]string, tags []string) RustyResponse
//	Post(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse
//	Patch(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse
//	Put(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse
//	Delete(ctx context.Context, url string, headers map[string]string, params map[string]interface{}, tags []string) RustyResponse
//}
//
//// RustyClient is the implementation of IRustyClient
//type RustyClient struct {
//	baseURL    string
//	headers    map[string]string
//	httpClient *HttpClient
//}
//
//// RustyResponse represents the response from an HTTP request
//type RustyResponse struct {
//	Body       []byte
//	StatusCode int
//	Error      error
//}
//
//// NewRustyClient creates a new RustyClient instance
//func NewRustyClient(baseURL string, opts ...Option) *RustyClient {
//	client := New(baseURL, opts...)
//	return &RustyClient{
//		baseURL:    baseURL,
//		headers:    make(map[string]string),
//		httpClient: client,
//	}
//}
//
//// WithDefaultHeaders sets default headers for all requests
//func (c *RustyClient) WithDefaultHeaders(headers map[string]string) *RustyClient {
//	for k, v := range headers {
//		c.headers[k] = v
//	}
//	return c
//}
//
//func (c *RustyClient) addHeaders(builder *RequestBuilder, headers map[string]string) *RequestBuilder {
//	// Add default headers first
//	for k, v := range c.headers {
//		builder = builder.WithHeader(k, v)
//	}
//	// Add/override with request-specific headers
//	for k, v := range headers {
//		builder = builder.WithHeader(k, v)
//	}
//	return builder
//}
//
//func (c *RustyClient) addQueryParams(builder *RequestBuilder, queryParams map[string]string) *RequestBuilder {
//	if len(queryParams) == 0 {
//		return builder
//	}
//
//	params := url.Values{}
//	for k, v := range queryParams {
//		params.Add(k, v)
//	}
//
//	// Get the current path and append query params
//	path := builder.path
//	if strings.Contains(path, "?") {
//		path += "&" + params.Encode()
//	} else {
//		path += "?" + params.Encode()
//	}
//
//	// Create a new builder with the updated path
//	return &RequestBuilder{
//		client:  builder.client,
//		method:  builder.method,
//		path:    path,
//		header:  builder.header,
//		body:    builder.body,
//		context: builder.context,
//	}
//}
//
//func (c *RustyClient) logRequest(ctx context.Context, method, url string, headers map[string]string, body interface{}, tags []string) {
//	var bodyStr string
//	if body != nil {
//		if b, ok := body.([]byte); ok {
//			bodyStr = string(b)
//		} else if b, err := json.Marshal(body); err == nil {
//			bodyStr = string(b)
//		}
//	}
//
//	tags = append(tags,
//		fmt.Sprintf("method:%s", method),
//		fmt.Sprintf("url:%s", url),
//		fmt.Sprintf("headers:%v", headers),
//	)
//
//	if bodyStr != "" {
//		tags = append(tags, fmt.Sprintf("body:%s", bodyStr))
//	}
//
//	log.Info(ctx, fmt.Sprintf("%s - %s", method, url), buildFields(tags)...)
//}
//
//func (c *RustyClient) handleResponse(ctx context.Context, resp *http.Response, err error, tags []string) RustyResponse {
//	result := RustyResponse{
//		Error:      err,
//		StatusCode: http.StatusFailedDependency, // Default to 424 if no status code is set
//	}
//
//	if err != nil {
//		log.Error(ctx, fmt.Sprintf("API ERROR: %v", err), buildFields(tags)...)
//		return result
//	}
//
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Error(ctx, fmt.Sprintf("Error reading response body: %v", err), buildFields(tags)...)
//		result.Error = err
//		return result
//	}
//
//	tags = append(tags,
//		fmt.Sprintf("status_code:%d", resp.StatusCode),
//		fmt.Sprintf("response:%.500s", string(body)), // Limit log size
//	)
//	log.Info(ctx, "Response received", buildFields(tags)...)
//
//	result.Body = body
//	result.StatusCode = resp.StatusCode
//
//	// Handle non-2xx status codes
//	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
//		errMsg := fmt.Sprintf("request failed with status %d: %s", resp.StatusCode, string(body))
//		result.Error = fmt.Errorf(errMsg)
//	}
//
//	return result
//}
//
//func buildFields(tags []string) []interface{} {
//	fields := make([]interface{}, 0, len(tags)*2)
//	for _, tag := range tags {
//		parts := strings.SplitN(tag, ":", 2)
//		if len(parts) == 2 {
//			fields = append(fields, parts[0], parts[1])
//		}
//	}
//	return fields
//}
//
//// Get performs an HTTP GET request
//func (c *RustyClient) Get(ctx context.Context, url string, headers map[string]string, queryParams map[string]string, tags []string) RustyResponse {
//	c.logRequest(ctx, http.MethodGet, url, headers, nil, tags)
//
//	builder := c.httpClient.Get(url)
//	builder = c.addHeaders(builder, headers)
//	builder = c.addQueryParams(builder, queryParams)
//
//	resp, err := builder.WithContext(ctx).Do()
//	return c.handleResponse(ctx, resp, err, tags)
//}
//
//// Post performs an HTTP POST request
//func (c *RustyClient) Post(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse {
//	c.logRequest(ctx, http.MethodPost, url, headers, body, tags)
//
//	headers["Content-Type"] = "application/json"
//	builder := c.httpClient.Post(url)
//	builder = c.addHeaders(builder, headers)
//
//	if body != nil {
//		builder = builder.WithJSON(body)
//	}
//
//	resp, err := builder.WithContext(ctx).Do()
//	return c.handleResponse(ctx, resp, err, tags)
//}
//
//// Put performs an HTTP PUT request
//func (c *RustyClient) Put(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse {
//	c.logRequest(ctx, http.MethodPut, url, headers, body, tags)
//
//	headers["Content-Type"] = "application/json"
//	builder := c.httpClient.Put(url)
//	builder = c.addHeaders(builder, headers)
//
//	if body != nil {
//		builder = builder.WithJSON(body)
//	}
//
//	resp, err := builder.WithContext(ctx).Do()
//	return c.handleResponse(ctx, resp, err, tags)
//}
//
//// Patch performs an HTTP PATCH request
//func (c *RustyClient) Patch(ctx context.Context, url string, headers map[string]string, body interface{}, tags []string) RustyResponse {
//	c.logRequest(ctx, http.MethodPatch, url, headers, body, tags)
//
//	headers["Content-Type"] = "application/json"
//	builder := c.httpClient.Patch(url)
//	builder = c.addHeaders(builder, headers)
//
//	if body != nil {
//		builder = builder.WithJSON(body)
//	}
//
//	resp, err := builder.WithContext(ctx).Do()
//	return c.handleResponse(ctx, resp, err, tags)
//}
//
//// Delete performs an HTTP DELETE request
//func (c *RustyClient) Delete(ctx context.Context, url string, headers map[string]string, params map[string]interface{}, tags []string) RustyResponse {
//	c.logRequest(ctx, http.MethodDelete, url, headers, params, tags)
//
//	builder := c.httpClient.Delete(url)
//	builder = c.addHeaders(builder, headers)
//
//	// Add query parameters if any
//	if len(params) > 0 {
//		queryParams := make(map[string]string)
//		for k, v := range params {
//			if s, ok := v.(string); ok {
//				queryParams[k] = s
//			} else if b, err := json.Marshal(v); err == nil {
//				queryParams[k] = string(b)
//			}
//		}
//		builder = c.addQueryParams(builder, queryParams)
//	}
//
//	resp, err := builder.WithContext(ctx).Do()
//	return c.handleResponse(ctx, resp, err, tags)
//}
