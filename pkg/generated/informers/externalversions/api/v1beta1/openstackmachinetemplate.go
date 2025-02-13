/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	apiv1beta1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	clientset "sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/clientset/clientset"
	internalinterfaces "sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/informers/externalversions/internalinterfaces"
	v1beta1 "sigs.k8s.io/cluster-api-provider-openstack/pkg/generated/listers/api/v1beta1"
)

// OpenStackMachineTemplateInformer provides access to a shared informer and lister for
// OpenStackMachineTemplates.
type OpenStackMachineTemplateInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.OpenStackMachineTemplateLister
}

type openStackMachineTemplateInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewOpenStackMachineTemplateInformer constructs a new informer for OpenStackMachineTemplate type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewOpenStackMachineTemplateInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredOpenStackMachineTemplateInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredOpenStackMachineTemplateInformer constructs a new informer for OpenStackMachineTemplate type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredOpenStackMachineTemplateInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.InfrastructureV1beta1().OpenStackMachineTemplates(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.InfrastructureV1beta1().OpenStackMachineTemplates(namespace).Watch(context.TODO(), options)
			},
		},
		&apiv1beta1.OpenStackMachineTemplate{},
		resyncPeriod,
		indexers,
	)
}

func (f *openStackMachineTemplateInformer) defaultInformer(client clientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredOpenStackMachineTemplateInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *openStackMachineTemplateInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiv1beta1.OpenStackMachineTemplate{}, f.defaultInformer)
}

func (f *openStackMachineTemplateInformer) Lister() v1beta1.OpenStackMachineTemplateLister {
	return v1beta1.NewOpenStackMachineTemplateLister(f.Informer().GetIndexer())
}
