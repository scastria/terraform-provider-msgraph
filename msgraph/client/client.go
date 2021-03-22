package client

import (
	"bytes"
	"encoding/json"
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
	AzureAuthServer  = "login.microsoftonline.com"
)

type Client struct {
	tenantId     string
	accessToken  string
	clientId     string
	clientSecret string
	httpClient   *http.Client
}

func NewClient(tenantId string, accessToken string, clientId string, clientSecret string) (client *Client, err error) {
	c := &Client{
		tenantId:     tenantId,
		accessToken:  accessToken,
		clientId:     clientId,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
	}
	//Check for client credentials authentication and try to get access token
	if c.accessToken == "" {
		log.Print("Microsoft Graph API: Obtaining access token...")
		requestURL := fmt.Sprintf("https://%s/%s/oauth2/v2.0/token", AzureAuthServer, c.tenantId)
		requestForm := url.Values{
			"grant_type":    []string{"client_credentials"},
			"scope":         []string{"https://graph.microsoft.com/.default"},
			"client_id":     []string{c.clientId},
			"client_secret": []string{c.clientSecret},
		}
		req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBufferString(requestForm.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set(headers.ContentType, FormEncoded)
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Print("Microsoft Graph API:")
			log.Print(err)
		} else {
			log.Print("Microsoft Graph API: " + string(requestDump))
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
		//Parse body to extract access_token
		token := &OauthToken{}
		err = json.NewDecoder(resp.Body).Decode(token)
		if err != nil {
			return nil, err
		}
		log.Print("Microsoft Graph API: Received access token: " + token.AccessToken)
		//Inject token as access_token for client for all future calls
		c.accessToken = token.AccessToken
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
		log.Print("Microsoft Graph API:")
		log.Print(err)
	} else {
		log.Print("Microsoft Graph API: " + string(requestDump))
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
