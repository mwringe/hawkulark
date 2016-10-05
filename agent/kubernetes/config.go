package kubernetes

import (
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/rest"
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
)


const USERAGENT string = "HawkularK"

func GetKubernetesClient (url string, token string, caFile string) ( *v1.CoreClient, error) {

	var config *rest.Config

	// if no values are passed, assume that we are running within the container within the Kubernetes cluster
	if (url == "" || token == "") {
		var err error
		config, err = rest.InClusterConfig()
		if (err != nil) {
			return nil, err
		}
	} else {
		c := rest.Config{
			Host: url,
			BearerToken: token,
		}

		if (caFile != "") {
			tlsConfig := rest.TLSClientConfig{}
			tlsConfig.CAFile = caFile

			c.TLSClientConfig = tlsConfig
		}

		config = &c

	}

	// set our user agent
	config.UserAgent = USERAGENT

	client, err := kubernetes.NewForConfig(config)
	if (err != nil) {
		return nil, err
	}

	return client.CoreClient, nil
}