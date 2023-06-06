package cluster

import (
	"context"
	"encoding/json"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// GetNamespace ...
func (k *Kubernetes) GetNamespace(name string) (*coreV1.Namespace, error) {
	return k.Clientset.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
}
func (k *Kubernetes) CreateNamespace(name string) (*coreV1.Namespace, error) {
	return k.Clientset.CoreV1().Namespaces().Create(context.TODO(), &coreV1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}, metav1.CreateOptions{})
}

// PatchNamespaceLabel ...
func (k *Kubernetes) PatchNamespaceLabel(name, key, value string) (*coreV1.Namespace, error) {
	namespace, err := k.GetNamespace(name)
	if nil != err {
		namespace, err = k.CreateNamespace(name)
		if err != nil {
			return nil, err
		}
	}
	labels := namespace.Labels
	if nil == labels {
		labels = make(map[string]string, 1)
	}
	labels[key] = value
	patchData := map[string]interface{}{"metadata": map[string]map[string]string{"labels": labels}}
	playLoadBytes, _ := json.Marshal(patchData)
	return k.Clientset.CoreV1().Namespaces().Patch(context.TODO(), namespace.Name, types.StrategicMergePatchType, playLoadBytes, metav1.PatchOptions{})
}
