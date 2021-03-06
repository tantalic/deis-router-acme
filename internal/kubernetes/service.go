package kubernetes

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Service struct {
	Metadata ObjectMeta  `json:"metadata"`
	Spec     ServiceSpec `json:"Spec"`
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
	Metadata ObjectMeta `json:"metadata,omitempty"`
	Services []Service  `json:"items,omitempty"`
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
	err := c.request(http.MethodGet, "/api/v1/services", params, nil, nil, &results)
	if err != nil {
		return []Service{}, err
	}

	return results.Services, nil
}

func (c *Client) PatchService(secret Secret, patch Secret) error {
	return c.patch(secret.Metadata.Path, patch)
}
