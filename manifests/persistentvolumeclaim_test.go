package manifests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openlyinc/pointy"

	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	manifests "github.com/freeipa/freeipa-operator/manifests"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("LOCAL:PersistentVolumeClaim tests", func() {
	Describe("IsResourceRequirementsEmpty", func() {
		// GIVEN
		Context("IsResourceRequirementsEmpty positive case", func() {
			var resource = corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("10Gi"),
				},
			}
			// WHEN
			When("IsReourceRequirementsEmpty with non empty argument", func() {
				var result = manifests.IsResourceRequirementsEmpty(&resource)
				// EXPECT
				Expect(result).Should(BeFalse())
			})
		})

		// GIVEN
		Context("IsResourceRequirementsEmpty negative cases", func() {
			var resource = corev1.ResourceRequirements{}
			// WHEN
			When("IsReourceRequirementsEmpty with empty item", func() {
				var result = manifests.IsResourceRequirementsEmpty(&resource)
				// EXPECT
				Expect(result).Should(BeTrue())
			})
		})

		// GIVEN
		Context("IsResourceRequirementsEmpty called with nil", func() {
			// WHEN
			When("IsReourceRequirementsEmpty with empty item", func() {
				var result = manifests.IsResourceRequirementsEmpty(nil)
				// EXPECT
				Expect(result).Should(BeTrue())
			})
		})
	})

	// -----------------------------------

	Describe("IsPersistentVolumeClaimSpecEmpty", func() {
		Context("IsPersistentVolumeClaimSpecEmpty called with empty values", func() {
			var pvc *corev1.PersistentVolumeClaimSpec = nil
			When("IsPersistentVolumeClaimSpecEmpty called with nil", func() {
				var r = manifests.IsPersistentVolumeClaimSpecEmpty(pvc)
				Expect(r).Should(BeTrue())
			})
			When("IsPersistentVolumeClaimSpecEmpty is called with empty values", func() {
				var r = manifests.IsPersistentVolumeClaimSpecEmpty(&corev1.PersistentVolumeClaimSpec{})
				Expect(r).Should(BeTrue())
			})
		})

		Context("IsPersistentVolumeClaimSpecEmpty is called with no empty values", func() {
			// GIVEN
			var pvcs []corev1.PersistentVolumeClaimSpec = []corev1.PersistentVolumeClaimSpec{
				{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
				},
				{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "freeipa",
						},
					},
				},
				{
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("10Gi"),
						},
					},
				},
				{
					VolumeName: "freeipa",
				},
				{
					StorageClassName: pointy.String("standard"),
				},
				{
					VolumeMode: (*corev1.PersistentVolumeMode)(pointy.String(string(corev1.PersistentVolumeFilesystem))),
				},
			}

			for _, pvc := range pvcs {
				// WHEN
				var r = manifests.IsPersistentVolumeClaimSpecEmpty(&pvc)
				// EXPECT
				Expect(r).Should(BeFalse())
			}
		})
	})

	// -----------------------------------

	Describe("CheckVolumeInformation", func() {
		// Given
		var idms []*v1alpha1.IDM = []*v1alpha1.IDM{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1alpha1",
					Kind:       "IDM",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-freeipa",
				},
				Spec: v1alpha1.IDMSpec{
					Realm:          "freeipa.com",
					PasswordSecret: pointy.String("test-freeipa"),
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						VolumeMode: (*corev1.PersistentVolumeMode)(pointy.String(string(corev1.PersistentVolumeFilesystem))),
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("10Gi"),
							},
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1alpha1",
					Kind:       "IDM",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-freeipa",
				},
				Spec: v1alpha1.IDMSpec{
					Realm:               "freeipa.com",
					PasswordSecret:      pointy.String("test-freeipa"),
					VolumeClaimTemplate: nil,
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1alpha1",
					Kind:       "IDM",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-freeipa",
				},
				Spec: v1alpha1.IDMSpec{
					Realm:               "freeipa.com",
					PasswordSecret:      pointy.String("test-freeipa"),
					VolumeClaimTemplate: nil,
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1alpha1",
					Kind:       "IDM",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-freeipa",
				},
				Spec: v1alpha1.IDMSpec{
					Realm:               "freeipa.com",
					PasswordSecret:      pointy.String("test-freeipa"),
					VolumeClaimTemplate: nil,
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1alpha1",
					Kind:       "IDM",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-freeipa",
				},
				Spec: v1alpha1.IDMSpec{
					Realm:               "freeipa.com",
					PasswordSecret:      pointy.String("test-freeipa"),
					VolumeClaimTemplate: nil,
				},
			},
		}
		var ds []string = []string{
			"",
			"ephimeral",
			"hostpath",
			"",
			"wrongdefaultsorage",
		}
		var r []error = []error{
			nil,
			nil,
			nil,
			errors.New("no VolumeClaimTemplate nor defaultStorage found"),
			errors.New("defaultStorage has an invalid value; only 'ephimeral' and 'hostpath' are allowed"),
		}

		for index, item := range r {
			// When
			err := manifests.CheckVolumeInformation(idms[index], ds[index])
			// Expect
			if item == nil {
				Expect(err).Should(BeNil())
			} else {
				Expect(err.Error()).Should(Equal(item.Error()))
			}
		}
	})

	// -----------------------------------

	Describe("MainPersistentVolumeClaimForIDM", func() {
		// GIVEN
		var idms []*v1alpha1.IDM = []*v1alpha1.IDM{
			{
				Spec: v1alpha1.IDMSpec{
					VolumeClaimTemplate: nil,
				},
			},
			{
				Spec: v1alpha1.IDMSpec{
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Selector: nil,
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("10Gi"),
							},
						},
						StorageClassName: pointy.String("standard"),
						VolumeName:       "test",
						VolumeMode:       (*corev1.PersistentVolumeMode)(pointy.String((string)(corev1.PersistentVolumeFilesystem))),
						DataSource: &corev1.TypedLocalObjectReference{
							APIGroup: &corev1.SchemeGroupVersion.Group,
							Kind:     "Snapshot",
							Name:     "test",
						},
					},
				},
			},
		}
		var results []*corev1.PersistentVolumeClaim = []*corev1.PersistentVolumeClaim{
			nil,
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      manifests.GetMainPersistentVolumeClaimName(idms[0]),
					Namespace: idms[0].Namespace,
					Labels:    manifests.LabelsForIDM(idms[0]),
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					DataSource: &corev1.TypedLocalObjectReference{
						APIGroup: &corev1.SchemeGroupVersion.Group,
						Kind:     "Snapshot",
						Name:     "test",
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceName(corev1.ResourceStorage): resource.MustParse("10Gi"),
						},
					},
					StorageClassName: pointy.String("standard"),
					VolumeMode:       (*corev1.PersistentVolumeMode)(pointy.String((string)(corev1.PersistentVolumeFilesystem))),
					VolumeName:       "test",
				},
			},
		}
		for index, r := range results {
			// WHEN
			result := manifests.MainPersistentVolumeClaimForIDM(idms[index])

			if r == nil {
				Expect(result).Should(BeNil())
				continue
			}

			// EXPECT
			Expect(result.ObjectMeta.Name).Should(Equal(r.ObjectMeta.Name))
			Expect(result.ObjectMeta.Namespace).Should(Equal(r.ObjectMeta.Namespace))
			Expect(result.ObjectMeta.Labels).Should(Equal(r.ObjectMeta.Labels))

			Expect(result.Spec.AccessModes).Should(Equal(r.Spec.AccessModes))
			if r.Spec.DataSource == nil {
				Expect(result.Spec.DataSource).Should(BeNil())
			} else {
				Expect(*result.Spec.DataSource).Should(Equal(*r.Spec.DataSource))
			}
			Expect(result.Spec.Resources).Should(Equal(r.Spec.Resources))
			if r.Spec.StorageClassName == nil {
				Expect(result.Spec.StorageClassName).Should(BeNil())
			} else {
				Expect(*result.Spec.StorageClassName).Should(Equal(*r.Spec.StorageClassName))
			}
			if r.Spec.VolumeMode == nil {
				Expect(result.Spec.VolumeMode).Should(BeNil())
			} else {
				Expect(*result.Spec.VolumeMode).Should(Equal(*r.Spec.VolumeMode))
			}
			Expect(result.Spec.VolumeName).Should(Equal(r.Spec.VolumeName))
		}
	})
})
