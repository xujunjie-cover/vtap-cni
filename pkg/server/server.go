package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/xujunjie-cover/vtap-cni/pkg/config"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	http.Server
	rundir     string
	kubeclient *kubernetes.Interface
}

func FilesystemPreRequirements(rundir string) error {
	if err := os.RemoveAll(rundir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove old pod info socket directory %s: %v", rundir, err)
	}
	if err := os.MkdirAll(rundir, 0700); err != nil {
		return fmt.Errorf("failed to create pod info socket directory %s: %v", rundir, err)
	}
	return nil
}

// GetListener creates a listener to a unix socket located in `socketPath`
func GetListener(socketPath string) (net.Listener, error) {
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on pod info socket: %v", err)
	}
	if err := os.Chmod(socketPath, 0600); err != nil {
		_ = l.Close()
		return nil, fmt.Errorf("failed to listen on pod info socket: %v", err)
	}
	return l, nil
}

func NewCNIServer(daemonConfig *config.Configuration, handler http.Handler) (*Server, error) {
	s := &Server{
		Server: http.Server{
			Handler: handler,
		},
		rundir:     daemonConfig.SocketDir,
		kubeclient: &daemonConfig.KubeClient,
	}
	return s, nil
}
