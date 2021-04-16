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

package manifests_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	// "k8s.io/client-go/rest"

	// "sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

// var cfg *rest.Config
// var k8sClient client.Client
// var k8sManager ctrl.Manager
// var testEnv *envtest.Environment
// var reconciler *IDMReconciler

// func GetK8sManager() ctrl.Manager {
// 	return k8sManager
// }

// func GetReconciler() *IDMReconciler {
// 	return reconciler
// }

func TestManifests(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter(fmt.Sprintf("../junit/TEST-ginkgo-junit_manifests_%d.xml", config.GinkgoConfig.ParallelNode))
	RunSpecsWithDefaultAndCustomReporters(t,
		"Manifests Suite",
		[]Reporter{printer.NewlineReporter{}, junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter)))
	// SetUpK8s()
	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	// err := testEnv.Stop()
	// Expect(err).ToNot(HaveOccurred())
})
