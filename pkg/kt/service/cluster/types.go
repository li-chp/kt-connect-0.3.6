package cluster

import (
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	extV1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// KubernetesInterface kubernetes interface
type KubernetesInterface interface {
	GetPod(name string, namespace string) (*coreV1.Pod, error)
	GetAllPodInNamespace(namespace string) (*coreV1.PodList, error)
	GetPodsByLabel(labels map[string]string, namespace string) (*coreV1.PodList, error)
	UpdatePod(pod *coreV1.Pod) (*coreV1.Pod, error)
	RemovePod(name, namespace string) error
	GetOrCreateShadow(name string, labels, annotations, envs map[string]string, portsToExpose string, portNameDict map[int]string) (string, string, string, error)
	CreateRouterPod(name string, labels, annotations map[string]string, ports map[int]int) (*coreV1.Pod, error)
	CreateRectifierPod(name string) (*coreV1.Pod, error)
	UpdatePodHeartBeat(name, namespace string)
	WaitPodReady(name, namespace string, timeoutSec int) (*coreV1.Pod, error)
	WaitPodTerminate(name, namespace string) (*coreV1.Pod, error)
	WatchPod(name, namespace string, fAdd, fDel, fMod func(*coreV1.Pod))
	ExecInPod(containerName, podName, namespace string, cmd ...string) (string, string, error)
	AddEphemeralContainer(containerName, podName string, envs map[string]string) (string, error)
	RemoveEphemeralContainer(containerName, podName string, namespace string) error
	IncreasePodRef(name, namespace string) error
	DecreasePodRef(name, namespace string) (bool, error)

	GetDeployment(name string, namespace string) (*appV1.Deployment, error)
	GetDeploymentsByLabel(labels map[string]string, namespace string) (*appV1.DeploymentList, error)
	GetAllDeploymentInNamespace(namespace string) (*appV1.DeploymentList, error)
	UpdateDeployment(deployment *appV1.Deployment) (*appV1.Deployment, error)
	RemoveDeployment(name, namespace string) error
	IncreaseDeploymentRef(name, namespace string) error
	DecreaseDeploymentRef(name, namespace string) (bool, error)
	ScaleTo(deployment, namespace string, replicas *int32) (err error)

	GetService(name, namespace string) (*coreV1.Service, error)
	GetServicesBySelector(matchLabels map[string]string, namespace string) ([]coreV1.Service, error)
	GetAllServiceInNamespace(namespace string) (*coreV1.ServiceList, error)
	GetServicesByLabel(labels map[string]string, namespace string) (*coreV1.ServiceList, error)
	CreateService(metaAndSpec *SvcMetaAndSpec) (*coreV1.Service, error)
	UpdateService(svc *coreV1.Service) (*coreV1.Service, error)
	RemoveService(name, namespace string) (err error)
	UpdateServiceHeartBeat(name, namespace string)
	WatchService(name, namespace string, fAdd, fDel, fMod func(*coreV1.Service))

	GetConfigMap(name, namespace string) (*coreV1.ConfigMap, error)
	GetConfigMapsByLabel(labels map[string]string, namespace string) (*coreV1.ConfigMapList, error)
	RemoveConfigMap(name, namespace string) (err error)
	UpdateConfigMapHeartBeat(name, namespace string)

	GetAllIngressInNamespace(namespace string) (*extV1.IngressList, error)

	GetKtResources(namespace string) ([]coreV1.Pod, []coreV1.ConfigMap, []appV1.Deployment, []coreV1.Service, error)
	GetAllNamespaces() (*coreV1.NamespaceList, error)
	ClusterCidr(namespace string) (cidr []string, excludeCidr []string)

	GetVirtualServiceList(namespace string) (*v1alpha3.VirtualServiceList, error)
	GetVirtualService(name, namespace string) (*v1alpha3.VirtualService, error)
	GetDestinationRule(name, namespace string) (*v1alpha3.DestinationRule, error)
	PatchVirtualService(name, service, namespace, op, meshKey, meshVersion string) (*v1alpha3.VirtualService, error)
	PatchDestinationRule(name, namespace string, op, meshKey, meshVersion string) (*v1alpha3.DestinationRule, error)

	CreateByFile(filebyte []byte, namespace string) error
	DeleteByFile(filebyte []byte, namespace string) error

	GetNamespace(name string) (*coreV1.Namespace, error)
	PatchNamespaceLabel(namespace, key, value string) (*coreV1.Namespace, error)
}

// Kubernetes implements KubernetesInterface
type Kubernetes struct {
	Clientset     kubernetes.Interface
	IstioClient   versionedclient.Interface
	DynamicClient dynamic.Interface
}

// Cli the singleton type
var instance *Kubernetes

// Ins get singleton instance
func Ins() KubernetesInterface {
	if instance == nil {
		instance = &Kubernetes{
			Clientset:     opt.Store.Clientset,
			IstioClient:   opt.Store.IstioClient,
			DynamicClient: opt.Store.DynamicClient,
		}
	}
	return instance
}
