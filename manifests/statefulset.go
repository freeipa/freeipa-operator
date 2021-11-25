package manifests

import (
	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/openlyinc/pointy"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetEphimeralVolumeForMainStatefulset Return the Volume definition when using ephimeral
// storage.
func GetEphimeralVolumeForMainStatefulset(m *v1alpha1.IDM) corev1.Volume {
	return corev1.Volume{
		Name: GetMainPersistentVolumeClaimName(m),
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

var ExecStopSystemd []string = []string{"kill", "-s", "SIGRTMIN+3", "1"}

// GetVolumeListForMainStatefulset Return the VolumeList for the Pod Spec embeded into
// the Statefulset definition, giveng an IDM object.
func GetVolumeListForMainStatefulset(m *v1alpha1.IDM) []corev1.Volume {
	var result []corev1.Volume = []corev1.Volume{}
	if m.Spec.VolumeClaimTemplate == nil {
		result = append(result, GetEphimeralVolumeForMainStatefulset(m))
	}
	result = append(result, []corev1.Volume{
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
	}...)
	return result
}

// MainStatefulsetForIDM return a master pod for an IDM CRD
func MainStatefulsetForIDM(m *v1alpha1.IDM, baseDomain string, workload string, defaultStorage string) *appsv1.StatefulSet {

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetMainStatefulsetName(m),
			Namespace: m.Namespace,
			Labels: map[string]string{
				"app":  "idm",
				"role": "main",
				"idm":  m.Name,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  "idm",
					"role": "main",
					"idm":  m.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      GetMainPodName(m),
					Namespace: m.Namespace,
					Labels: map[string]string{
						"app":  "idm",
						"role": "main",
						"idm":  m.Name,
					},
					Annotations: map[string]string{
						"openshift.io/scc": "idm",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: GetServiceAccountName(m),
					Containers: []corev1.Container{
						{
							Name:      "main",
							Image:     workload,
							TTY:       true,
							Resources: m.Spec.Resources,
							SecurityContext: &corev1.SecurityContext{
								Privileged: pointy.Bool(false),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{
										"NET_RAW",
										"SYS_CHROOT",
										"AUDIT_CONTROL",
										"AUDIT_READ",
										"BLOCK_SUSPEND",
										"DAC_READ_SEARCH",
										"IPC_LOCK",
										"IPC_OWNER",
										"LEASE",
										"LINUX_IMMUTABLE",
										"MAC_ADMIN",
										"MAC_OVERRIDE",
										"NET_ADMIN",
										"NET_BROADCAST",
										"SYS_BOOT",
										"SYS_MODULE",
										"SYS_NICE",
										"SYS_PACCT",
										"SYS_PTRACE",
										"SYS_RAWIO",
										"SYS_TIME",
										"SYS_TTY_CONFIG",
										"SYSLOG",
										"WAKE_ALARM",
										"SYS_RAWIO",
										"MKNOD",
									},
									Add: []corev1.Capability{
										"CHOWN",
										"FOWNER",
										"DAC_OVERRIDE",
										"SETUID",
										"SETGID",
										"KILL",
										"NET_BIND_SERVICE",
										"SETPCAP",
										"SETFCAP",
										"SYS_ADMIN",
										"SYS_RESOURCE",
										// Without FSETID the ccache files created has the permissions below:
										//   -rw-rw----. 1 apache apache 3464 Nov 23 21:00 host~freeipa.apps.crc.testing@APPS.CRC.TESTING-z0WzhM
										// With FSETID the ccache files created has the permissions below:
										//   -rw-------. 1 ipaapi ipaapi 5068 Nov 23 21:12 host~freeipa.apps.crc.testing@APPS.CRC.TESTING-YQgWH4
										// In both cases the /run/ipa/ccaches directory has the permissions: drwsrws---+
										//
										// The running process has effective uid:gid=ipaapi:ipaapi so FSETID makes the files
										// to inherit the owner and group from the directory because the set-user-id and set-group-id
										// are set for the directory that contain them.
										//
										// Without FSETID capability the ccache files created are not accessible by the wsgi process
										// evoking the kerberos authentication to fail. It evoked sssd client installation step to
										// fail too as in the step is invoked a RPC method on the wsgi service that has to be
										// authenticated.
										"FSETID",
									},
								},
							},
							Command: []string{"/usr/local/sbin/init"},
							Args: []string{
								"no-exit",
								"ipa-server-install",
								"-U",
								"--realm",
								GetRealm(m, baseDomain),
								"--ca-subject=" + GetCaSubject(m, baseDomain),
								"--no-ntp",
								"--no-sshd",
								"--no-ssh",
								// TODO Add --verbose if some indicator for debugging
								//      is set up such as '--verbose-freeipa' to avoid
								//      enable it isolated from '-debug' flag which is
								//      passed to the controller.
							},
							EnvFrom: buildEnvFrom(m),
							Env: []corev1.EnvVar{
								{
									Name:  "NAMESPACE",
									Value: m.Namespace,
								},
								{
									Name:  "IPA_SERVER_HOSTNAME",
									Value: GetIpaServerHostname(m, baseDomain),
								},
								{
									Name:  "container_uuid",
									Value: GenerateSystemdUUIDString(),
								},
								{
									Name:  "SYSTEMD_OFFLINE",
									Value: "1",
								},
								{
									Name:  "SYSTEMD_NSPAWN_API_VFS_WRITABLE",
									Value: "network",
								},
							},
							// https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/
							// https://www.freedesktop.org/software/systemd/man/systemd.html#SIGRTMIN+3
							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: ExecStopSystemd,
									},
								},
							},
							Ports: []corev1.ContainerPort{
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
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									// Name:      "data",
									Name:      GetMainPersistentVolumeClaimName(m),
									MountPath: "/data",
								},
								{
									Name:      "systemd-tmp",
									MountPath: "/tmp",
								},
								{
									Name:      "systemd-var-run",
									MountPath: "/var/run",
								},
								{
									Name:      "systemd-var-dirsrv",
									MountPath: "/var/run/dirsrv",
								},
							},
						},
					},
					Volumes: GetVolumeListForMainStatefulset(m),
				},
			},
			VolumeClaimTemplates: MainPersistentVolumeClaimTemplatesForIDM(m),
		},
	}

	// ctrl.SetControllerReference(m, pod, r.Scheme)
	return statefulset
}
