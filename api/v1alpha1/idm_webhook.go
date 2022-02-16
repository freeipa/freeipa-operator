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
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var (
	idmlog     = logf.Log.WithName("idm-resource")
	clientInst client.Client
)

// SetupWebhookWithManager Set up webhooks for the
// indicated controller manager.
// mgr The controller manager where to add the webhook.
func (r *IDM) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-idmocp-redhat-com-v1alpha1-idm,mutating=true,failurePolicy=fail,groups=idmocp.redhat.com,resources=idms,verbs=create;update,versions=v1alpha1,sideEffects=None,admissionReviewVersions=v1,name=midm.kb.io

var _ webhook.Defaulter = &IDM{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *IDM) Default() {
	idmlog.Info("default", "name", r.Name)

	// TODO(user): set here the default values for not specified fields or
	//             empty fields, such as, use a REALM that match the cluster
	//             base domain if empty, or a hostname that match the namespace
	//    		   and the ingressDomain by default or the minimal resource
	//             limits for running the workload.
}

//+kubebuilder:webhook:verbs=create;update,path=/validate-idmocp-redhat-com-v1alpha1-idm,mutating=false,failurePolicy=fail,groups=idmocp.redhat.com,resources=idms,versions=v1alpha1,sideEffects=None,admissionReviewVersions=v1,name=vidm.kb.io

var _ webhook.Validator = &IDM{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
// For this CRD we checks the below:
// - The referenced secret exists and is immutable.
func (r *IDM) ValidateCreate() error {
	idmlog.Info("validate create", "name", r.Name)

	// TODO(user): add here validation when creating the custom resource
	//             such as values for attributes belong to their domains,
	//             or checking that referenced resources exists or
	//             any other check to warranty the values for the custom
	//             resource are valids.

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
// This method validate that the new resource is not changed for
// the fields considered immutable.
// - Realm can not be changed once the IDM has been created.
// - Resources can not be changed as the PodSpec does not allow it.
// - VolumeClaimTemplate can not be changes once the IDM resource
//   has been created.
func (r *IDM) ValidateUpdate(oldRaw runtime.Object) error {
	idmlog.Info("validate update", "name", r.Name)

	old, ok := oldRaw.(*IDM)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected an IDM but got a %T", oldRaw))
	}

	if old.Spec.Realm != r.Spec.Realm {
		return apierrors.NewBadRequest("IDM.Spec.Realm is immutable")
	}

	if !reflect.DeepEqual(r.Spec.PasswordSecret, old.Spec.PasswordSecret) {
		return apierrors.NewBadRequest("IDM.Spec.PasswordSecret is immutable")
	}

	if !reflect.DeepEqual(r.Spec.VolumeClaimTemplate, old.Spec.VolumeClaimTemplate) {
		return apierrors.NewBadRequest("IDM.Spec.VolumeClaimTemplate is immutable")
	}

	if !reflect.DeepEqual(r.Spec.Resources, old.Spec.Resources) {
		return apierrors.NewBadRequest("IDM.Spec.Resources is immutable")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *IDM) ValidateDelete() error {
	idmlog.Info("validate delete", "name", r.Name)

	// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.

	// TODO(user): add here any validation to be made before to proceed to
	//             delete the custom resource
	return nil
}
