package kubernetes

import (
	"github.com/golang/glog"
	"k8s.io/client-go/1.4/pkg/api/v1"
	core "k8s.io/client-go/1.4/kubernetes/typed/core/v1"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/fields"
)

type Monitor interface {
	Start()
	Stop()
}


type monitor struct {
	client *core.CoreClient
	nodename string
	monitoredPods map[string]v1.Pod
}

func NewMonitor(client *core.CoreClient, nodename string) Monitor {
	m := monitor{
		client: client,
		nodename: nodename,
	}

	monitoredPods := make(map[string]v1.Pod)
	m.monitoredPods = monitoredPods

	return &m
}

func (m *monitor) Start() {
	glog.V(4).Info("Starting the Kubernetes Monitor")

	// we only want to listen to pods on our own node
	fieldSelector := fields.OneTermEqualSelector("spec.nodeName", m.nodename)

	listOptions := api.ListOptions{
		Watch: true,
		FieldSelector: fieldSelector,
	}

	watcher, err := m.client.Pods(v1.NamespaceAll).Watch(listOptions)
	if (err != nil) {
		glog.Fatal(err)
	}

	go func() {
		for event := range watcher.ResultChan() {
			pod := event.Object.(*v1.Pod)
			glog.Warning("EVENT :", event.Type, "  --- POD NAME : ", pod.Name, " --- ANNOTATIONS : ", pod.Annotations)
		}
	}()
}

func (m *monitor) Stop() {
	glog.V(4).Info("Shuting down the Kubernetes Monitor")
}
