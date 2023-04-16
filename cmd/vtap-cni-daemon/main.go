package main

import (
	"os"

	"github.com/xujunjie-cover/vtap-cni/pkg/config"
	"github.com/xujunjie-cover/vtap-cni/pkg/daemon"

	"k8s.io/klog/v2"
	"k8s.io/sample-controller/pkg/signals"
)

const (
	defaultSocketFileName = "vtap-cni.sock"
)

func main() {
	config, err := config.ParseFlags()
	if err != nil {
		klog.Fatalf("parse config failed %v", err)
	}
	stopCh := signals.SetupSignalHandler()

	if err := daemon.StartVtapCniDaemon(config, stopCh); err != nil {
		klog.Fatalf("failed start the vtap cni listener: %v", err)
		os.Exit(3)
	}
}
