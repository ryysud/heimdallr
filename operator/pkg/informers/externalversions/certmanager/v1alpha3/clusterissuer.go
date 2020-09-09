/*

MIT License

Copyright (c) 2019 Fumihiro Ito

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha3

import (
	"context"
	time "time"

	certmanagerv1alpha3 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha3"
	versioned "go.f110.dev/heimdallr/operator/pkg/client/versioned"
	internalinterfaces "go.f110.dev/heimdallr/operator/pkg/informers/externalversions/internalinterfaces"
	v1alpha3 "go.f110.dev/heimdallr/operator/pkg/listers/certmanager/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ClusterIssuerInformer provides access to a shared informer and lister for
// ClusterIssuers.
type ClusterIssuerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha3.ClusterIssuerLister
}

type clusterIssuerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewClusterIssuerInformer constructs a new informer for ClusterIssuer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewClusterIssuerInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredClusterIssuerInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredClusterIssuerInformer constructs a new informer for ClusterIssuer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredClusterIssuerInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CertmanagerV1alpha3().ClusterIssuers().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CertmanagerV1alpha3().ClusterIssuers().Watch(context.TODO(), options)
			},
		},
		&certmanagerv1alpha3.ClusterIssuer{},
		resyncPeriod,
		indexers,
	)
}

func (f *clusterIssuerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredClusterIssuerInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *clusterIssuerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&certmanagerv1alpha3.ClusterIssuer{}, f.defaultInformer)
}

func (f *clusterIssuerInformer) Lister() v1alpha3.ClusterIssuerLister {
	return v1alpha3.NewClusterIssuerLister(f.Informer().GetIndexer())
}
