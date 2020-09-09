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

package v1beta1

import (
	v1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
	"go.f110.dev/heimdallr/operator/pkg/client/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type CertmanagerV1beta1Interface interface {
	RESTClient() rest.Interface
	CertificatesGetter
	CertificateRequestsGetter
	ClusterIssuersGetter
	IssuersGetter
}

// CertmanagerV1beta1Client is used to interact with features provided by the cert-manager.io group.
type CertmanagerV1beta1Client struct {
	restClient rest.Interface
}

func (c *CertmanagerV1beta1Client) Certificates(namespace string) CertificateInterface {
	return newCertificates(c, namespace)
}

func (c *CertmanagerV1beta1Client) CertificateRequests(namespace string) CertificateRequestInterface {
	return newCertificateRequests(c, namespace)
}

func (c *CertmanagerV1beta1Client) ClusterIssuers() ClusterIssuerInterface {
	return newClusterIssuers(c)
}

func (c *CertmanagerV1beta1Client) Issuers(namespace string) IssuerInterface {
	return newIssuers(c, namespace)
}

// NewForConfig creates a new CertmanagerV1beta1Client for the given config.
func NewForConfig(c *rest.Config) (*CertmanagerV1beta1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &CertmanagerV1beta1Client{client}, nil
}

// NewForConfigOrDie creates a new CertmanagerV1beta1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *CertmanagerV1beta1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new CertmanagerV1beta1Client for the given RESTClient.
func New(c rest.Interface) *CertmanagerV1beta1Client {
	return &CertmanagerV1beta1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1beta1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *CertmanagerV1beta1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
