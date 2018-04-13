package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	types "github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/apis/etcdaas/v1alpha1"
	client "github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/client/clientset/versioned/typed/etcdaas/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var create, list, delete bool
	var name string
	var namespace string
	var replicas int
	var tls bool
	flag.BoolVar(&create, "create", false, "create an etcd")
	flag.BoolVar(&list, "list", false, "list etcds")
	flag.BoolVar(&delete, "delete", false, "delete an etcd")
	flag.BoolVar(&tls, "enable-tls", false, "enable tls")
	flag.StringVar(&name, "name", "", "etcd name")
	flag.StringVar(&namespace, "namespace", "default", "k8s namespace")
	flag.IntVar(&replicas, "replicas", 1, "etcd replicas")
	flag.Parse()

	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	cfg, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		panic(err)
	}

	c := client.NewForConfigOrDie(cfg)

	switch {
	case create:
		i := &types.ETCDInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: types.ETCDInstanceSpec{
				Replicas:      replicas,
				WithTLSBundle: tls,
			},
		}
		_, err := c.ETCDInstances(namespace).Create(i)
		if err != nil {
			panic(err)
		}
	case list:
		lst, err := c.ETCDInstances(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, item := range lst.Items {
			fmt.Printf("%s:\n\tspec: %#v,\n\tstatus: %#v\n", item.Name, item.Spec, item.Status)
		}
	case delete:
		if err := c.ETCDInstances(namespace).Delete(name, nil); err != nil {
			panic(err)
		}
	}
}
