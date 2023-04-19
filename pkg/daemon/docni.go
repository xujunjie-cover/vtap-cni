package daemon

import (
	"fmt"
	"net/http"

	current "github.com/containernetworking/cni/pkg/types/100"
	"github.com/emicklei/go-restful"
	"github.com/vishvananda/netlink"
	"github.com/xujunjie-cover/vtap-cni/pkg/request"
	"k8s.io/klog"
)

func (controller Controller) handleAdd(req *restful.Request, resp *restful.Response) {

	_, err := controller.CreateVtapInterface(req)
	if err != nil {
		errMsg := fmt.Errorf("parse add request failed %v", err)
		if err := resp.WriteHeaderAndEntity(http.StatusBadRequest, request.CniResponse{Err: errMsg.Error()}); err != nil {
			klog.Errorf("failed to write response, %v", err)
		}
		return
	}

	// result := &current.Result{
	// 	CNIVersion: cniVersion,
	// 	Interfaces: []*current.Interface{vtapInterface},
	// }

}

func (controller Controller) handleDel(req *restful.Request, resp *restful.Response) {
	//todo
}

func (controller Controller) CreateVtapInterface(req *restful.Request) (*current.Interface, error) {
	podRequest := request.CniRequest{}
	if err := req.ReadEntity(&podRequest); err != nil {
		return nil, err
	}

	klog.Errorf("provider %s %v", podRequest.Provider, podRequest)

	_, err := controller.podsLister.Pods(podRequest.PodNamespace).Get(podRequest.PodName)
	if err != nil {
		return nil, err
	}

	_, err = netlink.LinkByName(podRequest.DefaultMaster)
	return nil, nil
}
