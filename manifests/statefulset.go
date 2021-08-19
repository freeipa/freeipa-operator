package manifests

import (
	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/openlyinc/pointy"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MainStatefulsetForIDM return a master pod for an IDM CRD
func MainStatefulsetForIDM(m *v1alpha1.IDM, baseDomain string, defaultStorage string) *appsv1.StatefulSet {
	sDirectoryOrCreate := corev1.HostPathDirectoryOrCreate
	sDirectory := corev1.HostPathDirectory

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetMainStatefulsetName(m),
			Namespace: m.Namespace,
			Labels: map[string]string{
				"app":  "idm",
				"role": "main",
				"idm":  m.Name,
			},
			// Annotations: map[string]string{
			// 	"openshift.io/scc": "idm",
			// },
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
							Name:  "main",
							Image: "quay.io/freeipa/freeipa-openshift-container:freeipa-server",
							TTY:   true,
							SecurityContext: &corev1.SecurityContext{
								Privileged: pointy.Bool(false),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{
										"NET_RAW",
										"SYS_CHROOT",
										"FSETID",
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
									},
								},
							},
							Command: []string{"/usr/sbin/init"},
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
								"--verbose",
							},
							EnvFrom: buildEnvFrom(m),
							Env: []corev1.EnvVar{
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
							},
						},
					},
					Volumes: []corev1.Volume{
						// GetDataVolumeForMainStatefulset(m, defaultStorage),
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
					},
				},
			},
			VolumeClaimTemplates: MainPersistentVolumeClaimTemplatesForIDM(m),
		},
	}

	// ctrl.SetControllerReference(m, pod, r.Scheme)
	return statefulset
}
