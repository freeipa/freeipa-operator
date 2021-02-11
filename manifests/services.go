package helper

import (
	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ServiceKerberosForIDM Create the Service manifest for the kerberos service
func ServiceKerberosForIDM(m *v1alpha1.IDM) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-kerberos",
			Namespace: m.Namespace,
			Labels:    LabelsForIDM(m),
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "freeipa",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "kerberos-tcp",
					Port: 88,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 88,
					},
					Protocol: "TCP",
				},
				{
					Name: "kerberos-udp",
					Port: 88,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 88,
					},
					Protocol: "UDP",
				},
				{
					Name: "kerberos-adm-tcp",
					Port: 749,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 749,
					},
					Protocol: "TCP",
				},
				{
					Name: "kerberos-adm-udp",
					Port: 749,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 749,
					},
					Protocol: "UDP",
				},
			},
		},
	}
	return service
}

// ServiceWebForIDM Create the Service manifest for the web interface
func ServiceWebForIDM(m *v1alpha1.IDM) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetWebServiceName(m),
			Namespace: m.Namespace,
			Labels:    LabelsForIDM(m),
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "idm",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "http-tcp",
					Port: 80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "http-tcp",
					},
					Protocol: "TCP",
				},
				{
					Name: "https-tcp",
					Port: 443,
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "https-tcp",
					},
					Protocol: "TCP",
				},
			},
		},
	}
	return service
}
