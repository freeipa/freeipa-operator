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
	"fmt"
	"os"
	"strings"

	configv1 "github.com/openshift/api/config/v1"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	idmv1alpha1 "github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/freeipa/freeipa-operator/controllers"
	arguments "github.com/freeipa/freeipa-operator/internal/arguments"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(idmv1alpha1.AddToScheme(scheme))
	utilruntime.Must(routev1.AddToScheme(scheme))
	utilruntime.Must(configv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

// getWatchNamespace returns the Namespace the operator should be watching for changes
func getWatchNamespace() (string, error) {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	const watchNamespaceEnvVar = "WATCH_NAMESPACE"

	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}
	return ns, nil
}

// getEnableWebhooks returns the value for ENABLE_WEBHOOKS env var or
// true by default.
func getEnableWebhooks() (bool, error) {
	const enableWebhooksEnvVar = "ENABLE_WEBHOOKS"

	enableWebhooks, found := os.LookupEnv(enableWebhooksEnvVar)
	if !found {
		return true, nil
	}
	if strings.ToLower(enableWebhooks) != "false" {
		return true, nil
	}
	return false, nil
}

// getWorkloadImage return t
func getWorkloadImage() (string, error) {
	if workload, exist := os.LookupEnv("RELATED_IMAGE_FREEIPA"); exist && workload != "" {
		return workload, nil
	}
	return "", fmt.Errorf("RELATED_IMAGE_FREEIPA environment variable must be specified and not empty")
}

func main() {
	var err error
	var ctrlArguments *arguments.Arguments
	var workloadImage string

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Load and check arguments
	ctrlArguments, err = arguments.New()
	if err != nil {
		setupLog.Error(err, "invalid controller arguments")
		os.Exit(1)
	}

	// Get RELATED_IMAGE_FREEIPA
	workloadImage, err = getWorkloadImage()
	if err != nil {
		setupLog.Error(err, "Reading RELATED_IMAGE_FREEIPA")
		os.Exit(2)
	}

	// Get WATCH_NAMESPACE value
	watchNamespace, err := getWatchNamespace()
	if err != nil || watchNamespace == "" {
		setupLog.Error(err, "unable to get WatchNamespace, "+
			"the manager will watch and manage resources in all namespaces")
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     ctrlArguments.GetMetricsAddr(),
		Port:                   9443,
		HealthProbeBindAddress: ctrlArguments.GetProbeAddr(),
		LeaderElection:         ctrlArguments.GetEnableLeaderElection(),
		LeaderElectionID:       "42b6c26c.redhat.com",
		Namespace:              watchNamespace, // namespaced-scope when the value is not an empty string
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
	// Conditionally disable webhooks, this is useful when debugging locally
	if enable, _ := getEnableWebhooks(); enable {
		if err = (&idmv1alpha1.IDM{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "IDM")
			os.Exit(4)
		}
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(5)
	}
}
