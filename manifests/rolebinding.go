package helper

import (
	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	// "github.com/openlyinc/pointy"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RoleBindingForIDM return a Role for an IDM CRD
func RoleBindingForIDM(m *v1alpha1.IDM) *rbacv1.RoleBinding {

	rolebinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Namespace,
			Name:      GetRoleBindingName(m),
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     GetRoleName(m),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      GetServiceAccountName(m),
				Namespace: m.Namespace,
			},
		},
	}

	return rolebinding
}
