package cluster

import (
	"bytes"
	"context"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"strings"
)

func (k *Kubernetes) CreateByFile(filebytes []byte) error {

	//filebytes, err := ioutil.ReadFile(filepath)
	//if err != nil {
	//	return err
	//}
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(filebytes), 1000)
	gr, err := restmapper.GetAPIGroupResources(k.Clientset.Discovery())
	if err != nil {
		return err
	}
	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			break
		}
		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}
		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			//if unstructuredObj.GetNamespace() == "" {
			//	unstructuredObj.SetNamespace("cert-manager")
			//}
			dri = k.DynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = k.DynamicClient.Resource(mapping.Resource)
		}
		createdObj, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Warn().Msgf(err.Error())
			} else {
				return err
			}
		} else {
			log.Info().Msgf("%s/%s created", createdObj.GetKind(), createdObj.GetName())
		}
	}
	return nil
}

func (k *Kubernetes) DeleteByFile(filebytes []byte) error {

	//filebytes, err := ioutil.ReadFile(filepath)
	//if err != nil {
	//	return err
	//}
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(filebytes), 1000)
	discovery := k.Clientset.Discovery()
	gr, err := restmapper.GetAPIGroupResources(discovery)
	if err != nil {
		return err
	}

	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			break
		}
		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return err
		}
		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			//if unstructuredObj.GetNamespace() == "" {
			//	unstructuredObj.SetNamespace("cert-manager")
			//}
			dri = k.DynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = k.DynamicClient.Resource(mapping.Resource)
		}
		err = dri.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no matches") {
				log.Warn().Msgf(err.Error())
			} else {
				return err
			}
		} else {
			log.Info().Msgf("%s/%s deleted", unstructuredObj.GetKind(), unstructuredObj.GetName())
		}
	}
	return nil
}
