package kubernetes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	Scheme     string
	Host       string
	HTTPClient *http.Client
}

func (c *Client) request(method string, path string, params url.Values, v interface{}) error {

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

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return err
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

	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return nil
}

type Service struct {
	Metadata Metadata    `json:"metadata"`
	Spec     ServiceSpec `json:"Spec"`
}

type Metadata struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Path        string            `json:"selfLink,omitempty"`
	UID         string            `json:"uid,omitempty"`
	Version     string            `json:"resourceVersion,omitempty"`
	Created     string            `json:"creationTimestam,omitemptyp`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ServiceSpec struct {
	Ports                    []ServicePort     `json:"ports,omitempty"`
	Selector                 map[string]string `json:"selector,omitempty"`
	ClusterIP                string            `json:"clusterIP,omitempty"`
	Type                     string            `json:"type,omitempty"`
	ExternalIPs              []string          `json:"externalIPs,omitempty"`
	SessionAffinity          string            `json:"sessionAffinity,omitempty"`
	LoadBalancerIP           string            `json:"loadBalancerIP,omitempty"`
	LoadBalancerSourceRanges []string          `json:"loadBalancerSourceRanges,omitempty"`
	ExternalName             string            `json:"externalName,omitempty"`
}

type ServicePort struct {
	Name       string      `json:"name,omitempty"`
	Protocol   string      `json:"protocol,omitempty"`
	Port       int32       `json:"port,omitempty"`
	TargetPort json.Number `json:"targetPort,omitempty"`
	NodePort   json.Number `json:"nodePort,omitempty"`
}

type ServiceList struct {
	Metadata Metadata  `json:"metadata,omitempty"`
	Services []Service `json:"items,omitempty"`
}

func (c *Client) AllServices() ([]Service, error) {
	return c.ServicesMatchingSelector("")
}

func (c *Client) ServicesMatchingSelector(labelSelector string) ([]Service, error) {
	params := url.Values{}
	if labelSelector != "" {
		params.Add("labelSelector", labelSelector)
	}

	var results ServiceList
	err := c.request(http.MethodGet, "/api/v1/services", params, &results)
	if err != nil {
		return []Service{}, err
	}

	return results.Services, nil
}
