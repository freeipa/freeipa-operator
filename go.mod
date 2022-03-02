// https://blog.golang.org/using-go-modules
module github.com/freeipa/freeipa-operator

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.2.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/openlyinc/pointy v1.1.2
	github.com/openshift/api v0.0.0-20210309190949-7d6cac66d2a4
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	sigs.k8s.io/controller-runtime v0.10.0
)
