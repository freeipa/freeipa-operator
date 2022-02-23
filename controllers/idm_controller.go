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

package controllers

import (
	"context"
	generalerr "errors"
	"fmt"

	"github.com/go-logr/logr"
	configv1 "github.com/openshift/api/config/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/freeipa/freeipa-operator/api/v1alpha1"
	arguments "github.com/freeipa/freeipa-operator/internal/arguments"
	manifests "github.com/freeipa/freeipa-operator/manifests"
)

// IDMReconciler reconciles a IDM object
type IDMReconciler struct {
	client.Client

	Log           logr.Logger
	Scheme        *runtime.Scheme
	IngressDomain string
	Arguments     *arguments.Arguments
	WorkloadImage string
}

var (
	metricsAddr string
)

// ReadIngressDomainFromOpenshiftConfig This method encapsulate the operation below:
//
// ```raw
// oc get dnses.config.openshift.io/cluster -o json | jq -r ".spec.baseDomain"
// ```
//
// So that it returns the baseDomain that was configured meanwhile installing the
// cluster. This value cannot be changed once the cluster has been installed.
// ctx A memory context used for the operation.
// Return If the operation is executed successfully, the baseDomain and nil for
// error, else return empty string for baseDomain and the error object.
func (r *IDMReconciler) ReadIngressDomainFromOpenshiftConfig(ctx context.Context) (string, error) {
	namespacedName := types.NamespacedName{
		Namespace: "",
		Name:      "cluster",
	}
	ingressConfig := &configv1.Ingress{}
	if err := r.Get(ctx, namespacedName, ingressConfig); err != nil {
		return "", err
	}
	return ingressConfig.Spec.Domain, nil
}

// InitIngressDomain Initialize the cache for the IngressDomain that is
// used by the cluster.
// ctx The memory context to be used for the operation.
// Return nil if it was initialized, else an error object.
func (r *IDMReconciler) InitIngressDomain(ctx context.Context) error {
	var err error
	if r.IngressDomain == "" {
		log := r.Log.WithValues("idm_controller", "InitBaseDomain")
		log.Info("IngressDomain is empty, retrieving")

		if r.IngressDomain, err = r.ReadIngressDomainFromOpenshiftConfig(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Reconcile Read the current of the cluster for IDM object and makes the
// necessary changes to bring the system to the requested state.
// +kubebuilder:rbac:groups=idmocp.redhat.com,resources=idms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=idmocp.redhat.com,resources=idms/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;create;update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;create;update;watch
func (r *IDMReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	var idm v1alpha1.IDM = v1alpha1.IDM{}
	log := r.Log.WithValues("idm", req.NamespacedName)

	// Fetch the IDM instance
	err = r.Get(ctx, req.NamespacedName, &idm)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("IDM resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get IDM")
		return ctrl.Result{}, err
	}

	if err := r.InitIngressDomain(ctx); err != nil {
		log.Error(err, "Failed initializing the IngressDomain attribute")
		return ctrl.Result{}, err
	}

	// FIXME Remove this after creating the specific ServiceAccount for the workload
	if err := r.CreateServiceAccount(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	// FIXME Remove this after creating the specific ServiceAccount for the workload
	if err := r.CreateRole(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	// FIXME Remove this after creating the specific ServiceAccount for the workload
	if err := r.CreateRoleBinding(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateSecret(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateStatefulsetMain(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateWebService(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateRoute(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// CreateRoleBinding Create the role
// TODO Remove this method when the RoleBinding is decoupled
func (r *IDMReconciler) CreateRoleBinding(ctx context.Context, item *v1alpha1.IDM) error {
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetRoleBindingName(item),
	}
	log := r.Log.WithValues("idm", namespacedName)
	found := &rbacv1.RoleBinding{}
	err := r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating RoleBinding")
			manifest := manifests.RoleBindingForIDM(item)
			if err := r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		log.Info("Currently the RoleBinding exists")
	}

	return nil
}

// CreateRole Create the role
// TODO Remove this method when the Role is decoupled
func (r *IDMReconciler) CreateRole(ctx context.Context, item *v1alpha1.IDM) error {
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetRoleName(item),
	}
	log := r.Log.WithValues("idm", namespacedName)
	found := &rbacv1.Role{}
	err := r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Role")
			manifest := manifests.RoleForIDM(item)
			if err := r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		log.Info("Currently the Role exists")
	}

	return nil
}

// CreateServiceAccount Create the service account
// TODO Remove this method when the ServiceAccount is decoupled
func (r *IDMReconciler) CreateServiceAccount(ctx context.Context, item *v1alpha1.IDM) error {
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetServiceAccountName(item),
	}
	log := r.Log.WithValues("idm", namespacedName)
	found := &corev1.ServiceAccount{}
	err := r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Service Account")
			manifest := manifests.ServiceAccountForIDM(item)
			if err := r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		log.Info("Currently the ServiceAccount exists")
	}

	return nil
}

// CheckStatusSecret Check the status indicated into the secret exists or not.
func (r *IDMReconciler) CheckStatusSecret(ctx context.Context, item *v1alpha1.IDM) error {
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      item.Status.SecretName,
	}
	found := &corev1.Secret{}
	if err := r.Get(ctx, namespacedName, found); err != nil {
		return err
	}
	if _, ok := found.Data["IPA_ADMIN_PASSWORD"]; !ok {
		return fmt.Errorf("IPA_ADMIN_PASSWORD does not exist in '%s' secret", item.Status.SecretName)
	}
	if _, ok := found.Data["IPA_DM_PASSWORD"]; !ok {
		return fmt.Errorf("IPA_DM_PASSWORD does not exist in '%s' secret", item.Status.SecretName)
	}
	return r.Get(ctx, namespacedName, found)
}

// UpdateStatusSecretNameWith Update the secretName status field
func (r *IDMReconciler) UpdateStatusSecretNameWith(ctx context.Context, secretName string, item *v1alpha1.IDM) error {
	item.Status.SecretName = secretName
	if err := r.Status().Update(ctx, item); err != nil {
		return err
	}
	return nil
}

// CreateSecret Create a secret for he encrypted information
// ctx is the context for this request.
// item is a reference to the custom resource to sync to.
// Return no nil error when any situation happened, else it returns nil.
func (r *IDMReconciler) CreateSecret(ctx context.Context, item *v1alpha1.IDM) error {
	var err error
	log := r.Log.WithValues("idm", item.Namespace)

	// If it was assigned it returns check status secret
	if item.Status.SecretName != "" {
		log.Info("Checking Current Secret Name")
		if err = r.CheckStatusSecret(ctx, item); err != nil {
			return err
		}
	}

	if item.Spec.PasswordSecret != nil {
		namespacedName := types.NamespacedName{
			Namespace: item.Namespace,
			Name:      *item.Spec.PasswordSecret,
		}
		found := &corev1.Secret{}
		err = r.Get(ctx, namespacedName, found)
		if err == nil {
			log.Info("Updating secret name in status")
			if err = r.UpdateStatusSecretNameWith(ctx, manifests.GetSecretName(item), item); err != nil {
				return err
			}
			return nil
		}
		if !errors.IsNotFound(err) {
			return err
		}
	}

	// If a password was not specified, it creates a secret with a random password
	log.Info("Creating Secret with random password")
	manifest := manifests.SecretForIDM(item, manifests.GenerateRandomPassword(), manifests.GenerateRandomPassword())
	ctrl.SetControllerReference(item, manifest, r.Scheme)
	if err = r.Create(ctx, manifest); err != nil {
		return err
	}
	if err = r.UpdateStatusSecretNameWith(ctx, manifests.GetSecretName(item), item); err != nil {
		return err
	}
	return nil
}

// CreatePersistentVolumeClaim Create the volum claim that will be used by the worload.
func (r *IDMReconciler) CreatePersistentVolumeClaim(ctx context.Context, item *v1alpha1.IDM) error {
	var err error
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetMainPersistentVolumeClaimName(item),
	}
	log := r.Log.WithValues(item.Name, namespacedName)
	if item.Spec.VolumeClaimTemplate == nil {
		if r.Arguments.GetDefaultStorage() == "ephemeral" {
			log.Info("VolumeTemplateClaim not defined, using an 'ephemeral' volume for storing data")
			return nil
		}
		return generalerr.New("VolumeClaimTemplate is missed")
	}
	found := &corev1.PersistentVolumeClaim{}
	err = r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Volume Claim")
			manifest := manifests.MainPersistentVolumeClaimForIDM(item)
			ctrl.SetControllerReference(item, manifest, r.Scheme)
			if err = r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// TODO Update changes if any that affect to the PersistentVolumeClaim
		log.Info("Currently the PersistentVolumeClaim exists")
	}

	return nil
}

// CreateMainPod Create the master freeipa pod
func (r *IDMReconciler) CreateMainPod(ctx context.Context, item *v1alpha1.IDM) error {
	var err error
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetMainPodName(item),
	}
	log := r.Log.WithValues("idm", namespacedName)

	var defaultStorage = r.Arguments.GetDefaultStorage()
	// Check volume storage information
	err = manifests.CheckVolumeInformation(item, defaultStorage)
	if err != nil {
		log.Info("Checking Volume Information")
		return err
	}

	found := &corev1.Pod{}
	err = r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Master Pod")
			manifest := manifests.MainPodForIDM(item, r.IngressDomain, r.WorkloadImage, defaultStorage)
			ctrl.SetControllerReference(item, manifest, r.Scheme)
			if err = r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// TODO Update changes if any that affect to the Pod
		log.Info("Currently the Main Pod exists")
	}
	return nil
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// CreateStatefulsetMain Create the Statefulset object for the workload
func (r *IDMReconciler) CreateStatefulsetMain(ctx context.Context, item *v1alpha1.IDM) error {
	var err error
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetMainPodName(item),
	}
	log := r.Log.WithValues("idm", namespacedName)

	var defaultStorage = r.Arguments.GetDefaultStorage()

	found := &appsv1.StatefulSet{}
	err = r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Main Statefulset")
			manifest := manifests.MainStatefulsetForIDM(item, r.IngressDomain, r.WorkloadImage, defaultStorage)
			ctrl.SetControllerReference(item, manifest, r.Scheme)
			if err = r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// TODO Update changes if any that affect to the Pod
		log.Info("Currently the Main Pod exists")
	}
	return nil
}

// CreateWebService Create the service to access the web frontend running on Apache
func (r *IDMReconciler) CreateWebService(ctx context.Context, item *v1alpha1.IDM) error {
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetWebServiceName(item),
	}
	log := r.Log.WithValues("idm", namespacedName)
	found := &corev1.Service{}
	err := r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Service for Web access")
			manifest := manifests.ServiceWebForIDM(item)
			ctrl.SetControllerReference(item, manifest, r.Scheme)
			if err := r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// TODO Update changes if any that affect to the Pod
		log.Info("Currently the Service for Web Interface exists")
	}

	return nil
}

// CreateRoute Create the service to access the web frontend running on Apache
func (r *IDMReconciler) CreateRoute(ctx context.Context, item *v1alpha1.IDM) error {
	var err error
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      item.Name,
	}
	log := r.Log.WithValues("idm", namespacedName)
	found := &routev1.Route{}
	err = r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Route to web service")
			manifest := manifests.RouteForIDM(item, manifests.GenerateDefaultRoute(item.Namespace, r.IngressDomain))
			ctrl.SetControllerReference(item, manifest, r.Scheme)
			if err = r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// TODO Update changes if any that affect to the Pod
		log.Info("Currently the Route exists")
	}

	return nil
}

// SetupWithManager Specifies how the controller is built to watch a CR and
// other resources that are owned and managed by that controller.
func (r *IDMReconciler) SetupWithManager(mgr ctrl.Manager, args *arguments.Arguments) error {
	// A build pattern is used here, so that the controller
	// is not 100% initialized until the Complete method has
	// finished.
	r.Arguments = args
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.IDM{}).
		Complete(r)
}
