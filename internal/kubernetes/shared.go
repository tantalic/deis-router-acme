package kubernetes

// ObjectMeta is metadata that all persisted resources must have.
type ObjectMeta struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Path        string            `json:"selfLink,omitempty"`
	UID         string            `json:"uid,omitempty"`
	Version     string            `json:"resourceVersion,omitempty"`
	Created     string            `json:"creationTimestam,omitemptyp`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Status represents an API response when an API operation is not
// succesful or when a HTTP DELETE call is sucessful. The status
// object contains fields for humans and machine consumers of the
// API to get more detailed information for the cause of the failure.
type Status struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Reason  string                 `json:"reason"`
	Details map[string]interface{} `json:"status"`
	Code    int                    `json:"code"`
}

// Error returns the message string for Status
func (s Status) Error() string {
	return s.Message
}
