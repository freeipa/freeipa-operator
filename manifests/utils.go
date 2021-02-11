package helper

import (
	"math/rand"

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

// GetMasterPodName Return the MasterPodName for the requested IDM resource
func GetMasterPodName(m *v1alpha1.IDM) string {
	return m.Name + "-master"
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

// GetCaSubject Get the CA Subject for the POD
func GetCaSubject(m *v1alpha1.IDM) string {
	return m.Namespace + "-" + RandStringBytes(7)
}

// GetWebServiceName Return the MasterPodName for the requested IDM resource
func GetWebServiceName(m *v1alpha1.IDM) string {
	return m.Name + "-web"
}
