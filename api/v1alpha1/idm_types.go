/*
Copyright 2020 Red Hat.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IDMSpec defines the desired state of IDM
type IDMSpec struct {
	// +kubebuilder:validation:MaxLength:=64
	Host           string                      `json:"host,omitempty"`
	Realm          string                      `json:"realm,omitempty"`
	PasswordSecret *string                     `json:"passwordSecret"`
	Resources      corev1.ResourceRequirements `json:"resources"`
	// +optional
	VolumeClaimTemplate *corev1.PersistentVolumeClaimSpec `json:"volumeClaimTemplate,omitempty"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of IDM. Edit idm_types.go to remove/update
}

// IDMStatus defines the observed state of IDM
type IDMStatus struct {
	SecretName string   `json:"secretName,omitempty"`
	MasterPod  string   `json:"master"`
	ReplicaPod []string `json:"replicas,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=idms,scope=Namespaced

// IDM is the Schema for the idms API
type IDM struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IDMSpec   `json:"spec,omitempty"`
	Status IDMStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IDMList contains a list of IDM
type IDMList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IDM `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IDM{}, &IDMList{})
}
