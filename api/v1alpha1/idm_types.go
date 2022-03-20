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

// IDMSpec defines the desired state of IDM.
type IDMSpec struct {
	// +kubebuilder:validation:MaxLength:=64
	// The hostname to be used when installing ipa-server; The default value
	// is composed by <namespace>.<ingressDomain>
	Host string `json:"host,omitempty"`
	// The Realm to be managed by the freeipa instance
	Realm string `json:"realm,omitempty"`
	// The password secret which store the admin and dm passwords
	PasswordSecret *string `json:"passwordSecret"`
	// Resource requirements for the deployment
	Resources corev1.ResourceRequirements `json:"resources"`
	// +optional
	// Volume template for the persistent storage to use
	VolumeClaimTemplate *corev1.PersistentVolumeClaimSpec `json:"volumeClaimTemplate,omitempty"`
}

// IDMStatus defines the observed state of IDM
type IDMStatus struct {
	// The secret name that was used or generated
	SecretName string   `json:"secretName,omitempty"`
	MasterPod  string   `json:"master"`
	ReplicaPod []string `json:"replicas,omitempty"`
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
