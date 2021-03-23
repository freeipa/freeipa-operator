package helper

import (
	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/openlyinc/pointy"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ServiceAccountForIDM(m *v1alpha1.IDM) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Namespace,
			Name:      GetServiceAccountName(m),
		},
		AutomountServiceAccountToken: pointy.Bool(true),
	}
	return sa
}
