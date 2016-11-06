package kubernetes

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	Scheme     string
	Host       string
	HTTPClient *http.Client
}

func (c *Client) request(method string, path string, params url.Values, headers http.Header, body interface{}, v interface{}) error {
	if c.Scheme == "" {
		c.Scheme = "http"
	}

	if c.Host == "" {
		c.Host = "127.0.0.1:8001"
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	u := url.URL{
		Scheme:   c.Scheme,
		Host:     c.Host,
		Path:     path,
		RawQuery: params.Encode(),
	}

	var bodyReader io.ReadWriter = nil
	if body != nil {
		bodyReader = new(bytes.Buffer)
		encoder := json.NewEncoder(bodyReader)
		err := encoder.Encode(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return err
	}

	if headers != nil {
		for key, header := range headers {
			for _, value := range header {
				req.Header.Add(key, value)
			}
		}
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		var status Status
		err = json.Unmarshal(data, &status)
		if err != nil {
			return err
		}
		return status
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) post(path string, item interface{}) error {
	headers := make(http.Header, 1)
	headers.Set("Content-Type", "application/json")
	return c.request(http.MethodPost, path, nil, headers, item, nil)
}

func (c *Client) patch(path string, patch interface{}) error {
	headers := make(http.Header, 1)
	headers.Set("Content-Type", "application/strategic-merge-patch+json")

	return c.request(http.MethodPatch, path, nil, headers, patch, nil)
}
