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

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/freeipa/freeipa-operator/api/v1alpha1"
	manifests "github.com/freeipa/freeipa-operator/manifests"
)

// IDMReconciler reconciles a IDM object
type IDMReconciler struct {
	client.Client

	Log        logr.Logger
	Scheme     *runtime.Scheme
	BaseDomain string
}

var (
	metricsAddr string
)

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

	if err := r.CreateServiceAccount(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateRole(ctx, &idm); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.CreateRoleBinding(ctx, &idm); err != nil {
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

// CreateMainPod Create the master freeipa pod
func (r *IDMReconciler) CreateMainPod(ctx context.Context, item *v1alpha1.IDM) error {
	var err error
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetMainPodName(item),
	}
	log := r.Log.WithValues("idm", namespacedName)
	found := &corev1.Pod{}
	err = r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Master Pod")
			manifest := manifests.MainPodForIDM(item, r.BaseDomain)
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
func (r *IDMReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// A build pattern is used here, so that the controller
	// is not 100% initialized until the Complete method has
	// finished.
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.IDM{}).
		Complete(r)
}
