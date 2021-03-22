package client

import (
	"bytes"
	"fmt"
	"github.com/go-http-utils/headers"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	FormEncoded      = "application/x-www-form-urlencoded"
	ApplicationJson  = "application/json"
	ApplicationXml   = "application/xml"
	IdSeparator      = ":"
	Basic            = "Basic"
	Bearer           = "Bearer"
	AzureGraphServer = "graph.microsoft.com"
)

type Client struct {
	tenantId    string
	accessToken string
	httpClient  *http.Client
}

func NewClient(tenantId string, accessToken string) (client *Client, err error) {
	c := &Client{
		tenantId:    tenantId,
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}
	return c, nil
}

func (c *Client) HttpRequest(method string, path string, query url.Values, headerMap http.Header, body *bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), body)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	//Handle query values
	if query != nil {
		requestQuery := req.URL.Query()
		for key, values := range query {
			for _, value := range values {
				requestQuery.Add(key, value)
			}
		}
		req.URL.RawQuery = requestQuery.Encode()
	}
	//Handle header values
	if headerMap != nil {
		for key, values := range headerMap {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	//Handle authentication
	if c.accessToken != "" {
		req.Header.Set(headers.Authorization, Bearer+" "+c.accessToken)
		//} else {
		//	req.SetBasicAuth(c.username, c.password)
	}
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Print("Apigee Management API:")
		log.Print(err)
	} else {
		log.Print("Apigee Management API: " + string(requestDump))
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, &RequestError{StatusCode: resp.StatusCode, Err: err}
		}
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("https://%s/v1.0/%s", AzureGraphServer, path)
}
