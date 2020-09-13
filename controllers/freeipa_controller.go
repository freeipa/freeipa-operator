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

	freeipav1alpha1 "github.com/freeipa/freeipa-operator/api/v1alpha1"
)

// FreeipaReconciler reconciles a Freeipa object
type FreeipaReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile Read the current of the cluster for Freeipa object and makes the
// necessary changes to bring the system to the requested state.
// +kubebuilder:rbac:groups=freeipa.redhat.com,resources=freeipas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=freeipa.redhat.com,resources=freeipas/status,verbs=get;update;patch
func (r *FreeipaReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("freeipa", req.NamespacedName)

	// Fetch the Freeipa instance
	freeipa := &freeipav1alpha1.Freeipa{}
	err := r.Get(ctx, req.NamespacedName, freeipa)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Freeipa resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Freeipa")
		return ctrl.Result{}, err
	}

	found := &corev1.Pod{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      freeipa.Name,
		Namespace: freeipa.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new pod
		pod := r.podForFreeipa(freeipa)
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
	// err = r.Get(ctx, types.NamespacedName{Name: freeipa.Name, Namespace: freeipa.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	// Define a new deployment
	// 	dep := r.deploymentForFreeipa(freeipa)
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

	// Update the Freeipa status with the pod names
	// List the pods for this freeipa's deployment
	// podList := &corev1.PodList{}
	// listOpts := []client.ListOption{
	// 	client.InNamespace(freeipa.Namespace),
	// 	client.MatchingLabels(labelsForFreeipa(freeipa.Name)),
	// }
	// if err = r.List(ctx, podList, listOpts...); err != nil {
	// 	log.Error(err, "Failed to list pods", "Freeipa.Namespace", freeipa.Namespace, "Freeipa.Name")
	// 	return ctrl.Result{}, err
	// }
	// podNames := getPodNames(podList.Items)

	return ctrl.Result{}, nil
}

// podForFreeipa return a pod for Freeipa
func (r *FreeipaReconciler) podForFreeipa(m *freeipav1alpha1.Freeipa) *corev1.Pod {
	ls := labelsForFreeipa(m.Name)

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

// deploymentForFreeipa returns a freeipa Deployment object
func (r *FreeipaReconciler) deploymentForFreeipa(m *freeipav1alpha1.Freeipa) *appsv1.Deployment {
	ls := labelsForFreeipa(m.Name)
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
						Command: []string{"freeipa"},
					}},
				},
			},
		},
	}
	// Set Freeipa instace as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// labelsForFreeipa returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForFreeipa(name string) map[string]string {
	return map[string]string{"app": "freeipa", "freeipa_cr": name}
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
func (r *FreeipaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&freeipav1alpha1.Freeipa{}).
		Complete(r)
}
