package kubernetes

import (
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
	"fmt"
	"os"
)

const (
	POD_NAME  = "POD_NAME"
	NAMESPACE = "NAMESPACE"
)


func GetNodeName(client *v1.CoreClient, podName, namespace string) (string, error) {

	// TODO: move the section here to hawkularkagent.go, the specifics about the envars and parameters should go there
	if (podName == "") {
		podName = os.Getenv(POD_NAME)
		if (podName == "") {
			return "", fmt.Errorf("Could not determine the pod name. The pod name must be passed via --hawkulark_pod_name or set with the %v environment variable", POD_NAME)
		}
	}

	if (namespace == "") {
		namespace = os.Getenv(NAMESPACE)
		if (namespace == "") {
			return "", fmt.Errorf("Could not determine the namespace. The pod name must be passed via --hawkulark_pod_namespace or set with the %v environment variable", NAMESPACE)
		}
	}
	// /TODO

	pod, err := client.Pods(namespace).Get(podName)
	if (err != nil) {
		return "", fmt.Errorf("Error detemining the namespace that the Hawkulark pod is running under: %v", err)
	}

	return pod.Spec.NodeName, nil
}