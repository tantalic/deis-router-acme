package main

import (
	"strings"

	"tantalic.com/deis-router-acme/internal/kubernetes"
)

func serviceNeedsCert(service kubernetes.Service, opts Options) bool {
	needed := certsMissingForService(service, opts)
	return len(needed) > 0
}

func certsMissingForService(service kubernetes.Service, opts Options) []string {
	domains := domainsForService(service, opts)
	certificates := certificatesForService(service, opts)

	var val []string
	for _, domain := range domains {
		_, exists := certificates[domain]
		if !exists {
			val = append(val, domain)
		}
	}

	return val
}

func domainsForService(service kubernetes.Service, opt Options) []string {
	list := service.Metadata.Annotations[opt.DomainsAnnotation]
	domains := strings.Split(list, ",")

	var val []string
	for _, domain := range domains {
		// Considered a fully qualified domain if it contains a "."
		if strings.Contains(domain, ".") {
			val = append(val, domain)
			continue
		}

		// subdomains only included if platform domain is set
		if opt.PlatformDomain != "" {
			val = append(val, domain+"."+opt.PlatformDomain)
			continue
		}
	}

	return val
}

func certificatesForService(service kubernetes.Service, opt Options) map[string]string {
	list := service.Metadata.Annotations[opt.CertificatesAnnotation]
	domains := strings.Split(list, ",")

	val := make(map[string]string)
	for _, domain := range domains {
		x := strings.Split(domain, ":")

		// There should always be exactly two values after splitting
		if len(x) != 2 {
			continue
		}

		domain := x[0]
		cert := x[1]
		val[domain] = cert
	}

	return val
}
