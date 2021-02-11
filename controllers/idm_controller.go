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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	// k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/freeipa/freeipa-operator/api/v1alpha1"
	manifests "github.com/freeipa/freeipa-operator/manifests"
	// go get k8s.io/client-go@v0.20.0
)

// IDMReconciler reconciles a IDM object
type IDMReconciler struct {
	client.Client

	Log    logr.Logger
	Scheme *runtime.Scheme
}

var (
	metricsAddr string
)

// Reconcile Read the current of the cluster for IDM object and makes the
// necessary changes to bring the system to the requested state.
// +kubebuilder:rbac:groups=idmocp.redhat.com,resources=idms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=idmocp.redhat.com,resources=idms/status,verbs=get;update;patch
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

	if err := r.CreateMasterPod(ctx, &idm); err != nil {
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

// CreateMasterPod Create the master freeipa pod
func (r *IDMReconciler) CreateMasterPod(ctx context.Context, item *v1alpha1.IDM) error {
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      manifests.GetMasterPodName(item),
	}
	log := r.Log.WithValues("idm", namespacedName)
	found := &corev1.Pod{}
	err := r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Master Pod")
			manifest := manifests.MasterPodForIDM(item)
			ctrl.SetControllerReference(item, manifest, r.Scheme)
			if err := r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// TODO Update changes if any that affect to the Pod
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
	}

	return nil
}

// CreateRoute Create the service to access the web frontend running on Apache
func (r *IDMReconciler) CreateRoute(ctx context.Context, item *v1alpha1.IDM) error {
	namespacedName := types.NamespacedName{
		Namespace: item.Namespace,
		Name:      item.Name,
	}
	log := r.Log.WithValues("idm", namespacedName)
	found := &routev1.Route{}
	err := r.Get(ctx, namespacedName, found)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Route to web service")
			manifest := manifests.RouteForIDM(item)
			ctrl.SetControllerReference(item, manifest, r.Scheme)
			if err := r.Create(ctx, manifest); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// TODO Update changes if any that affect to the Pod
	}

	return nil
}

// SetupWithManager Specifies how the controller is built to watch a CR and
// other resources that are owned and managed by that controller.
func (r *IDMReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.IDM{}).
		Complete(r)
}
