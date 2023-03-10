package options

import (
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var Store = &RuntimeStore{}

// RuntimeStore ...
type RuntimeStore struct {
	// Clientset for kubernetes operation
	Clientset kubernetes.Interface
	// IstioClient for kubernetes operation
	IstioClient versionedclient.Interface
	// IstioClient for kubernetes operation
	DynamicClient dynamic.Interface
	// RestConfig kubectl config
	RestConfig *rest.Config
	// Version et version
	Version string
	// Component current sub-command (connect, exchange, mesh or preview)
	Component string
	// Shadow pod name
	Shadow string
	// Router pod name
	Router string
	// Mesh version of mesh pod
	Mesh string
	// Mesh version of mesh pod
	MeshDebug string
	// Origin the origin deployment or service name
	Origin string
	// Replicas the origin replicas
	Replicas int32
	// Service exposed service name
	Service string
	// VS
	VirtualServicePatch bool
	// DR
	DestinationRulePatch bool
	// isIpv6Cluster
	Ipv6Cluster bool
}
