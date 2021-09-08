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

package main

import (
	"os"

	configv1 "github.com/openshift/api/config/v1"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	idmv1alpha1 "github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/freeipa/freeipa-operator/controllers"
	arguments "github.com/freeipa/freeipa-operator/internal/arguments"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const (
	ENV_DEFAULT_STORAGE = "DEFAULT_STORAGE"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(idmv1alpha1.AddToScheme(scheme))
	utilruntime.Must(routev1.AddToScheme(scheme))
	utilruntime.Must(configv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func getWorkloadImage() string {
	if workload, exist := os.LookupEnv("WORKLOAD_IMAGE"); exist {
		return workload
	}
	return "quay.io/freeipa/freeipa-openshift-container:freeipa-server"
}

func main() {
	var err error
	var ctrlArguments *arguments.Arguments
	var workloadImage string = getWorkloadImage()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Load and check arguments
	ctrlArguments, err = arguments.New()
	if err != nil {
		setupLog.Error(err, "invalid controller arguments")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: ctrlArguments.GetMetricsAddr(),
		Port:               9443,
		LeaderElection:     ctrlArguments.GetEnableLeaderElection(),
		LeaderElectionID:   "42b6c26c.redhat.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(2)
	}

	if err = (&controllers.IDMReconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName("controllers").WithName("IDM"),
		Scheme:        mgr.GetScheme(),
		WorkloadImage: workloadImage,
	}).SetupWithManager(mgr, ctrlArguments); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "IDM")
		os.Exit(3)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(4)
	}
}
