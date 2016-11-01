package main

import (
	"log"
	"time"

	"github.com/pkg/errors"

	"tantalic.com/deis-workflow-acme/internal/kubernetes"
)

func main() {
	certNeededChan := make(chan kubernetes.Service)
	errorChan := make(chan error)

	go certNeededLoop(30*time.Second, certNeededChan, errorChan)

	for {
		select {
		case service := <-certNeededChan:
			log.Printf("Service %s needs a certificate\n", service.Metadata.Name)
		case err := <-errorChan:
			log.Printf("ERROR: %s\n", err)
		}
	}
}

func certNeededLoop(sleep time.Duration, certNeededChan chan kubernetes.Service, errorChan chan error) {
	for {
		c := kubernetes.Client{}
		services, err := c.ServicesMatchingSelector("router.deis.io/routable=true")
		if err != nil {
			errorChan <- errors.Wrap(err, "kubernetes service lookup failed")
			continue
		}

		for _, service := range services {
			if serviceNeedsCert(service) {
				certNeededChan <- service
			}
		}

		time.Sleep(sleep)
	}
}

func serviceNeedsCert(service kubernetes.Service) bool {
	certificates := service.Metadata.Annotations["router.deis.io/certificates"]
	if certificates == "" {
		return true
	}

	//TODO: Compare the router.deis.io/domains annotation to
	//the router.deis.io/certificates annotation to determine
	//if a certificate is needed.

	return false
}