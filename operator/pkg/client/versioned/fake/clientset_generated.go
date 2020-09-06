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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	clientset "go.f110.dev/heimdallr/operator/pkg/client/versioned"
	certmanagerv1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/certmanager/v1"
	fakecertmanagerv1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/certmanager/v1/fake"
	certmanagerv1alpha2 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/certmanager/v1alpha2"
	fakecertmanagerv1alpha2 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/certmanager/v1alpha2/fake"
	certmanagerv1alpha3 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/certmanager/v1alpha3"
	fakecertmanagerv1alpha3 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/certmanager/v1alpha3/fake"
	certmanagerv1beta1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/certmanager/v1beta1"
	fakecertmanagerv1beta1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/certmanager/v1beta1/fake"
	etcdv1alpha1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/etcd/v1alpha1"
	fakeetcdv1alpha1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/etcd/v1alpha1/fake"
	monitoringv1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/monitoring/v1"
	fakemonitoringv1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/monitoring/v1/fake"
	proxyv1alpha1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/proxy/v1alpha1"
	fakeproxyv1alpha1 "go.f110.dev/heimdallr/operator/pkg/client/versioned/typed/proxy/v1alpha1/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &Clientset{tracker: o}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})

	return cs
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
	tracker   testing.ObjectTracker
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *Clientset) Tracker() testing.ObjectTracker {
	return c.tracker
}

var _ clientset.Interface = &Clientset{}

// CertmanagerV1 retrieves the CertmanagerV1Client
func (c *Clientset) CertmanagerV1() certmanagerv1.CertmanagerV1Interface {
	return &fakecertmanagerv1.FakeCertmanagerV1{Fake: &c.Fake}
}

// CertmanagerV1alpha2 retrieves the CertmanagerV1alpha2Client
func (c *Clientset) CertmanagerV1alpha2() certmanagerv1alpha2.CertmanagerV1alpha2Interface {
	return &fakecertmanagerv1alpha2.FakeCertmanagerV1alpha2{Fake: &c.Fake}
}

// CertmanagerV1alpha3 retrieves the CertmanagerV1alpha3Client
func (c *Clientset) CertmanagerV1alpha3() certmanagerv1alpha3.CertmanagerV1alpha3Interface {
	return &fakecertmanagerv1alpha3.FakeCertmanagerV1alpha3{Fake: &c.Fake}
}

// CertmanagerV1beta1 retrieves the CertmanagerV1beta1Client
func (c *Clientset) CertmanagerV1beta1() certmanagerv1beta1.CertmanagerV1beta1Interface {
	return &fakecertmanagerv1beta1.FakeCertmanagerV1beta1{Fake: &c.Fake}
}

// EtcdV1alpha1 retrieves the EtcdV1alpha1Client
func (c *Clientset) EtcdV1alpha1() etcdv1alpha1.EtcdV1alpha1Interface {
	return &fakeetcdv1alpha1.FakeEtcdV1alpha1{Fake: &c.Fake}
}

// MonitoringV1 retrieves the MonitoringV1Client
func (c *Clientset) MonitoringV1() monitoringv1.MonitoringV1Interface {
	return &fakemonitoringv1.FakeMonitoringV1{Fake: &c.Fake}
}

// ProxyV1alpha1 retrieves the ProxyV1alpha1Client
func (c *Clientset) ProxyV1alpha1() proxyv1alpha1.ProxyV1alpha1Interface {
	return &fakeproxyv1alpha1.FakeProxyV1alpha1{Fake: &c.Fake}
}
