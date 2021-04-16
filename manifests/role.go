package manifests

import (
	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	// "github.com/openlyinc/pointy"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RoleForIDM return a Role for an IDM CRD
func RoleForIDM(m *v1alpha1.IDM) *rbacv1.Role {
	// - apiGroups: ["policy"]
	// resources:
	//   - "serviceaccounts"
	//   - "roles"
	//   - "rolebindings"
	// resourceNames:
	//   - idm
	// verbs:
	//   - use

	//   - apiGroups: [""]
	//     resources:
	//       - "pods/finalizers"
	//       - "services/finalizers"
	//     verbs:
	//       - "*"
	//   - apiGroups: ["route.openshift.io"]
	//     resources:
	//       - "routes/finalizers"
	//     verbs:
	//       - "*"

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: m.Namespace,
			Name:      GetRoleName(m),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"security.openshift.io"},
				Resources: []string{
					"securitycontextconstraints",
				},
				ResourceNames: []string{
					"idm-operator-idm",
				},
				Verbs: []string{"use"},
			},
			{
				APIGroups: []string{"authorization.openshift.io"},
				Resources: []string{
					"roles",
					"rolebindings",
				},
				Verbs: []string{"use"},
			},
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{
					"roles",
					"rolebindings",
				},
				Verbs: []string{"use"},
			},
		},
	}

	return role
}
