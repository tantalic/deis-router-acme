package main

import (
	"os"
)

const (
	DefaultRoutableServiceSelector = "router.deis.io/routable=true"
	DefaultPlatformDomain          = ""
	DefaultCertificatesAnnotation  = "router.deis.io/certificates"
	DefaultDomainsAnnotation       = "router.deis.io/domains"
)

type options struct {
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

func OptionFromEnv() options {
	return options{
		RoutableServiceSelector: getEnv("ROUTABLE_SERVICE_SELECTOR", DefaultRoutableServiceSelector),
		PlatformDomain:          getEnv("PLATFORM_DOMAIN", DefaultPlatformDomain),
		CertificatesAnnotation:  getEnv("CERTIFICATES_ANNOTATION", DefaultCertificatesAnnotation),
		DomainsAnnotation:       getEnv("DOMAINS_ANNOTATION", DefaultDomainsAnnotation),
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}

	return val
}
