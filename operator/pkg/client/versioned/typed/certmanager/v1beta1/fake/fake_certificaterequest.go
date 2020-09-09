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
	"context"

	v1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCertificateRequests implements CertificateRequestInterface
type FakeCertificateRequests struct {
	Fake *FakeCertmanagerV1beta1
	ns   string
}

var certificaterequestsResource = schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1beta1", Resource: "certificaterequests"}

var certificaterequestsKind = schema.GroupVersionKind{Group: "cert-manager.io", Version: "v1beta1", Kind: "CertificateRequest"}

// Get takes name of the certificateRequest, and returns the corresponding certificateRequest object, and an error if there is any.
func (c *FakeCertificateRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.CertificateRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(certificaterequestsResource, c.ns, name), &v1beta1.CertificateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CertificateRequest), err
}

// List takes label and field selectors, and returns the list of CertificateRequests that match those selectors.
func (c *FakeCertificateRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.CertificateRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(certificaterequestsResource, certificaterequestsKind, c.ns, opts), &v1beta1.CertificateRequestList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.CertificateRequestList{ListMeta: obj.(*v1beta1.CertificateRequestList).ListMeta}
	for _, item := range obj.(*v1beta1.CertificateRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested certificateRequests.
func (c *FakeCertificateRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(certificaterequestsResource, c.ns, opts))

}

// Create takes the representation of a certificateRequest and creates it.  Returns the server's representation of the certificateRequest, and an error, if there is any.
func (c *FakeCertificateRequests) Create(ctx context.Context, certificateRequest *v1beta1.CertificateRequest, opts v1.CreateOptions) (result *v1beta1.CertificateRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(certificaterequestsResource, c.ns, certificateRequest), &v1beta1.CertificateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CertificateRequest), err
}

// Update takes the representation of a certificateRequest and updates it. Returns the server's representation of the certificateRequest, and an error, if there is any.
func (c *FakeCertificateRequests) Update(ctx context.Context, certificateRequest *v1beta1.CertificateRequest, opts v1.UpdateOptions) (result *v1beta1.CertificateRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(certificaterequestsResource, c.ns, certificateRequest), &v1beta1.CertificateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CertificateRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCertificateRequests) UpdateStatus(ctx context.Context, certificateRequest *v1beta1.CertificateRequest, opts v1.UpdateOptions) (*v1beta1.CertificateRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(certificaterequestsResource, "status", c.ns, certificateRequest), &v1beta1.CertificateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CertificateRequest), err
}

// Delete takes name of the certificateRequest and deletes it. Returns an error if one occurs.
func (c *FakeCertificateRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(certificaterequestsResource, c.ns, name), &v1beta1.CertificateRequest{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCertificateRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(certificaterequestsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.CertificateRequestList{})
	return err
}

// Patch applies the patch and returns the patched certificateRequest.
func (c *FakeCertificateRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.CertificateRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(certificaterequestsResource, c.ns, name, pt, data, subresources...), &v1beta1.CertificateRequest{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CertificateRequest), err
}
