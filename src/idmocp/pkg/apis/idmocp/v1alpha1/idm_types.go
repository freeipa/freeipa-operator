package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IDMSpec defines the desired state of IDM
type IDMSpec struct {
	Realm string `json:"realm"`
}

// IDMStatus defines the observed state of IDM
type IDMStatus struct {
	Servers []string `json:"servers"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IDM is the Schema for the idms API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=idms,scope=Namespaced
type IDM struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IDMSpec   `json:"spec,omitempty"`
	Status IDMStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IDMList contains a list of IDM
type IDMList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IDM `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IDM{}, &IDMList{})
}
