package main

import (
	"context"
	"net"

	reg "github.com/simonferquel/devoxx-2018-k8s-workshop/cmd/etcdaas-api/registry"
	"github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/apis/etcdaas/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/util/logs"
)

// path-prefix for objects instances
const defaultEtcdPathPrefix = "/registry/etcdaas/etcdinstances"

type apiServerOptions struct {
	// these are given by k8s apiserver sdk
	RecommendedOptions *genericoptions.RecommendedOptions
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	// this codec is the one used to actually store data in etcd. This means "store data in format v1alpha1"
	codec := Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion)
	o := &apiServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(defaultEtcdPathPrefix, codec),
	}
	rootCmd := &cobra.Command{
		Short: "Launch api server",
		RunE: func(c *cobra.Command, args []string) error {
			errors := []error{}
			errors = append(errors, o.RecommendedOptions.Validate()...)
			if err := utilerrors.NewAggregate(errors); err != nil {
				return err
			}
			ctx := context.Background()
			return runEtcdAPI(o, ctx.Done())

		},
	}
	// this add recommended API server flags to the cobra command
	o.RecommendedOptions.AddFlags(rootCmd.Flags())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func runEtcdAPI(o *apiServerOptions, stopCh <-chan struct{}) error {
	// if no explicit tls configuration is provided, fallback to auto-generated self signed certs
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return err
	}

	// generate the apiserver config from flags
	serverConfig := genericapiserver.NewRecommendedConfig(Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig, Scheme); err != nil {
		return err
	}

	// prepare the generic api server to run
	completeConfig := serverConfig.Complete()
	completeConfig.Version = &version.Info{
		Major: "1",
		Minor: "0",
	}

	srv, err := completeConfig.New("etcdAPIServer", genericapiserver.EmptyDelegate)
	if err != nil {
		return err
	}

	// populate apiGroupInfos we want to expose (as native k8s APIs)
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(v1alpha1.GroupName, registry, Scheme, metav1.ParameterCodec, Codecs)
	apiGroupInfo.GroupMeta.GroupVersion = v1alpha1.SchemeGroupVersion
	apiGroupInfo.GroupMeta.GroupVersions = []schema.GroupVersion{
		v1alpha1.SchemeGroupVersion,
	}

	v1alpha1storage := map[string]rest.Storage{}
	r, err := reg.NewInstanceREST(Scheme, serverConfig.RESTOptionsGetter)
	if err != nil {
		return err
	}
	v1alpha1storage["etcdinstances"] = r
	apiGroupInfo.VersionedResourcesStorageMap["v1alpha1"] = v1alpha1storage

	// install APIs and run the api server
	if err = srv.InstallAPIGroup(&apiGroupInfo); err != nil {
		return err
	}
	return srv.PrepareRun().Run(stopCh)
}
