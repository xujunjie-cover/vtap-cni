package cni

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/100"
	"github.com/xujunjie-cover/vtap-cni/pkg/request"
)

func CmdAdd(args *skel.CmdArgs) error {
	netConf, cniVersion, err := loadConf(args.StdinData)
	if err != nil {
		return err
	}
	podName, err := parseValueFromArgs("K8S_POD_NAME", args.Args)
	if err != nil {
		return err
	}
	podNamespace, err := parseValueFromArgs("K8S_POD_NAMESPACE", args.Args)
	if err != nil {
		return err
	}

	client := request.NewCniServerClient(netConf.ServerSocket)
	_, err = client.Add(request.CniRequest{
		CniType:         netConf.Type,
		PodName:         podName,
		PodNamespace:    podNamespace,
		ContainerID:     args.ContainerID,
		DefaultMaster:   netConf.DefaultMaster,
		IsMasterInNetNs: netConf.IsMasterInNetNs,
		NetNs:           args.Netns,
		IfName:          args.IfName,
	})
	if err != nil {
		return err
	}

	result := &current.Result{
		CNIVersion: cniVersion,
		//TODO: response
		Interfaces: []*current.Interface{},
	}

	return types.PrintResult(result, cniVersion)
}

func CmdDel(args *skel.CmdArgs) error {
	netConf, _, err := loadConf(args.StdinData)
	if err != nil {
		return err
	}

	client := request.NewCniServerClient(netConf.ServerSocket)
	podName, err := parseValueFromArgs("K8S_POD_NAME", args.Args)
	if err != nil {
		return err
	}
	podNamespace, err := parseValueFromArgs("K8S_POD_NAMESPACE", args.Args)
	if err != nil {
		return err
	}

	return client.Del(request.CniRequest{
		CniType:         netConf.Type,
		PodName:         podName,
		PodNamespace:    podNamespace,
		ContainerID:     args.ContainerID,
		NetNs:           args.Netns,
		DefaultMaster:   netConf.DefaultMaster,
		IsMasterInNetNs: netConf.IsMasterInNetNs,
		IfName:          args.IfName,
	})
}

type NetConf struct {
	types.NetConf
	ServerSocket    string `json:"server_socket"`
	DefaultMaster   string `json:"default_master"`
	MTU             int    `json:"mtu,omitempty"`
	IsPromiscuous   bool   `json:"promiscMode,omitempty"`
	IsMasterInNetNs bool   `json:"masterInNetNs,omitempty"`
}

func loadConf(bytes []byte) (NetConf, string, error) {
	n := NetConf{}
	if err := json.Unmarshal(bytes, &n); err != nil {
		return n, "", fmt.Errorf("failed to load netconf: %v", err)
	}

	return n, n.CNIVersion, nil
}

func parseValueFromArgs(key, argString string) (string, error) {
	if argString == "" {
		return "", errors.New("CNI_ARGS is required")
	}
	args := strings.Split(argString, ";")
	for _, arg := range args {
		if strings.HasPrefix(arg, fmt.Sprintf("%s=", key)) {
			value := strings.TrimPrefix(arg, fmt.Sprintf("%s=", key))
			if len(value) > 0 {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("%s is required in CNI_ARGS", key)
}
