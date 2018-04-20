package registry

import (
	"fmt"
	"net/http"

	"github.com/simonferquel/devoxx-2018-k8s-workshop/pkg/apis/etcdaas/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

func NewBackupREST() rest.Storage {
	return &backupREST{}
}

type backupREST struct {
}

func (b *backupREST) New() runtime.Object {
	return &v1alpha1.ETCDInstance{}
}

// Connect returns an http.Handler that will handle the request/response for a given API invocation.
// The provided responder may be used for common API responses. The responder will write both status
// code and body, so the ServeHTTP method should exit after invoking the responder. The Handler will
// be used for a single API request and then discarded. The Responder is guaranteed to write to the
// same http.ResponseWriter passed to ServeHTTP.
func (b *backupREST) Connect(ctx genericapirequest.Context, id string, options runtime.Object, r rest.Responder) (http.Handler, error) {
	ns, _ := genericapirequest.NamespaceFrom(ctx)
	fmt.Printf("received backup connect request for %s in namespace %s", id, ns)
	return &backupHandler{id: id, ns: ns}, nil
}

// NewConnectOptions returns an empty options object that will be used to pass
// options to the Connect method. If nil, then a nil options object is passed to
// Connect. It may return a bool and a string. If true, the value of the request
// path below the object will be included as the named string in the serialization
// of the runtime object.
func (b *backupREST) NewConnectOptions() (runtime.Object, bool, string) {
	return nil, false, ""
}

// ConnectMethods returns the list of HTTP methods handled by Connect
func (b *backupREST) ConnectMethods() []string {
	return []string{http.MethodGet}
}

var _ rest.Storage = &backupREST{}
var _ rest.Connecter = &backupREST{}

type backupHandler struct {
	id string
	ns string
}

func (h *backupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// **** Serve me
}
