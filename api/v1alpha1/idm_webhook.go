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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var idmlog = logf.Log.WithName("idm-resource")

func (r *IDM) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-idmocp-redhat-com-redhat-com-v1alpha1-idm,mutating=true,failurePolicy=fail,groups=idmocp.redhat.com.redhat.com,resources=idms,verbs=create;update,versions=v1alpha1,name=midm.kb.io

var _ webhook.Defaulter = &IDM{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *IDM) Default() {
	idmlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:verbs=create;update,path=/validate-idmocp-redhat-com-redhat-com-v1alpha1-idm,mutating=false,failurePolicy=fail,groups=idmocp.redhat.com.redhat.com,resources=idms,versions=v1alpha1,name=vidm.kb.io

var _ webhook.Validator = &IDM{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *IDM) ValidateCreate() error {
	idmlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *IDM) ValidateUpdate(old runtime.Object) error {
	idmlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *IDM) ValidateDelete() error {
	idmlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
