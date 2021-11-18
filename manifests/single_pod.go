package manifests

import (
	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/openlyinc/pointy"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func buildEnvFrom(m *v1alpha1.IDM) []corev1.EnvFromSource {
	var result []corev1.EnvFromSource

	if m != nil && m.Spec.PasswordSecret != nil {
		result = append(result, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: *m.Spec.PasswordSecret,
				},
			},
		})
	}

	return result
}

func needsPersistentVolumeClaim(m *v1alpha1.IDM) bool {
	return m.Spec.VolumeClaimTemplate != nil
}

// GetDataVolumeForMainPod Return a corev1.Volume block for the PVC to be mounted
// Return a corev1.Volume structure properly filled.
func GetDataVolumeForMainPod(m *v1alpha1.IDM, defaultStorage string) corev1.Volume {
	if needsPersistentVolumeClaim(m) {
		return corev1.Volume{
			Name: "data",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: GetMainPersistentVolumeClaimName(m),
				},
			},
		}
	}

	// Set /data volume according to defaultStorage
	if defaultStorage == "ephemeral" {
		return corev1.Volume{
			Name: "data",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium: corev1.StorageMediumDefault,
				},
			},
		}
	}

	// By default return ephemeral
	return corev1.Volume{
		Name: "data",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium: corev1.StorageMediumDefault,
			},
		},
	}
}

// MainPodForIDM return a master pod for an IDM CRD
func MainPodForIDM(m *v1alpha1.IDM, baseDomain string, workload string, defaultStorage string) *corev1.Pod {
	pod := &corev1.Pod{
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
							Name:      "data",
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
						// This fix 'Error - Problem accessing the lockfile /var/lock/dirsrv/slapd-*/lock'
						{
							Name:      "dirsrv-var-lock-dirsrv",
							MountPath: "/var/lock/dirsrv",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				GetDataVolumeForMainPod(m, defaultStorage),
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
				{
					Name: "dirsrv-var-lock-dirsrv",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium: corev1.StorageMedium("Memory"),
						},
					},
				},
			},
		},
	}

	// ctrl.SetControllerReference(m, pod, r.Scheme)
	return pod
}
