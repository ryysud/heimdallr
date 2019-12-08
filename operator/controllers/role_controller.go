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

package controllers

import (
	"context"

	proxyv1 "github.com/f110/lagrangian-proxy/operator/api/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RoleReconciler reconciles a Role object
type RoleReconciler struct {
	client.Client
	Log               logr.Logger
	Scheme            *runtime.Scheme
	ProcessRepository *ProcessRepository
}

// +kubebuilder:rbac:groups=proxy.f110.dev,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=proxy.f110.dev,resources=roles/status,verbs=get;update;patch

func (r *RoleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&proxyv1.Role{}).
		Complete(r)
}

func (r *RoleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	role := &proxyv1.Role{}
	if err := r.Get(context.Background(), req.NamespacedName, role); err != nil {
		return ctrl.Result{}, err
	}

	defList := &proxyv1.ProxyList{}
	if err := r.List(context.Background(), defList); err != nil {
		return ctrl.Result{}, err
	}

	targets := make([]proxyv1.Proxy, 0)
Item:
	for _, v := range defList.Items {
		for k := range v.Spec.RoleSelector.MatchLabels {
			if value, ok := role.ObjectMeta.Labels[k]; !ok && v.Spec.RoleSelector.MatchLabels[k] != value {
				continue Item
			}
		}

		targets = append(targets, v)
	}

	for _, v := range targets {
		lp := r.ProcessRepository.Get(&v)
		if err := r.reconcileConfig(lp); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *RoleReconciler) reconcileConfig(lp *LagrangianProxy) error {
	lp.Lock()
	defer lp.Unlock()

	configMap, err := lp.ReverseProxyConfig()
	if err != nil {
		return err
	}
	orig := configMap.DeepCopy()
	_, err = ctrl.CreateOrUpdate(context.Background(), r, configMap, func() error {
		configMap.Data = orig.Data

		return ctrl.SetControllerReference(lp.Object, configMap, r.Scheme)
	})
	if err != nil {
		return err
	}

	return nil
}
