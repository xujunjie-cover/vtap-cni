package main

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
	"github.com/xujunjie-cover/vtap-cni/pkg/cni"
)

func main() {
	skel.PluginMain(cni.CmdAdd, nil, cni.CmdDel, version.All, bv.BuildString("vtap-cni"))
}
