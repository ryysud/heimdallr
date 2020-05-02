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
// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/f110/lagrangian-proxy/operator/pkg/api/proxy/v1"
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
