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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	idmv1alpha1 "github.com/freeipa/freeipa-operator/api/v1alpha1"
)

// IDMReconciler reconciles a IDM object
type IDMReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile Read the current of the cluster for IDM object and makes the
// necessary changes to bring the system to the requested state.
// +kubebuilder:rbac:groups=idmocp.redhat.com,resources=idms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=idmocp.redhat.com,resources=idms/status,verbs=get;update;patch
func (r *IDMReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("idm", req.NamespacedName)

	// Fetch the IDM instance
	idm := &idmv1alpha1.IDM{}
	err := r.Get(ctx, req.NamespacedName, idm)
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

	found := &corev1.Pod{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      idm.Name,
		Namespace: idm.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new pod
		pod := r.podForIDM(idm)
		log.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.Create(ctx, pod)
		if err != nil {
			log.Error(err, "Failed to create new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			return ctrl.Result{}, err
		}
		// Pod created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Pod")
		return ctrl.Result{}, err
	}

	// Check if the deploymen already exists, if not create a new one
	// found := &appsv1.Deployment{}
	// err = r.Get(ctx, types.NamespacedName{Name: idm.Name, Namespace: idm.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	// Define a new deployment
	// 	dep := r.deploymentForIDM(idm)
	// 	log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	// 	err = r.Create(ctx, dep)
	// 	if err != nil {
	// 		log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	// 		return ctrl.Result{}, err
	// 	}
	// 	// Deployment created successfully - return and requeue
	// 	return ctrl.Result{Requeue: true}, nil
	// } else if err != nil {
	// 	log.Error(err, "Failed to get Deployment")
	// 	return ctrl.Result{}, err
	// }

	// TODO Implement here the changes to bring the spec to the
	//      requested state

	// Update the IDM status with the pod names
	// List the pods for this idm's deployment
	// podList := &corev1.PodList{}
	// listOpts := []client.ListOption{
	// 	client.InNamespace(idm.Namespace),
	// 	client.MatchingLabels(IDM(idm.Name)),
	// }
	// if err = r.List(ctx, podList, listOpts...); err != nil {
	// 	log.Error(err, "Failed to list pods", "idm.Namespace", idm.Namespace, "idm.Name", idm.Name)
	// 	return ctrl.Result{}, err
	// }
	// podNames := getPodNames(podList.Items)

	return ctrl.Result{}, nil
}

// podForIDM return a pod for IDM
func (r *IDMReconciler) podForIDM(m *idmv1alpha1.IDM) *corev1.Pod {
	ls := labelsForIDM(m.Name)

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "freeipa-master",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}

	ctrl.SetControllerReference(m, pod, r.Scheme)
	return pod
}

// deploymentForIDM returns a idm Deployment object
func (r *IDMReconciler) deploymentForIDM(m *idmv1alpha1.IDM) *appsv1.Deployment {
	ls := labelsForIDM(m.Name)
	// realm := m.Spec.Realm

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "freeipa-operator:dev-test",
						Name:    "freeipa",
						Command: []string{"sleep 3600"},
					}},
				},
			},
		},
	}
	// Set IDM instace as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// labelsForIDM returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForIDM(name string) map[string]string {
	return map[string]string{"app": "idm", "idm_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager Specifies how the controller is built to watch a CR and
// other resources that are owned and managed by that controller.
func (r *IDMReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&idmv1alpha1.IDM{}).
		Complete(r)
}
