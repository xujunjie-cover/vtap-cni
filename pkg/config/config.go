package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

const (
	defaultSocketDir = "/run/vtap-cni/"
)

type Configuration struct {
	SocketDir      string
	KubeConfigFile string
	NodeName       string
	KubeClient     kubernetes.Interface
}

func ParseFlags() (*Configuration, error) {
	var (
		argSocketDir      = pflag.String("socketdir", defaultSocketDir, "Specify the path to the vtap-cni socket")
		argKubeConfigFile = pflag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information. If not set use the inCluster token.")
	)
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)

	// Sync the glog and klog flags.
	pflag.CommandLine.VisitAll(func(f1 *pflag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			if err := f2.Value.Set(value); err != nil {
				klog.Fatalf("failed to set pflag, %v", err)
			}
		}
	})

	pflag.CommandLine.AddGoFlagSet(klogFlags)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()

	config := &Configuration{
		KubeConfigFile: *argKubeConfigFile,
		SocketDir:      *argSocketDir,
	}
	if err := config.initKubeClient(); err != nil {
		return nil, err
	}

	return config, nil
}

func (config *Configuration) initKubeClient() error {
	var cfg *rest.Config
	var err error
	if config.KubeConfigFile != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", config.KubeConfigFile)
	} else if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		klog.Infof("no --kubeconfig, use in-cluster kubernetes config")
		cfg, err = rest.InClusterConfig()
	} else {
		err = fmt.Errorf("No kubernetes config")
	}

	if err != nil {
		klog.Errorf("failed to build kubeconfig %v", err)
		return err
	}

	// Specify that we use gRPC
	cfg.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	cfg.ContentType = "application/vnd.kubernetes.protobuf"
	// Set the config timeout to one minute.
	cfg.Timeout = time.Minute

	return nil
}
