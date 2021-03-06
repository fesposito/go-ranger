package ranger_http

import (
	"fmt"
	"net/http"
	"time"
)

// APIClientInterface is an interface for api clients. It allows us to mock the basic http client.
type APIClientInterface interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}

type apiClient struct {
	client *http.Client
}

// NewAPIClient is the factory method for api clients.
func NewAPIClient(requestTimeout int) APIClientInterface {
	return &apiClient{
		client: &http.Client{
			Timeout: time.Second * time.Duration(requestTimeout),
		},
	}
}

// Get is issueing a GET request to the given url
func (client *apiClient) Get(url string) (*http.Response, error) {
	res, err := client.client.Get(url)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"ApiClient.Get=Bad request,StatusCode=%d, URL=%s, Header: %+v", res.StatusCode, url, res.Header,
		)
	}

	return res, err
}

// Do sends an HTTP request and returns an HTTP response, following
// policy (such as redirects, cookies, auth) as configured on the
// client.
func (client *apiClient) Do(req *http.Request) (*http.Response, error) {
	res, err := client.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"ApiClient.Do=Cannot execute request, URL=%s, Header=%+v", req.URL, req.Header,
		)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"ApiClient.Do=Bad request,StatusCode=%d, URL=%s, Header: %+v", res.StatusCode, req.URL, res.Header,
		)
	}

	return res, err
}
