package main

import (
	"github.com/golang/glog"
	"flag"
	"strings"
	"os"
	"time"
	"fmt"
	"net/http"
	"crypto/tls"
	"github.com/hawkular/hawkulark/agent/kubernetes"
	"github.com/hawkular/hawkulark/agent/manager"
	"os/signal"
)

// Note: the actual version and gitcommit are set via ldflags during the build.
// Do not set the values here
var version = "unknown"
var gitcommit = "unknown"

var (
	argKubernetesMasterURL = flag.String("master_url", "https://kubernetes.default.svc:443", "The URL to connect to the Kubernetes master.")
	argKubernetesToken = flag.String("kubernetes_token", "", "The token to be used to authenticate with Kubernetes with")
	argKubernetesCA = flag.String("kubernetes_ca_file", "", "The CA file used to sign the Kubernetes Master API with")

	argHawkularkPodName = flag.String("hawkulark_pod_name", "hawkulark", "If running in a container, this is the name of the Hawkular pod")
	argHawkularkPodNamespace = flag.String("hawkulark_pod_namespace", "", "If running in a container, this is the name of the project the Hawkulark pod is running under")

	argMinMetricsResolution = flag.Duration("min_metric_resolution", 10*time.Second, "The minimum resolution custom metrics can request")
	argDefaultMetricResolution = flag.Duration("default_metrics_resolution", 30*time.Second, "The default resolution for collecting metrics")

	argCollectorClientCertFile = flag.String("collector-client-cert-file", "", "The file containing the client certificate for the collector")
	argCollectorClientPrivateKeyFile = flag.String("collector-client-private-key-file", "", "The file containing the corresponding key for --collector-client-cert-file")
)

func init() {
	// log everything to stderr so that it can be easily gathered by logs, separate log files are problematic with containers
	flag.Set("logtostderr", "true")
}

func main() {
	defer glog.Flush()
	flag.Parse()

	glog.Info(strings.Join(os.Args, " "))
	glog.Infof("Starting the Hawkular on Kubernetes Agent [%v:%v]", version, gitcommit)

	if err:= validateFlags(); err != nil {
		glog.Fatal(err)
	}

	httpClient,err := getHttpClient(*argCollectorClientCertFile, *argCollectorClientPrivateKeyFile)
	if err != nil {
		glog.Fatalf("Error trying to configure the Collector's client certificate : %v", err)
	}

	//TODO: the httpClient here should be used by what is doing the collecting
        //TODO: remove the printing once we are actually using it
	glog.Warning("HTTPCLIENT ", httpClient)


	//TODO: we should be handling this in a better manner. Specifying the values here is really only
	// useful when dealing with development. When running within a container we should use the
	client,err := kubernetes.GetKubernetesClient(*argKubernetesMasterURL, *argKubernetesToken, *argKubernetesCA)
	if err != nil {
		glog.Fatal("Error trying to get the Kubernetes Client : %v", err)
	}

	nodeName, err := kubernetes.GetNodeName(client, *argHawkularkPodName, *argHawkularkPodNamespace)
	if err != nil {
		glog.Fatal(err)
	}

	monitor := kubernetes.NewMonitor(client, nodeName)
	monitor.Start()

	manager, err := manager.NewManager()
	if err != nil {
		glog.Fatal("Failed creating the manager : %v", err)
	}
	manager.Start()

	signalChan := make (chan os.Signal, 1)
	doneChan := make (chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			glog.Warning("Termination Signal Received")
			doneChan <- true
		}

	}()
	<-doneChan

	monitor.Stop()
	glog.Warning("DONE")
}

func validateFlags() error {
	if *argMinMetricsResolution < 5*time.Second {
		return fmt.Errorf("The minimum metrics resolution (%d) cannot be less 5 seconds.", *argMinMetricsResolution)
	}

	if *argDefaultMetricResolution < * argMinMetricsResolution {
		return fmt.Errorf("The default metrics resolution (%d) cannot be less than the minimum resolution (%d).", *argDefaultMetricResolution, *argMinMetricsResolution)
	}

	if (*argCollectorClientCertFile != "" && *argCollectorClientPrivateKeyFile == "") ||
		(*argCollectorClientCertFile == "" && *argCollectorClientPrivateKeyFile != "") {
		return fmt.Errorf("Both the --collector-client-cert-file and --collector-client-private-key-file must be configured if one is specified.")
	}

	return nil
}

func getHttpClient(collectorClientCertFile, collectorClientPrivateKeyFile string) (*http.Client, error) {
	//Enable accessing insecure endpoints. We should be able to access metrics from any endpoint
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	if collectorClientCertFile != "" {
		cert, err := tls.LoadX509KeyPair(collectorClientCertFile, collectorClientPrivateKeyFile)
		if err != nil {
			return nil, fmt.Errorf("Error loading the collector client certificates: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := http.Client{Transport: transport}

	return &httpClient, nil
}