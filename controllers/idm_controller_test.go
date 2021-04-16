// https://itnext.io/testing-kubernetes-operators-with-ginkgo-gomega-and-the-operator-runtime-6ad4c2492379

// https://semaphoreci.com/community/tutorials/getting-started-with-bdd-in-go-using-ginkgo

package controllers_test

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	. "github.com/freeipa/freeipa-operator/controllers"

	. "github.com/onsi/ginkgo"

	// . "github.com/onsi/ginkgo/extensions/table"
	manifests "github.com/freeipa/freeipa-operator/manifests"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
)

var namespace = "idm-system"

// var _ = Describe("IDMReconcillerDeployment Controller", func() {

// 	const timeout = time.Second * 60
// 	const interval = time.Second * 1

// 	ctx := context.Background()
// 	var toCreate *v1alpha1.IDM
// 	// var generatedObject *v1alpha1.IDM = &v1alpha1.IDM{}
// 	var err error

// 	Context("When loading template", func() {
// 		It("Should load properly", func() {
// 			toCreate, err = loadIDMTemplate("test_idm.yaml")
// 			Expect(err).Should(BeNil())
// 			Expect(k8sClient.Create(ctx, toCreate)).Should(Succeed())
// 		})
// 	})

// 	// It("Get the generated object", func() {
// 	// 	Eventually(func() bool {
// 	// 		err = k8sClient.Get(ctx, types.NamespacedName{Name: toCreate.Name, Namespace: namespace}, generatedObject)
// 	// 		return err == nil
// 	// 	}, timeout, interval).Should(BeTrue())
// 	// 	Expect(toCreate.Name).To(Equal("freeipa"))
// 	// })
// })

var _ = Describe("LOCAL:IDMReconciller", func() {
	var (
		sut *IDMReconciler
	)

	// ctx := context.Background()

	BeforeEach(func() {
		// SetUpK8s()
		// Expect(k8sManager).ShouldNot(BeNil())
		// sut = &IDMReconciler{}
		// err := (&IDMReconciler{
		// 	Client: k8sClient,
		// 	Log:    ctrl.Log.WithName("controllers").WithName("IDM"),
		// 	Scheme: &runtime.Scheme{},
		// }).SetupWithManager(GetK8sManager())
		// Expect(err).ToNot(HaveOccurred())
		// Expect(sut).ShouldNot(BeNil())
		sut = GetReconciler()
		Expect(sut).ShouldNot(BeNil())

		// absdir, err := filepath.Abs("../config/templates")
		// Expect(err).ToNot(HaveOccurred())
		// os.Setenv("TEMPLATES_PATH", absdir)
	})

	AfterEach(func() {
		sut = nil
	})

	Context("IDMReconciler.loadPodTemplates", func() {
		When("create Por for IDM", func() {
			idm := &v1alpha1.IDM{
				ObjectMeta: v1.ObjectMeta{
					Name:      "test-create-pod-for-idm",
					Namespace: "tests",
				},
				Spec: v1alpha1.IDMSpec{
					Realm: "FREEIPA",
				},
			}
			Expect(idm).ShouldNot(BeNil())
			pod := manifests.MainPodForIDM(idm, "localhost")
			It("return nil error and valid Pod Object", func() {
				Expect(pod).ShouldNot(BeNil())
			})
		})
	})
})

func loadIDMTemplate(filename string) (*v1alpha1.IDM, error) {
	manifest, err := ioutil.ReadFile(os.Getenv("SAMPLES_PATH") + "/" + filename)
	if err != nil {
		return nil, err
	}
	idm := &v1alpha1.IDM{}
	dec := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(manifest), 1000)
	if err := dec.Decode(&idm); err != nil {
		return nil, err
	}
	return idm, nil
}
