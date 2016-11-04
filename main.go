package main

import (
	"log"
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

func certNeededLoop(sleep time.Duration, opts Options, certNeededChan chan kubernetes.Service, errorChan chan error) {
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
