package manifests

import (
	b64 "encoding/base64"

	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/openlyinc/pointy"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateRandomPassword Generate a random password
// Return a random string following the pattern:
// XXXXX-XXXXX-XXXXX-XXXXX
func GenerateRandomPassword() string {
	return RandStringBytes(5) + "-" + RandStringBytes(5) + "-" + RandStringBytes(5) + "-" + RandStringBytes(5)
}

// Get the name from the ObjectMeta.
// Return the secret name.
func GetSecretName(m *v1alpha1.IDM) string {
	if m.Spec.PasswordSecret != nil {
		return *m.Spec.PasswordSecret
	}
	return m.ObjectMeta.Name
}

// SecretForIDM Create a secret for the freeipa password
// m It is the idm manifest that triggered the event.
// password It is the password that will be used for deploying
// the freeipa workload.
func SecretForIDM(m *v1alpha1.IDM, adminPassword string, dmPassword string) *corev1.Secret {
	if adminPassword == "" {
		adminPassword = GenerateRandomPassword()
	}
	if dmPassword == "" {
		dmPassword = GenerateRandomPassword()
	}
	adminPassword = b64.StdEncoding.EncodeToString([]byte(adminPassword))
	dmPassword = b64.StdEncoding.EncodeToString([]byte(dmPassword))
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GetSecretName(m),
			Namespace: m.Namespace,
			Labels:    LabelsForIDM(m),
		},
		Immutable: pointy.Bool(true),
		Type:      corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"IPA_ADMIN_PASSWORD": adminPassword,
			"IPA_DM_PASSWORD":    dmPassword,
		},
	}
	return secret
}
