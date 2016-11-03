package main

import (
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"

	"tantalic.com/deis-router-acme/internal/kubernetes"
)

func main() {
	opts := OptionFromEnv()

	certNeededChan := make(chan kubernetes.Service)
	errorChan := make(chan error)

	go certNeededLoop(30*time.Second, opts, certNeededChan, errorChan)

	for {
		select {
		case service := <-certNeededChan:
			log.Printf("Service %s needs certificate(s):\n", service.Metadata.Name)
			for _, domain := range certsMissingForService(service, opts) {
				log.Printf("    %s", domain)
			}
		case err := <-errorChan:
			log.Printf("ERROR: %s\n", err)
		}
	}
}

func certNeededLoop(sleep time.Duration, opts options, certNeededChan chan kubernetes.Service, errorChan chan error) {
	for {
		c := kubernetes.Client{}
		services, err := c.ServicesMatchingSelector(opts.RoutableServiceSelector)
		if err != nil {
			errorChan <- errors.Wrap(err, "kubernetes service lookup failed")
			continue
		}

		for _, service := range services {
			if serviceNeedsCert(service, opts) {
				certNeededChan <- service
			}
		}

		time.Sleep(sleep)
	}
}

func serviceNeedsCert(service kubernetes.Service, opts options) bool {
	needed := certsMissingForService(service, opts)
	return len(needed) > 0
}

func certsMissingForService(service kubernetes.Service, opts options) []string {
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

func domainsForService(service kubernetes.Service, opt options) []string {
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

func certificatesForService(service kubernetes.Service, opt options) map[string]string {
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
