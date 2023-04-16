package api

import "path/filepath"

const (
	CNIAPIEndpoint    = "/cni"
	HealthAPIEndpoint = "/healthz"

	serverSocketName = "vtap-cni.sock"
)

// SocketPath returns the path of the multus CNI socket
func SocketPath(rundir string) string {
	return filepath.Join(rundir, serverSocketName)
}
