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

package v1

import (
	internalinterfaces "go.f110.dev/heimdallr/operator/pkg/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Certificates returns a CertificateInformer.
	Certificates() CertificateInformer
	// CertificateRequests returns a CertificateRequestInformer.
	CertificateRequests() CertificateRequestInformer
	// ClusterIssuers returns a ClusterIssuerInformer.
	ClusterIssuers() ClusterIssuerInformer
	// Issuers returns a IssuerInformer.
	Issuers() IssuerInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Certificates returns a CertificateInformer.
func (v *version) Certificates() CertificateInformer {
	return &certificateInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// CertificateRequests returns a CertificateRequestInformer.
func (v *version) CertificateRequests() CertificateRequestInformer {
	return &certificateRequestInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ClusterIssuers returns a ClusterIssuerInformer.
func (v *version) ClusterIssuers() ClusterIssuerInformer {
	return &clusterIssuerInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// Issuers returns a IssuerInformer.
func (v *version) Issuers() IssuerInformer {
	return &issuerInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
