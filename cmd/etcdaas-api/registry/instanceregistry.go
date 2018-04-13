package registry

import (
	"fmt"

	iv "github.com/simonferquel/devoxx-2018-k8s-workshop/cmd/etcdaas-api/internalversion"
	"github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/apis/etcdaas/v1alpha1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

func NewInstanceREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*genericregistry.Store, error) {
	strategy := &instanceStrategy{ObjectTyper: scheme, NameGenerator: names.SimpleNameGenerator}

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &v1alpha1.ETCDInstance{} },
		NewListFunc:              func() runtime.Object { return &v1alpha1.ETCDInstanceList{} },
		PredicateFunc:            matchEtcd,
		DefaultQualifiedResource: iv.SchemeGroupVersion.WithResource("etcdinstances").GroupResource(),
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: getEtcdAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return store, nil

}
func getEtcdAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	etcd, ok := obj.(*v1alpha1.ETCDInstance)
	if !ok {
		return nil, nil, false, fmt.Errorf("given object is not an etcd")
	}
	return labels.Set(etcd.ObjectMeta.Labels), etcdToSelectableFields(etcd), etcd.Initializers != nil, nil
}
func matchEtcd(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: getEtcdAttrs,
	}
}
func etcdToSelectableFields(obj *v1alpha1.ETCDInstance) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type instanceStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (s *instanceStrategy) NamespaceScoped() bool {
	return true
}

func (s *instanceStrategy) PrepareForCreate(ctx genericapirequest.Context, obj runtime.Object) {
	o := obj.(*v1alpha1.ETCDInstance)
	fmt.Printf("PrepareForCreate %s spec: %#v, status: %#v\n", o.Name, o.Spec, o.Status)
}

func (s *instanceStrategy) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*v1alpha1.ETCDInstance)
	fmt.Printf("Validate %s spec: %#v, status: %#v\n", o.Name, o.Spec, o.Status)
	return nil
}

func (s *instanceStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
	o := obj.(*v1alpha1.ETCDInstance)
	fmt.Printf("PrepareForUpdate %s spec: %#v, status: %#v\n", o.Name, o.Spec, o.Status)
}

func (s *instanceStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	o := obj.(*v1alpha1.ETCDInstance)
	fmt.Printf("ValidateUpdate %s spec: %#v, status: %#v\n", o.Name, o.Spec, o.Status)

	return nil

}

func (s *instanceStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (s *instanceStrategy) AllowUnconditionalUpdate() bool {
	return true
}

func (s *instanceStrategy) Canonicalize(obj runtime.Object) {
	o := obj.(*v1alpha1.ETCDInstance)
	fmt.Printf("Canonicalize %s spec: %#v, status: %#v\n", o.Name, o.Spec, o.Status)
}
