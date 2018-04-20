package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ETCDInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ETCDInstanceSpec   `json:"spec,omitempty"`
	Status            ETCDInstanceStatus `json:"status,omitempty"`
}

type ETCDInstanceSpec struct {
	// *** Fill me
}

type ETCDInstanceState string

const (
	ETCDNone      = ETCDInstanceState("")
	ETCDDeploying = ETCDInstanceState("deploying")
	ETCDRunning   = ETCDInstanceState("running")
	ETCDFailed    = ETCDInstanceState("failed")
)

type ETCDInstanceStatus struct {
	State   ETCDInstanceState `json:"state,omitempty"`
	Message string            `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ETCDInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ETCDInstance `json:"items"`
}
