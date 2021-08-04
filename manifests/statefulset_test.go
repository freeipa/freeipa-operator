package manifests_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openlyinc/pointy"

	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	manifests "github.com/freeipa/freeipa-operator/manifests"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func assertStringListsEqual(data1 []string, data2 []string) {
	eq := reflect.DeepEqual(data1, data2)
	Expect(eq).Should(BeTrue())
}

func assertEnvVarEqual(data1 *corev1.EnvVar, data2 *corev1.EnvVar) {
	if data1 == nil && data2 == nil {
		return
	}
	Expect(data1 != nil && data2 != nil).Should(BeTrue())
	eq := reflect.DeepEqual(data1, data2)
	Expect(eq).Should(BeTrue())
}

func assertContainerPortEqual(data1 *corev1.ContainerPort, data2 *corev1.ContainerPort) {
	if data1 == nil && data2 == nil {
		return
	}
	Expect(data1 != nil && data2 != nil).Should(BeTrue())
	eq := reflect.DeepEqual(data1, data2)
	Expect(eq).Should(BeTrue())
}

func assertVolumeMountEqual(data1 *corev1.VolumeMount, data2 *corev1.VolumeMount) {
	if data1 == nil && data2 == nil {
		return
	}
	Expect(data1 != nil && data2 != nil).Should(BeTrue())
	eq := reflect.DeepEqual(data1, data2)
	Expect(eq).Should(BeTrue())
}

func assertVolumeEqual(data1 *corev1.Volume, data2 *corev1.Volume) {
	if data1 == nil && data2 == nil {
		return
	}
	Expect(data1 != nil && data2 != nil).Should(BeTrue())
	eq := reflect.DeepEqual(data1, data2)
	Expect(eq).Should(BeTrue())
}

func assertStringStringMapsEqual(data1 map[string]string, data2 map[string]string) {
	// https://stackoverflow.com/questions/18208394/how-to-test-the-equivalence-of-maps-in-golang
	eq := reflect.DeepEqual(data1, data2)
	Expect(eq).Should(BeTrue())
}

func assertPodSecurityContext(data *corev1.SecurityContext) {
	Expect(data).ShouldNot(BeNil())
	Expect(data.Privileged).ShouldNot(BeNil())
	Expect(*data.Privileged).Should(BeFalse())
}

var _ = Describe("LOCAL:Statefulset tests", func() {
	Describe("GetDataVolumeForMainStatefulset", func() {

		AssertDataVolumeSourcesAreNilBut := func(data *corev1.VolumeSource, but interface{}) {
			Expect(data).ShouldNot(BeNil())
			Expect(but).ShouldNot(BeNil())
			if data.AWSElasticBlockStore != but {
				Expect(data.AWSElasticBlockStore).Should(BeNil())
			}
			if data.AzureDisk != but {
				Expect(data.AzureDisk).Should(BeNil())
			}
			if data.AzureFile != but {
				Expect(data.AzureFile).Should(BeNil())
			}
			if data.CSI != but {
				Expect(data.CSI).Should(BeNil())
			}
			if data.CephFS != but {
				Expect(data.CephFS).Should(BeNil())
			}
			if data.Cinder != but {
				Expect(data.Cinder).Should(BeNil())
			}
			if data.ConfigMap != but {
				Expect(data.ConfigMap).Should(BeNil())
			}
			if data.DownwardAPI != but {
				Expect(data.DownwardAPI).Should(BeNil())
			}
			if data.EmptyDir != but {
				Expect(data.EmptyDir).Should(BeNil())
			}
			if data.Ephemeral != but {
				Expect(data.Ephemeral).Should(BeNil())
			}
			if data.FC != but {
				Expect(data.FC).Should(BeNil())
			}
			if data.FlexVolume != but {
				Expect(data.FlexVolume).Should(BeNil())
			}
			if data.Flocker != but {
				Expect(data.Flocker).Should(BeNil())
			}
			if data.GCEPersistentDisk != but {
				Expect(data.GCEPersistentDisk).Should(BeNil())
			}
			if data.GitRepo != but {
				Expect(data.GitRepo).Should(BeNil())
			}
			if data.Glusterfs != but {
				Expect(data.Glusterfs).Should(BeNil())
			}
			if data.HostPath != but {
				Expect(data.HostPath).Should(BeNil())
			}
			if data.ISCSI != but {
				Expect(data.ISCSI).Should(BeNil())
			}
			if data.NFS != but {
				Expect(data.NFS).Should(BeNil())
			}
			if data.PersistentVolumeClaim != but {
				Expect(data.PersistentVolumeClaim).Should(BeNil())
			}
			if data.PhotonPersistentDisk != but {
				Expect(data.PhotonPersistentDisk).Should(BeNil())
			}
			if data.PortworxVolume != but {
				Expect(data.PortworxVolume).Should(BeNil())
			}
			if data.Projected != but {
				Expect(data.Projected).Should(BeNil())
			}
			if data.Quobyte != but {
				Expect(data.Quobyte).Should(BeNil())
			}
			if data.RBD != but {
				Expect(data.RBD).Should(BeNil())
			}
			if data.ScaleIO != but {
				Expect(data.ScaleIO).Should(BeNil())
			}
			if data.Secret != but {
				Expect(data.Secret).Should(BeNil())
			}
			if data.StorageOS != but {
				Expect(data.StorageOS).Should(BeNil())
			}
			if data.VsphereVolume != but {
				Expect(data.VsphereVolume).Should(BeNil())
			}
		}

		AssertDataVolumeWithPVC := func(data *corev1.Volume, withClaimName string) {
			Expect(data).ShouldNot(BeNil())
			Expect(data.Name).Should(Equal("data"))
			AssertDataVolumeSourcesAreNilBut(&data.VolumeSource, data.VolumeSource.PersistentVolumeClaim)
			Expect(data.VolumeSource.PersistentVolumeClaim.ClaimName).Should(Equal(withClaimName))
		}

		AssertDataVolumeEphimeral := func(data *corev1.Volume) {
			Expect(data).ShouldNot(BeNil())
			Expect(data.Name).Should(Equal("data"))
			AssertDataVolumeSourcesAreNilBut(&data.VolumeSource, data.VolumeSource.EmptyDir)
			Expect(data.VolumeSource.EmptyDir.Medium).Should(Equal(corev1.StorageMediumDefault))
		}

		// GIVEN
		Context("an IDM with PVC information", func() {
			var idm = v1alpha1.IDM{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"app": "freeipa",
					},
				},
				Spec: v1alpha1.IDMSpec{
					Realm:          "TEST.COM",
					PasswordSecret: nil,
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadOnlyMany,
							corev1.ReadWriteMany,
							corev1.ReadWriteOnce,
						},
						VolumeName: "test-volume-name",
						VolumeMode: (*corev1.PersistentVolumeMode)(pointy.String((string)(corev1.PersistentVolumeFilesystem))),
					},
				},
			}
			// WHEN
			When("GetDataVolumeForMainStatefulset is called", func() {
				var result = manifests.GetDataVolumeForMainStatefulset(&idm, "")
				// EXPECT
				AssertDataVolumeWithPVC(&result, manifests.GetMainPersistentVolumeClaimName(&idm))
			})
		})

		// GIVEN
		Context("an IDM with no PVC information", func() {
			var idm = v1alpha1.IDM{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"app": "freeipa",
					},
				},
				Spec: v1alpha1.IDMSpec{
					Realm:               "TEST.COM",
					PasswordSecret:      nil,
					VolumeClaimTemplate: nil,
				},
			}
			// WHEN
			When("GetDataVolumeForMainStatefulset is called", func() {
				var result = manifests.GetDataVolumeForMainStatefulset(&idm, "")
				// EXPECT
				AssertDataVolumeEphimeral(&result)
			})
		})

	})

	// -----------------------------------

	Describe("MainStatefulsetForIDM", func() {
		// GIVEN
		Context("an IDM with PVC information", func() {
			var idm = v1alpha1.IDM{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"app": "freeipa",
					},
				},
				Spec: v1alpha1.IDMSpec{
					Realm:          "TEST.COM",
					PasswordSecret: pointy.String("test"),
					VolumeClaimTemplate: &corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadOnlyMany,
							corev1.ReadWriteMany,
							corev1.ReadWriteOnce,
						},
						VolumeName: "test-volume-name",
						VolumeMode: (*corev1.PersistentVolumeMode)(pointy.String((string)(corev1.PersistentVolumeFilesystem))),
					},
				},
			}
			// WHEN
			When("MainStatefulsetForIDM is called", func() {
				var result = manifests.MainStatefulsetForIDM(&idm, "test.com", "")
				// EXPECT
				Expect(result).ShouldNot(BeNil())
				Expect(result.ObjectMeta.Name).Should(Equal(idm.Name + "-main"))
				Expect(result.ObjectMeta.Namespace).Should(Equal(idm.Namespace))
				It("has the labels expected", func(done Done) {
					go func() {
						defer GinkgoRecover()

						assertStringStringMapsEqual(result.ObjectMeta.Labels, map[string]string{
							"app":  "idm",
							"role": "main",
							"idm":  idm.Name,
						})

						close(done)
					}()
				})
				// assertMapsEqual(result.ObjectMeta.Annotations, map[string]string{
				// 	"openshift.io/scc": "idm",
				// })
				Expect(result.Spec.Template.ObjectMeta.Name).Should(Equal(idm.Name + "-main"))
				Expect(result.Spec.Template.ObjectMeta.Namespace).Should(Equal(idm.Namespace))
				It("has the pod template labels expected", func(done Done) {
					go func() {
						defer GinkgoRecover()
						assertStringStringMapsEqual(result.Spec.Template.ObjectMeta.Labels, map[string]string{
							"app":  "idm",
							"role": "main",
							"idm":  idm.Name,
						})
						close(done)
					}()
				})
				It("has the pod annotation for scc", func(done Done) {
					go func() {
						defer GinkgoRecover()
						assertStringStringMapsEqual(result.Spec.Template.ObjectMeta.Annotations, map[string]string{
							"openshift.io/scc": "idm",
						})
						close(done)
					}()
				})
				Expect(result.Spec.Template.Spec.ServiceAccountName).Should(Equal("idm"))
				Expect(len(result.Spec.Template.Spec.Containers)).Should(Equal(1))
				Expect(result.Spec.Template.Spec.Containers[0].Name).Should(Equal("main"))
				Expect(result.Spec.Template.Spec.Containers[0].Image).Should(Equal("quay.io/freeipa/freeipa-openshift-container:freeipa-server"))
				Expect(result.Spec.Template.Spec.Containers[0].TTY).Should(BeTrue())
				assertPodSecurityContext(result.Spec.Template.Spec.Containers[0].SecurityContext)
				assertStringListsEqual(result.Spec.Template.Spec.Containers[0].Command, []string{
					"/usr/sbin/init",
				})
				assertStringListsEqual(result.Spec.Template.Spec.Containers[0].Args[:4], []string{
					"no-exit",
					"ipa-server-install",
					"-U",
					"--realm",
					// manifests.GetRealm(&idm, "test.com"),
					// "--ca-subject=" + manifests.GetCaSubject(&idm, "test.com"),
				})
				assertStringListsEqual(result.Spec.Template.Spec.Containers[0].Args[6:], []string{
					// manifests.GetRealm(&idm, "test.com"),
					// "--ca-subject=" + manifests.GetCaSubject(&idm, "test.com"),
					"--no-ntp",
					"--no-sshd",
					"--no-ssh",
					"--verbose",
				})
				envList := []corev1.EnvVar{
					{
						Name:  "KRB5_TRACE",
						Value: "/dev/console",
					},
					{
						Name:  "SYSTEMD_LOG_LEVEL",
						Value: "debug",
					},
					{
						Name:  "SYSTEMD_LOG_COLOR",
						Value: "no",
					},
					{
						Name:  "INIT_WRAPPER",
						Value: "1",
					},
					{
						Name:  "DEBUG_TRACE",
						Value: "2",
					},
					{
						Name:  "NAMESPACE",
						Value: idm.Namespace,
					},
					{
						Name:  "IPA_SERVER_HOSTNAME",
						Value: manifests.GetIpaServerHostname(&idm, "test.com"),
					},
					{
						Name:  "container_uuid",
						Value: manifests.GenerateSystemdUUIDString(),
					},
					{
						Name:  "SYSTEMD_OFFLINE",
						Value: "1",
					},
					{
						Name:  "SYSTEMD_NSPAWN_API_VFS_WRITABLE",
						Value: "network",
					},
				}
				Expect(len(result.Spec.Template.Spec.Containers[0].Env)).Should(Equal(len(envList)))
				for index, item := range result.Spec.Template.Spec.Containers[0].Env[:7] {
					By("Checking result.Spec.Template.Spec.Containers[0].Env[].Name: " + item.Name)
					assertEnvVarEqual(&item, &envList[index])
				}
				for index, item := range result.Spec.Template.Spec.Containers[0].Env[8:] {
					By("Checking result.Spec.Template.Spec.Containers[0].Env[].Name: " + item.Name)
					assertEnvVarEqual(&item, &envList[index+8])
				}
				portList := []corev1.ContainerPort{
					{
						Name:          "http-tcp",
						Protocol:      "TCP",
						ContainerPort: 80,
					},
					{
						Name:          "https-tcp",
						Protocol:      "TCP",
						ContainerPort: 443,
					},
				}
				Expect(len(result.Spec.Template.Spec.Containers[0].Ports)).Should(Equal(len(portList)))
				for index, item := range result.Spec.Template.Spec.Containers[0].Ports {
					By("Checking result.Spec.Template.Spec.Containers[0].Ports[].Name: " + item.Name)
					assertContainerPortEqual(&item, &portList[index])
				}

				It("has the volumeMountList expected", func(done Done) {
					go func() {
						defer GinkgoRecover()
						volumeMountList := []corev1.VolumeMount{
							{
								Name:      "data",
								MountPath: "/data",
							},
							{
								Name:      "systemd-tmp",
								MountPath: "/tmp",
							},
							{
								Name:      "systemd-sys",
								MountPath: "/sys",
								ReadOnly:  true,
							},
							{
								Name:      "systemd-sys-fs-selinux",
								MountPath: "/sys/fs/selinux",
								ReadOnly:  true,
							},
							{
								Name:      "systemd-sys-firmware",
								MountPath: "/sys/firmware",
								ReadOnly:  true,
							},
							{
								Name:      "systemd-sys-kernel",
								MountPath: "/sys/kernel",
								ReadOnly:  true,
							},
							{
								Name:      "systemd-var-run",
								MountPath: "/var/run",
							},
							{
								Name:      "systemd-var-dirsrv",
								MountPath: "/var/run/dirsrv",
							},
						}
						By("Checking VolumeMounts length")
						Expect(len(result.Spec.Template.Spec.Containers[0].VolumeMounts)).Should(Equal(len(volumeMountList)))
						for index, item := range result.Spec.Template.Spec.Containers[0].VolumeMounts {
							By("Checking result.Spec.Template.Spec.Containers[0].VolumeMounts[].Name: " + item.Name)
							assertVolumeMountEqual(&item, &volumeMountList[index])
						}
						close(done)
					}()
				})

				It("has the volumeList expected", func(done Done) {
					go func() {
						defer GinkgoRecover()
						sDirectoryOrCreate := corev1.HostPathDirectoryOrCreate
						sDirectory := corev1.HostPathDirectory

						volumeList := []corev1.Volume{
							manifests.GetDataVolumeForMainPod(&idm, "ephimeral"),
							{
								Name: "systemd-sys",
								VolumeSource: corev1.VolumeSource{
									HostPath: &corev1.HostPathVolumeSource{
										Path: "/sys",
										Type: &sDirectoryOrCreate,
									},
								},
							},
							{
								Name: "systemd-sys-fs-selinux",
								VolumeSource: corev1.VolumeSource{
									HostPath: &corev1.HostPathVolumeSource{
										Path: "/sys/fs/selinux",
										Type: &sDirectory,
									},
								},
							},
							{
								Name: "systemd-sys-firmware",
								VolumeSource: corev1.VolumeSource{
									HostPath: &corev1.HostPathVolumeSource{
										Path: "/sys/firmware",
										Type: &sDirectory,
									},
								},
							},
							{
								Name: "systemd-sys-kernel",
								VolumeSource: corev1.VolumeSource{
									HostPath: &corev1.HostPathVolumeSource{
										Path: "/sys/kernel",
										Type: &sDirectory,
									},
								},
							},
							{
								Name: "systemd-var-run",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{
										Medium: corev1.StorageMedium("Memory"),
									},
								},
							},
							{
								Name: "systemd-var-dirsrv",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{
										Medium: corev1.StorageMedium("Memory"),
									},
								},
							},
							{
								Name: "systemd-run-rpcbind",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{
										Medium: corev1.StorageMedium("Memory"),
									},
								},
							},
							{
								Name: "systemd-tmp",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{
										Medium: corev1.StorageMedium("Memory"),
									},
								},
							},
						}
						By("Checking Volumes length")
						Expect(len(result.Spec.Template.Spec.Volumes)).Should(Equal(len(volumeList)))
						for index, item := range result.Spec.Template.Spec.Volumes {
							By("Checking result.Spec.Template.Spec.Volumes[].Name: " + item.Name)
							assertVolumeEqual(&item, &volumeList[index])
						}
						close(done)
					}()
				})
			})
		})

	})

})
