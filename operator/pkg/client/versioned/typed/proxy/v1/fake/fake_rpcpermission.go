/*

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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	proxyv1 "github.com/f110/lagrangian-proxy/operator/pkg/api/proxy/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRpcPermissions implements RpcPermissionInterface
type FakeRpcPermissions struct {
	Fake *FakeProxyV1
	ns   string
}

var rpcpermissionsResource = schema.GroupVersionResource{Group: "proxy.f110.dev", Version: "v1", Resource: "rpcpermissions"}

var rpcpermissionsKind = schema.GroupVersionKind{Group: "proxy.f110.dev", Version: "v1", Kind: "RpcPermission"}

// Get takes name of the rpcPermission, and returns the corresponding rpcPermission object, and an error if there is any.
func (c *FakeRpcPermissions) Get(name string, options v1.GetOptions) (result *proxyv1.RpcPermission, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(rpcpermissionsResource, c.ns, name), &proxyv1.RpcPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*proxyv1.RpcPermission), err
}

// List takes label and field selectors, and returns the list of RpcPermissions that match those selectors.
func (c *FakeRpcPermissions) List(opts v1.ListOptions) (result *proxyv1.RpcPermissionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(rpcpermissionsResource, rpcpermissionsKind, c.ns, opts), &proxyv1.RpcPermissionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &proxyv1.RpcPermissionList{ListMeta: obj.(*proxyv1.RpcPermissionList).ListMeta}
	for _, item := range obj.(*proxyv1.RpcPermissionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested rpcPermissions.
func (c *FakeRpcPermissions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(rpcpermissionsResource, c.ns, opts))

}

// Create takes the representation of a rpcPermission and creates it.  Returns the server's representation of the rpcPermission, and an error, if there is any.
func (c *FakeRpcPermissions) Create(rpcPermission *proxyv1.RpcPermission) (result *proxyv1.RpcPermission, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(rpcpermissionsResource, c.ns, rpcPermission), &proxyv1.RpcPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*proxyv1.RpcPermission), err
}

// Update takes the representation of a rpcPermission and updates it. Returns the server's representation of the rpcPermission, and an error, if there is any.
func (c *FakeRpcPermissions) Update(rpcPermission *proxyv1.RpcPermission) (result *proxyv1.RpcPermission, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(rpcpermissionsResource, c.ns, rpcPermission), &proxyv1.RpcPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*proxyv1.RpcPermission), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRpcPermissions) UpdateStatus(rpcPermission *proxyv1.RpcPermission) (*proxyv1.RpcPermission, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(rpcpermissionsResource, "status", c.ns, rpcPermission), &proxyv1.RpcPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*proxyv1.RpcPermission), err
}

// Delete takes name of the rpcPermission and deletes it. Returns an error if one occurs.
func (c *FakeRpcPermissions) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(rpcpermissionsResource, c.ns, name), &proxyv1.RpcPermission{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRpcPermissions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(rpcpermissionsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &proxyv1.RpcPermissionList{})
	return err
}

// Patch applies the patch and returns the patched rpcPermission.
func (c *FakeRpcPermissions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *proxyv1.RpcPermission, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(rpcpermissionsResource, c.ns, name, pt, data, subresources...), &proxyv1.RpcPermission{})

	if obj == nil {
		return nil, err
	}
	return obj.(*proxyv1.RpcPermission), err
}
