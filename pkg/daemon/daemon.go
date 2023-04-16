package daemon

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/xujunjie-cover/vtap-cni/pkg/api"
	"github.com/xujunjie-cover/vtap-cni/pkg/config"
	"github.com/xujunjie-cover/vtap-cni/pkg/request"
	srv "github.com/xujunjie-cover/vtap-cni/pkg/server"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/klog/v2"
)

const (
	requestLogFormat  = "[%s] Incoming %s %s %s request"
	responseLogFormat = "[%s] Outgoing response %s %s with %d status code in %vms"
)

func StartVtapCniDaemon(conf *config.Configuration, stopCh <-chan struct{}) error {
	if err := srv.FilesystemPreRequirements(conf.SocketDir); err != nil {
		return fmt.Errorf("failed to prepare the cni-socket for communicating: %w", err)
	}

	podInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(conf.KubeClient, 0,
		kubeinformers.WithTweakListOptions(func(listOption *v1.ListOptions) {
			listOption.FieldSelector = fmt.Sprintf("spec.nodeName=%s", conf.NodeName)
			listOption.AllowWatchBookmarks = true
		}))

	controller, err := NewController(conf, podInformerFactory)
	if err != nil {
		klog.Fatalf("create controller failed %v", err)
	}
	podInformerFactory.Start(stopCh)

	go controller.Run(stopCh)

	server, err := srv.NewCNIServer(conf, createHandler(controller))
	if err != nil {
		return fmt.Errorf("failed to create the server: %v", err)
	}

	l, err := srv.GetListener(api.SocketPath(conf.SocketDir))
	if err != nil {
		return fmt.Errorf("failed to start the CNI server using socket %s. Reason: %+v", api.SocketPath(conf.SocketDir), err)
	}

	server.SetKeepAlivesEnabled(false)
	go func() {
		wait.Until(func() {
			klog.Info("open for work")
			if err := server.Serve(l); err != nil {
				runtime.HandleError(fmt.Errorf("CNI server Serve() failed: %v", err))
			}
		}, 0, stopCh)
		server.Shutdown(context.TODO())
	}()

	return nil
}

func createHandler(controller *Controller) http.Handler {
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	ws := new(restful.WebService)
	ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(ws)

	ws.Route(
		ws.POST("/add").
			To(controller.handleAdd).
			Reads(request.CniRequest{}))
	ws.Route(
		ws.POST("/del").
			To(controller.handleDel).
			Reads(request.CniRequest{}))

	ws.Filter(requestAndResponseLogger)

	return wsContainer
}

func requestAndResponseLogger(request *restful.Request, response *restful.Response,
	chain *restful.FilterChain) {
	klog.Infof(formatRequestLog(request))
	start := time.Now()
	chain.ProcessFilter(request, response)
	elapsed := float64((time.Since(start)) / time.Millisecond)
	klog.Infof(formatResponseLog(response, request, elapsed))
}

func formatRequestLog(request *restful.Request) string {
	return fmt.Sprintf(requestLogFormat, time.Now().Format(time.RFC3339), request.Request.Proto,
		request.Request.Method, getRequestURI(request))
}

// formatResponseLog formats response log string.
func formatResponseLog(response *restful.Response, request *restful.Request, reqTime float64) string {
	return fmt.Sprintf(responseLogFormat, time.Now().Format(time.RFC3339),
		request.Request.Method, getRequestURI(request), response.StatusCode(), reqTime)
}

// getRequestURI get the request uri
func getRequestURI(request *restful.Request) (uri string) {
	if request.Request.URL != nil {
		uri = request.Request.URL.RequestURI()
	}
	return
}
