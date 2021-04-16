package helper

import (
	"github.com/freeipa/freeipa-operator/api/v1alpha1"
	"github.com/openlyinc/pointy"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// RouteForIDM Create the Route manifest for this IDM resource
// clusterDomain It is the subdomain associated to the cluster
//
func RouteForIDM(m *v1alpha1.IDM, clusterDomain string) *routev1.Route {
	route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Annotations: map[string]string{
				// https://docs.openshift.com/container-platform/4.6/networking/routes/route-configuration.html
				"openshift.io/host.generated":             "true",
				"haproxy.router.openshift.io/timeout":     "2s",
				"haproxy.router.openshift.io/hsts_header": "max-age=31536000;includeSubDomains;preload",
			},
			Labels: LabelsForIDM(m),
		},
		Spec: routev1.RouteSpec{
			Host: m.Namespace + ".apps." + clusterDomain,
			Port: &routev1.RoutePort{
				TargetPort: intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "https-tcp",
				},
			},
			To: routev1.RouteTargetReference{
				Kind:   "Service",
				Name:   GetWebServiceName(m),
				Weight: pointy.Int32(100),
			},
			TLS: &routev1.TLSConfig{
				Termination: routev1.TLSTerminationPassthrough,
			},
			WildcardPolicy: routev1.WildcardPolicyNone,
		},
	}
	return route
}
