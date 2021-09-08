// https://itnext.io/testing-kubernetes-operators-with-ginkgo-gomega-and-the-operator-runtime-6ad4c2492379

// https://semaphoreci.com/community/tutorials/getting-started-with-bdd-in-go-using-ginkgo

package manifests_test

import (
	b64 "encoding/base64"

	. "github.com/onsi/ginkgo"
	"github.com/openlyinc/pointy"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	manifests "github.com/freeipa/freeipa-operator/manifests"
	. "github.com/onsi/gomega"
)

const (
	// Helpful tool for the regular expression here:  https://regex101.com/
	regexpForGeneratedPassword = "((?:[[:digit:]]|[[:alpha:]]){5})-((?:[[:digit:]]|[[:alpha:]]){5})-((?:[[:digit:]]|[[:alpha:]]){5})-((?:[[:digit:]]|[[:alpha:]]){5})"
)

var _ = Describe("UNIT:GenerateRandomPassword", func() {

	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Context("Password generation for the secrets", func() {
		When("a password is generated", func() {
			password := manifests.GenerateRandomPassword()
			It("has a len of 23", func() {
				Expect(len(password)).Should(Equal(23))
			})
			It("contains 4 alphanumeric blocks, with len of 5 and splited by '-'", func() {
				Expect(password).Should(
					MatchRegexp(regexpForGeneratedPassword),
				)
			})
		})
	})
})

var _ = Describe("UNIT:GetSecretName", func() {

	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Context("Spec has PasswordSecret", func() {
		idm := v1alpha1.IDM{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "sample",
				Name:      "idm-sample",
			},
			Spec: v1alpha1.IDMSpec{
				PasswordSecret: pointy.String("sample-secret"),
			},
		}
		When("the secret name is retrieved for an IDM resource", func() {
			secretName := manifests.GetSecretName(&idm)
			It("matches the 'PasswordSecret' field", func() {
				Expect(secretName).Should(Equal(*idm.Spec.PasswordSecret))
			})
		})
	})

	Context("Spec has a nil PasswordSecret", func() {
		idm := v1alpha1.IDM{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "sample",
				Name:      "idm-sample",
			},
			Spec: v1alpha1.IDMSpec{
				PasswordSecret: nil,
			},
		}
		When("the secret name is retrieved for an IDM resource", func() {
			secretName := manifests.GetSecretName(&idm)
			It("matches the ObjectMeta.Name field", func() {
				Expect(secretName).Should(Equal(idm.ObjectMeta.Name))
			})
		})
	})
})

var _ = Describe("UNIT:SecretForIDM", func() {
	Context("Given an IDM resource", func() {
		var cr = v1alpha1.IDM{
			ObjectMeta: v1.ObjectMeta{
				Name:      "idm-sample",
				Namespace: "test",
			},
			Spec: v1alpha1.IDMSpec{
				PasswordSecret: pointy.String("secret-sample"),
			},
		}

		When("we generate a manifest by SecretForIDM and using a generated password", func() {
			var adminPassword = manifests.GenerateRandomPassword()
			var dsPassword = manifests.GenerateRandomPassword()
			var manifest = manifests.SecretForIDM(&cr, adminPassword, dsPassword)
			It("is not nil", func() {
				Expect(manifest).ShouldNot(BeNil())
			})
			It("is immutable", func() {
				Expect(manifest.Immutable).ShouldNot(BeNil())
				Expect(*manifest.Immutable).Should(BeTrue())
			})
			It("has the name provided at 'PasswordSecret'", func() {
				Expect(manifest.ObjectMeta.Name).Should(Equal(*cr.Spec.PasswordSecret))
			})
			It("belongs to the same namespace as the CR", func() {
				Expect(manifest.ObjectMeta.Namespace).Should(Equal(cr.ObjectMeta.Namespace))
			})
			It("has type 'SecretTypeOpaque'", func() {
				Expect(manifest.Type).Should(Equal(corev1.SecretTypeOpaque))
			})
			It("has a ADMIN_PASSWORD entry", func() {
				_, ok := manifest.StringData["ADMIN_PASSWORD"]
				Expect(ok).Should(BeTrue())
			})
			It("has a DS_PASSWORD entry", func() {
				_, ok := manifest.StringData["DS_PASSWORD"]
				Expect(ok).Should(BeTrue())
			})
			It("has the ADMIN_PASSWORD entry coded in base64", func() {
				val := manifest.StringData["ADMIN_PASSWORD"]
				decodedVal, err := b64.StdEncoding.DecodeString(string(val))
				decodedValString := string(decodedVal)
				Expect(err).Should(BeNil())
				Expect(adminPassword).Should(Equal(decodedValString))
			})
			It("has the DS_PASSWORD entry coded in base64", func() {
				val := manifest.StringData["DS_PASSWORD"]
				decodedVal, err := b64.StdEncoding.DecodeString(string(val))
				decodedValString := string(decodedVal)
				Expect(err).Should(BeNil())
				Expect(dsPassword).Should(Equal(decodedValString))
			})
		})

		When("we generate a manifest by SecretForIDM and providing an empty password", func() {
			var manifest = manifests.SecretForIDM(&cr, "", "")
			It("is not nil", func() {
				Expect(manifest).ShouldNot(BeNil())
			})
			It("is immutable", func() {
				Expect(manifest.Immutable).ShouldNot(BeNil())
				Expect(*manifest.Immutable).Should(BeTrue())
			})
			It("has the name provided at 'PasswordSecret'", func() {
				Expect(manifest.ObjectMeta.Name).Should(Equal(*cr.Spec.PasswordSecret))
			})
			It("belongs to the same namespace as the CR", func() {
				Expect(manifest.ObjectMeta.Namespace).Should(Equal(cr.ObjectMeta.Namespace))
			})
			It("has type 'SecretTypeOpaque'", func() {
				Expect(manifest.Type).Should(Equal(corev1.SecretTypeOpaque))
			})
			It("has a ADMIN_PASSWORD entry", func() {
				_, ok := manifest.StringData["ADMIN_PASSWORD"]
				Expect(ok).Should(BeTrue())
			})
			It("has a DS_PASSWORD entry", func() {
				_, ok := manifest.StringData["DS_PASSWORD"]
				Expect(ok).Should(BeTrue())
			})
			It("has a password entry coded in base64 that match the pattern 'xxxxx-xxxxx-xxxxx-xxxxx', where x are aphanumerics characters.", func() {
				val := manifest.StringData["ADMIN_PASSWORD"]
				decodedVal, err := b64.StdEncoding.DecodeString(string(val))
				decodedValString := string(decodedVal)
				Expect(err).Should(BeNil())
				Expect(decodedValString).Should(MatchRegexp(regexpForGeneratedPassword))
			})
			It("has a password entry coded in base64 that match the pattern 'xxxxx-xxxxx-xxxxx-xxxxx', where x are aphanumerics characters.", func() {
				val := manifest.StringData["DS_PASSWORD"]
				decodedVal, err := b64.StdEncoding.DecodeString(string(val))
				decodedValString := string(decodedVal)
				Expect(err).Should(BeNil())
				Expect(decodedValString).Should(MatchRegexp(regexpForGeneratedPassword))
			})
		})
	})
})
