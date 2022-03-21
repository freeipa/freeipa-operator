package v1alpha1_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openlyinc/pointy"

	v1alpha1 "github.com/freeipa/freeipa-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var _ = Describe("LOCAL:WebHook test create", func() {
	Describe("ValidateCreate", func() {
		var cpus3 = resource.NewQuantity(3, resource.Format(""))
		var memory4Gi = resource.NewQuantity(4, resource.Format("Giga"))

		type TestTableWebHookValidateCreate struct {
			Spec  v1alpha1.IDMSpec
			Error error
		}
		var table []TestTableWebHookValidateCreate = []TestTableWebHookValidateCreate{
			{
				// Success case
				Spec: v1alpha1.IDMSpec{
					Host:           "freeipa.example.testing",
					Realm:          "EXAMPLE.TESTING",
					PasswordSecret: nil,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"cpu":    *cpus3,
							"memory": *memory4Gi,
						},
					},
				},
				Error: nil,
			},
		}
		for idx, test := range table {
			var contextString = fmt.Sprintf("ValidateCreate %d", idx)
			// GIVEN
			Context(contextString, func() {
				var record *v1alpha1.IDM = &v1alpha1.IDM{
					Spec: test.Spec,
				}
				// WHEN
				When("ValidateCreate is called", func() {
					var result = record.ValidateCreate()
					// EXPECT
					if test.Error == nil {
						Expect(result).Should(BeNil())
					} else {
						Expect(result).Should(BeEquivalentTo(test.Error))
					}
				})
			})
		}
	})

	Describe("ValidateUpdate", func() {
		var cpus3 = resource.NewQuantity(3, resource.Format(""))
		var memory4Gi = resource.NewQuantity(4, resource.Format("Giga"))
		var cpus2 = resource.NewQuantity(2, resource.Format(""))
		var memory3Gi = resource.NewQuantity(3, resource.Format("Giga"))

		type TestTableWebHookValidateUpdate struct {
			SpecOld v1alpha1.IDMSpec
			SpecNew v1alpha1.IDMSpec
			Error   error
		}
		var _oldSpec = v1alpha1.IDMSpec{
			Host:           "freeipa.example.testing",
			Realm:          "EXAMPLE.TESTING",
			PasswordSecret: pointy.String("password-secret"),
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					"cpu":    *cpus3,
					"memory": *memory4Gi,
				},
			},
			VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
				VolumeName: "test",
			},
		}
		var table []TestTableWebHookValidateUpdate = []TestTableWebHookValidateUpdate{
			{
				// Success case
				SpecOld: _oldSpec,
				SpecNew: _oldSpec,
				Error:   nil,
			},
			{
				// Change Host
				SpecOld: _oldSpec,
				SpecNew: v1alpha1.IDMSpec{
					Host:           "changed.example.testing",
					Realm:          "EXAMPLE.TESTING",
					PasswordSecret: pointy.String("password-secret"),
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"cpu":    *cpus3,
							"memory": *memory4Gi,
						},
					},
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						VolumeName: "test",
					},
				},
				Error: fmt.Errorf("IDM.Spec.Host is immutable"),
			},
			{
				// Change Realm
				SpecOld: _oldSpec,
				SpecNew: v1alpha1.IDMSpec{
					Host:           "freeipa.example.testing",
					Realm:          "OTHER.TESTING",
					PasswordSecret: pointy.String("password-secret"),
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"cpu":    *cpus3,
							"memory": *memory4Gi,
						},
					},
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						VolumeName: "test",
					},
				},
				Error: fmt.Errorf("IDM.Spec.Realm is immutable"),
			},
			{
				// PasswordSecret
				SpecOld: _oldSpec,
				SpecNew: v1alpha1.IDMSpec{
					Host:           "freeipa.example.testing",
					Realm:          "EXAMPLE.TESTING",
					PasswordSecret: pointy.String("other-password-secret"),
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"cpu":    *cpus3,
							"memory": *memory4Gi,
						},
					},
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						VolumeName: "test",
					},
				},
				Error: fmt.Errorf("IDM.Spec.PasswordSecret is immutable"),
			},
			{
				// Resources
				SpecOld: _oldSpec,
				SpecNew: v1alpha1.IDMSpec{
					Host:           "freeipa.example.testing",
					Realm:          "EXAMPLE.TESTING",
					PasswordSecret: pointy.String("password-secret"),
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"cpu":    *cpus2,
							"memory": *memory3Gi,
						},
					},
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						VolumeName: "test",
					},
				},
				Error: fmt.Errorf("IDM.Spec.Resources is immutable"),
			},
			{
				// VolumeClaimTemplate
				SpecOld: _oldSpec,
				SpecNew: v1alpha1.IDMSpec{
					Host:           "freeipa.example.testing",
					Realm:          "EXAMPLE.TESTING",
					PasswordSecret: pointy.String("password-secret"),
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"cpu":    *cpus2,
							"memory": *memory3Gi,
						},
					},
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						VolumeName: "another-test-volume",
					},
				},
				Error: fmt.Errorf("IDM.Spec.VolumeClaimTemplate is immutable"),
			},
		}
		for idx, test := range table {
			var contextString = fmt.Sprintf("ValidateUpdate: %d", idx)
			// GIVEN
			Context(contextString, func() {
				var recordOld *v1alpha1.IDM = &v1alpha1.IDM{
					Spec: test.SpecOld,
				}
				var recordNew *v1alpha1.IDM = &v1alpha1.IDM{
					Spec: test.SpecNew,
				}
				// WHEN
				var whenString = fmt.Sprintf("ValidateUpdate is called: %d", idx)
				When(whenString, func() {
					var result = recordNew.ValidateUpdate(recordOld)
					// EXPECT
					if test.Error == nil {
						Expect(result).Should(BeNil())
					} else {
						Expect(result.Error()).Should(BeEquivalentTo(test.Error.Error()))
					}
				})
			})
		}
	})
})
