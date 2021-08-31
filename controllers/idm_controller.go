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

	"github.com/go-logr/logr"
	configv1 "github.com/openshift/api/config/v1"
	routev1 "github.com/openshift/api/route/v1"
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

	Log        logr.Logger
	Scheme     *runtime.Scheme
	BaseDomain string
	Arguments  *arguments.Arguments
}

var (
	metricsAddr string
)

// ReadBaseDomainFromOpenshiftConfig This method encapsulate the operation below:
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
func (r *IDMReconciler) ReadBaseDomainFromOpenshiftConfig(ctx context.Context) (string, error) {
	namespacedName := types.NamespacedName{
		Namespace: "",
		Name:      "cluster",
	}
	dnsConfig := &configv1.DNS{}
	if err := r.Get(ctx, namespacedName, dnsConfig); err != nil {
		return "", err
	}
	return dnsConfig.Spec.BaseDomain, nil
}

// InitBaseDomain Initialize the cache for the BaseDomain that is
// used by the cluster.
// ctx The memory context to be used for the operation.
// Return nil if it was initialized, else an error object.
func (r *IDMReconciler) InitBaseDomain(ctx context.Context) error {
	var err error
	if r.BaseDomain == "" {
		log := r.Log.WithValues("idm_controller", "GetClusterDomain")
		log.Info("BaseDomain is empty, retrieving")

		if r.BaseDomain, err = r.ReadBaseDomainFromOpenshiftConfig(ctx); err != nil {
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

	if err := r.InitBaseDomain(ctx); err != nil {
		log.Error(err, "Failed initializing the BaseDomain attribute")
		return ctrl.Result{}, err
	}

	if err := r.CreateServiceAccount(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateRole(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateRoleBinding(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateSecret(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreatePersistentVolumeClaim(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateMainPod(ctx, &idm); err != nil {
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
	return r.Get(ctx, namespacedName, found)
}

// UpdateStatusSecretNameWith Update the secretName status field
func (r *IDMReconciler) UpdateStatusSecretNameWith(secretName string, ctx context.Context, item *v1alpha1.IDM) error {
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
			if err = r.UpdateStatusSecretNameWith(manifests.GetSecretName(item), ctx, item); err != nil {
				return err
			}
			return nil
		}
		if !errors.IsNotFound(err) {
			return err
		}
		log.Info("Creating Secret")
		manifest := manifests.SecretForIDM(item, manifests.GenerateRandomPassword())
		ctrl.SetControllerReference(item, manifest, r.Scheme)
		if err = r.Create(ctx, manifest); err != nil {
			return err
		}
		if err = r.UpdateStatusSecretNameWith(manifests.GetSecretName(item), ctx, item); err != nil {
			return err
		}
		return nil
	}

	// If a password was not specified, it creates a secret with a random password
	log.Info("Creating Secret with random password")
	manifest := manifests.SecretForIDM(item, manifests.GenerateRandomPassword())
	ctrl.SetControllerReference(item, manifest, r.Scheme)
	if err = r.Create(ctx, manifest); err != nil {
		return err
	}
	if err = r.UpdateStatusSecretNameWith(manifests.GetSecretName(item), ctx, item); err != nil {
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
			manifest := manifests.MainPodForIDM(item, r.BaseDomain, defaultStorage)
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
			manifest := manifests.RouteForIDM(item, r.BaseDomain)
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
