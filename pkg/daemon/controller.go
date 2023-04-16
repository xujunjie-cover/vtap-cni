package daemon

import (
	"github.com/xujunjie-cover/vtap-cni/pkg/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"k8s.io/kubectl/pkg/scheme"
)

type Controller struct {
	config *config.Configuration

	podsLister listerv1.PodLister
	podsSynced cache.InformerSynced

	recorder record.EventRecorder
}

func NewController(config *config.Configuration, podInformerFactory informers.SharedInformerFactory) (*Controller, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: config.KubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: config.NodeName})
	podInformer := podInformerFactory.Core().V1().Pods()

	controller := &Controller{
		config: config,

		podsLister: podInformer.Lister(),
		podsSynced: podInformer.Informer().HasSynced,

		recorder: recorder,
	}

	return controller, nil
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	if ok := cache.WaitForCacheSync(stopCh, c.podsSynced); !ok {
		klog.Fatalf("failed to wait for caches to sync")
		return
	}

	klog.Info("Started controller worker")

	<-stopCh
	klog.Info("Shutting down controller worker")
}
