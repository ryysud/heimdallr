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
// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "go.f110.dev/heimdallr/operator/pkg/api/proxy/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RpcPermissionLister helps list RpcPermissions.
type RpcPermissionLister interface {
	// List lists all RpcPermissions in the indexer.
	List(selector labels.Selector) (ret []*v1.RpcPermission, err error)
	// RpcPermissions returns an object that can list and get RpcPermissions.
	RpcPermissions(namespace string) RpcPermissionNamespaceLister
	RpcPermissionListerExpansion
}

// rpcPermissionLister implements the RpcPermissionLister interface.
type rpcPermissionLister struct {
	indexer cache.Indexer
}

// NewRpcPermissionLister returns a new RpcPermissionLister.
func NewRpcPermissionLister(indexer cache.Indexer) RpcPermissionLister {
	return &rpcPermissionLister{indexer: indexer}
}

// List lists all RpcPermissions in the indexer.
func (s *rpcPermissionLister) List(selector labels.Selector) (ret []*v1.RpcPermission, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RpcPermission))
	})
	return ret, err
}

// RpcPermissions returns an object that can list and get RpcPermissions.
func (s *rpcPermissionLister) RpcPermissions(namespace string) RpcPermissionNamespaceLister {
	return rpcPermissionNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RpcPermissionNamespaceLister helps list and get RpcPermissions.
type RpcPermissionNamespaceLister interface {
	// List lists all RpcPermissions in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.RpcPermission, err error)
	// Get retrieves the RpcPermission from the indexer for a given namespace and name.
	Get(name string) (*v1.RpcPermission, error)
	RpcPermissionNamespaceListerExpansion
}

// rpcPermissionNamespaceLister implements the RpcPermissionNamespaceLister
// interface.
type rpcPermissionNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all RpcPermissions in the indexer for a given namespace.
func (s rpcPermissionNamespaceLister) List(selector labels.Selector) (ret []*v1.RpcPermission, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RpcPermission))
	})
	return ret, err
}

// Get retrieves the RpcPermission from the indexer for a given namespace and name.
func (s rpcPermissionNamespaceLister) Get(name string) (*v1.RpcPermission, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("rpcpermission"), name)
	}
	return obj.(*v1.RpcPermission), nil
}
