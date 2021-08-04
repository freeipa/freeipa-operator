package manifests

import (
	"math/rand"
	"strings"

	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/google/uuid"
)

// LabelsForIDM Returns the labels for selecting the resources
// belonging to the given memcached CR name.
func LabelsForIDM(m *v1alpha1.IDM) map[string]string {
	return map[string]string{
		"app": "idm",
		"idm": m.Name,
	}
}

// GenerateSystemdUUIDString Generate a UUID string that can be used for system-uuid
func GenerateSystemdUUIDString() string {
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx => xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	str := uuid.New().String()
	value := str[:8] + str[9:13] + str[14:18] + str[19:23] + str[24:]
	return value
}

// GetMainPodName Return the MasterPodName for the requested IDM resource
func GetMainPodName(m *v1alpha1.IDM) string {
	return m.Name + "-main"
}

// GetMainStatefulsetName Return the MainStatefulsetName for the requested IDM resource
func GetMainStatefulsetName(m *v1alpha1.IDM) string {
	return m.Name + "-main"
}

// RandStringBytes Read a string of n random chars
func RandStringBytes(n int) string {
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// GetRealm Get the REALM for the POD
func GetRealm(m *v1alpha1.IDM, baseDomain string) string {
	return "APPS." + strings.ToUpper(baseDomain)
}

// GetIpaServerHostname Get the hostname passed to ipa installation
func GetIpaServerHostname(m *v1alpha1.IDM, baseDomain string) string {
	return m.Namespace + ".apps." + baseDomain
}

// GetCaSubject Get the CA Subject for the POD
func GetCaSubject(m *v1alpha1.IDM, baseDomain string) string {
	cn := m.Namespace + "-" + RandStringBytes(7)
	o := GetRealm(m, baseDomain)
	return "CN=" + cn + ", O=" + o
}

// GetWebServiceName Return the MasterPodName for the requested IDM resource
func GetWebServiceName(m *v1alpha1.IDM) string {
	return m.Name + "-web"
}

// GetServiceAccountName Return the ServiceAccount name for the
// requested IDM resource.
func GetServiceAccountName(m *v1alpha1.IDM) string {
	return "idm"
}

// GetRoleName Return the Role name for the requested IDM resource.
func GetRoleName(m *v1alpha1.IDM) string {
	return "idm"
}

// GetRoleBindingName Return the Role name for the requested IDM resource.
func GetRoleBindingName(m *v1alpha1.IDM) string {
	return "idm"
}

// GetMainPersistentVolumeClaim Return the name for the PersistentVolumClaim
// used by the main pod.
func GetMainPersistentVolumeClaimName(m *v1alpha1.IDM) string {
	return m.Name + "-main"
}
