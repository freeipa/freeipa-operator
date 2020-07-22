package idm

import (
	"context"
	"reflect"

	idmocpv1alpha1 "idmocp/pkg/apis/idmocp/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_idm")

// Add creates a new IDM Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileIDM{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("idm-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource IDM
	err = c.Watch(&source.Kind{Type: &idmocpv1alpha1.IDM{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Our secondary resources are the Pods running the freeipa
	// server container.
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &idmocpv1alpha1.IDM{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileIDM implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileIDM{}

// ReconcileIDM reconciles a IDM object
type ReconcileIDM struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a IDM object and makes changes based on the state read
// and what is in the IDM.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
func (r *ReconcileIDM) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling IDM")

	// Fetch the IDM instance
	instance := &idmocpv1alpha1.IDM{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// List all pods owned by this PodSet instance
	podList := &corev1.PodList{}
	lbls := map[string]string{
		"app":     instance.Name,
	}
	listOps := &client.ListOptions{
		Namespace: instance.Namespace,
		LabelSelector: labels.SelectorFromSet(lbls),
	}
	if err = r.client.List(context.TODO(), podList, listOps); err != nil {
		return reconcile.Result{}, err
	}

	// Count the pods that are pending or running as available
	var pods []corev1.Pod
	for _, pod := range podList.Items {
		if pod.ObjectMeta.DeletionTimestamp != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending {
			pods = append(pods, pod)
		}
	}
	numPods := int32(len(pods))
	podNames := []string{}
	for _, pod := range pods {
		podNames = append(podNames, pod.ObjectMeta.Name)
	}

	// Update the status if necessary
	status := idmocpv1alpha1.IDMStatus{
		Servers: podNames,
	}
	if !reflect.DeepEqual(instance.Status, status) {
		instance.Status = status
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update IDM status")
			return reconcile.Result{}, err
		}
	}

	if numPods < 1 {
		reqLogger.Info("Deploying IDM")

		pod := newPodForCR(instance)
		if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			reqLogger.Error(err, "Failed to create pod", "pod.name", pod.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

func newPodForCR(cr *idmocpv1alpha1.IDM) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	// TODO create PersistentVolumeClaim?
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cr.Name + "-pod",  // adds random suffix
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "freeipa",  // dns name?
					Image:   "freeipa/freeipa-server:fedora-31",
					Command: []string{"sleep", "3600"}, // FIXME
				},
			},
		},
	}
}
