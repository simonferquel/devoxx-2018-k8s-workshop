package main

import (
	"github.com/simonferquel/devoxx-2018-k8s-workshop/cmd/etcdaas-api/internalversion"
	"github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/apis/etcdaas/v1alpha1"
	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// Internal variables
var (
	groupFactoryRegistry = make(announced.APIGroupFactoryRegistry)
	registry             = registered.NewOrDie("")

	// Scheme of the API Server: contains all types manipulated by the API server
	// (both externally e.g.: v1alpha.ETCDInstance, metav1.APIGroup... , and internally - in this case we use the same types for both external and internal version, but those are still 2 apigroupversions registered)
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	// this register the groupAPI with all its versions (and version preference)
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  v1alpha1.GroupName,
			VersionPreferenceOrder:     []string{v1alpha1.SchemeGroupVersion.Version},
			AddInternalObjectsToScheme: internalversion.AddToScheme,
		},
		announced.VersionToSchemeFunc{
			v1alpha1.SchemeGroupVersion.Version: v1alpha1.AddToScheme,
		},
	).Announce(groupFactoryRegistry).RegisterAndEnable(registry, Scheme); err != nil {
		panic(err)
	}

	// we also need types from metav1 for exchanging API metadata with the central API server
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// and some unversioned flavors as well
	// note: it is on k8s kube-aggregation TODO-list to remove these requiresments
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
	internalversion.AddToScheme(Scheme)
}
