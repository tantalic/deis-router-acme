package main

import (
	"os"
)

const (
	// DefaultRoutableServiceSelector is the default value for options.RoutableServiceSelector
	DefaultRoutableServiceSelector = "router.deis.io/routable=true"
	// DefaultPlatformDomain is the default value for options.PlatformDomain
	DefaultPlatformDomain = ""
	// DefaultCertificatesAnnotation is the default value for CertificatesAnnotation
	DefaultCertificatesAnnotation = "router.deis.io/certificates"
	// DefaultDomainsAnnotation is the default value for DomainsAnnotation
	DefaultDomainsAnnotation = "router.deis.io/domains"
)

// Options collects the options that can be used to control the execution
// of the application. Typically, an instance would be populated from the
// OptionFromEnv function.
type Options struct {
	// The kubernetes service selector for discovering services
	// that should have certificates.
	RoutableServiceSelector string

	// The router's platform domain. Any domains not containing
	// the . character will be assumed to be subdomains of this
	// domain.
	PlatformDomain string

	// The key for the annotation containing mappings between
	// domain names (see DomainsAnnotation) and the certificate
	// to be used.
	CertificatesAnnotation string

	// The key for the annotation containing a comma-delimited
	// list of domains that should routed to each application.
	DomainsAnnotation string
}

// OptionFromEnv returns a new options struct with values populated
// from the environment variables.
func OptionFromEnv() Options {
	return Options{
		RoutableServiceSelector: getEnv("ROUTABLE_SERVICE_SELECTOR", DefaultRoutableServiceSelector),
		PlatformDomain:          getEnv("PLATFORM_DOMAIN", DefaultPlatformDomain),
		DomainsAnnotation:       getEnv("DOMAINS_ANNOTATION", DefaultDomainsAnnotation),
	}
}

// getEnv retrieves the value of the environment variable named
// by the key or the provided default value if there is no
// environment variabled with the named key.
func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}

	return val
}
