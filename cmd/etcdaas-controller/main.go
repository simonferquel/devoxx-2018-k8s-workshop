package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	types "github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/apis/etcdaas/v1alpha1"
	clientset "github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/client/clientset/versioned"
	informers "github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/client/informers/externalversions"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var (
		kubeconfig string
	)
	flag.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig path (keep unset for using ambient config)")
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	// controllers are implemented as any client to the API
	// it uses informers to watch for changes in the object collection
	c := clientset.NewForConfigOrDie(cfg)
	informerFactory := informers.NewSharedInformerFactory(c, time.Minute)
	informer := informerFactory.Etcdaas().V1alpha1().ETCDInstances().Informer()

	ctx := context.Background()
	ctr := &controller{config: cfg, client: c}
	informer.AddEventHandler(ctr)
	informer.Run(ctx.Done())
}

type controller struct {
	config *restclient.Config
	client clientset.Interface
}

func (c *controller) OnAdd(obj interface{}) {
	etcd, ok := obj.(*types.ETCDInstance)
	if !ok {
		panic("unexpected object")
	}
	fmt.Printf("OnAdd:\n %#v\n", *etcd)

	etcd.Status = types.ETCDInstanceStatus{
		State:   types.ETCDDeploying,
		Message: "deploying",
	}
	etcd, _ = c.client.EtcdaasV1alpha1().ETCDInstances(etcd.Namespace).Update(etcd)

	time.Sleep(10 * time.Second)

	etcd.Status = types.ETCDInstanceStatus{
		State:   types.ETCDRunning,
		Message: "deployment successfull",
	}
	c.client.EtcdaasV1alpha1().ETCDInstances(etcd.Namespace).Update(etcd)
}

func (c *controller) OnUpdate(oldObj, newObj interface{}) {
	oldetcd, ok := oldObj.(*types.ETCDInstance)
	if !ok {
		panic("unexpected object")
	}
	newetcd, ok := newObj.(*types.ETCDInstance)
	if !ok {
		panic("unexpected object")
	}
	fmt.Printf("OnUpdate:\n %#v\n to\n %#v\n", *oldetcd, *newetcd)

	if oldetcd.Spec == newetcd.Spec {
		// don't do anything
		return
	}

	newetcd.Status = types.ETCDInstanceStatus{
		State:   types.ETCDDeploying,
		Message: "updating",
	}
	newetcd, _ = c.client.EtcdaasV1alpha1().ETCDInstances(newetcd.Namespace).Update(newetcd)

	time.Sleep(10 * time.Second)

	newetcd.Status = types.ETCDInstanceStatus{
		State:   types.ETCDRunning,
		Message: "update successfull",
	}
	c.client.EtcdaasV1alpha1().ETCDInstances(newetcd.Namespace).Update(newetcd)
}

func (c *controller) OnDelete(obj interface{}) {
	etcd, ok := obj.(*types.ETCDInstance)
	if !ok {
		panic("unexpected object")
	}
	fmt.Printf("OnDelete:\n %#v\n", *etcd)
}
