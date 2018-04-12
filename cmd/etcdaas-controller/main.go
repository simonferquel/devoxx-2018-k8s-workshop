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
	c := clientset.NewForConfigOrDie(cfg)
	informerFactory := informers.NewSharedInformerFactory(c, time.Minute)
	informer := informerFactory.Etcdaas().V1alpha1().ETCDInstances().Informer()

	ctx := context.Background()
	ctr := &controller{config: cfg}
	informer.AddEventHandler(ctr)
	informer.Run(ctx.Done())
}

type controller struct {
	config *restclient.Config
}

func (c *controller) OnAdd(obj interface{}) {
	etcd, ok := obj.(*types.ETCDInstance)
	if !ok {
		panic("unexpected object")
	}
	fmt.Printf("OnAdd:\n %#v\n", *etcd)
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
}
func (c *controller) OnDelete(obj interface{}) {
	etcd, ok := obj.(*types.ETCDInstance)
	if !ok {
		panic("unexpected object")
	}
	fmt.Printf("OnDelete:\n %#v\n", *etcd)
}
