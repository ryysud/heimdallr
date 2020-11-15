package controllers

import (
	"fmt"
	"sort"

	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"sigs.k8s.io/yaml"

	proxyv1alpha1 "go.f110.dev/heimdallr/operator/pkg/api/proxy/v1alpha1"
	"go.f110.dev/heimdallr/pkg/config/configv2"
)

type ConfigConverter struct {
}

func (ConfigConverter) Proxy(backends []*proxyv1alpha1.Backend, serviceLister listers.ServiceLister) ([]byte, error) {
	proxies := make([]*configv2.Backend, 0, len(backends))
	for _, v := range backends {
		_, virtualDashboard := v.Labels[labelKeyVirtualDashboard]

		var service *corev1.Service
		if !virtualDashboard && v.Spec.Upstream == "" {
			svc, err := findService(serviceLister, v.Spec.ServiceSelector, v.Namespace)
			if err != nil {
				// At this time, ignore error
				continue
			}
			if svc == nil {
				continue
			}
			service = svc
		}
		if v.Spec.Upstream == "" && service == nil {
			continue
		}

		name := v.Name + "." + v.Spec.Layer
		if v.Spec.Layer == "" {
			name = v.Name
		}
		upstream := v.Spec.Upstream
		if upstream == "" {
			for _, p := range service.Spec.Ports {
				if p.Name == v.Spec.ServiceSelector.Port {
					scheme := v.Spec.ServiceSelector.Scheme
					if scheme == "" {
						switch p.Name {
						case "http", "https":
							scheme = p.Name
						}
					}

					upstream = fmt.Sprintf("%s://%s.%s.svc:%d", scheme, service.Name, service.Namespace, p.Port)
					break
				}
			}
		}
		b := &configv2.Backend{
			Name:             name,
			FQDN:             v.Spec.FQDN,
			Upstream:         upstream,
			Permissions:      toConfigPermissions(v.Spec.Permissions),
			WebHook:          v.Spec.Webhook,
			WebHookPath:      v.Spec.WebhookPath,
			Agent:            v.Spec.Agent,
			Socket:           v.Spec.Socket,
			AllowRootUser:    v.Spec.AllowRootUser,
			DisableAuthn:     v.Spec.DisableAuthn,
			InsecureUpstream: v.Spec.Insecure,
			AllowHttp:        v.Spec.AllowHttp,
		}
		if v.Spec.SocketTimeout != nil {
			b.SocketTimeout = &configv2.Duration{Duration: v.Spec.SocketTimeout.Duration}
		}
		if v.Spec.MaxSessionDuration != nil {
			b.MaxSessionDuration = &configv2.Duration{Duration: v.Spec.MaxSessionDuration.Duration}
		}
		proxies = append(proxies, b)
	}
	sort.Slice(proxies, func(i, j int) bool {
		return proxies[i].Name < proxies[j].Name
	})
	proxyBinary, err := yaml.Marshal(proxies)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return proxyBinary, nil
}

func (ConfigConverter) Role(backends []*proxyv1alpha1.Backend, roleList []*proxyv1alpha1.Role, roleBindings []*proxyv1alpha1.RoleBinding) ([]byte, error) {
	backendMap := make(map[string]*proxyv1alpha1.Backend)
	for _, v := range backends {
		backendMap[v.Namespace+"/"+v.Name] = v
	}

	roles := make([]*configv2.Role, len(roleList))
	for i, role := range roleList {
		bindings := make([]*configv2.Binding, 0)

		matchedBindings := RoleBindings(roleBindings).Select(func(binding *proxyv1alpha1.RoleBinding) bool {
			if binding.RoleRef.Name != role.Name {
				return false
			}
			if binding.RoleRef.Namespace != "" && binding.RoleRef.Namespace == role.Namespace {
				return true
			}
			if binding.RoleRef.Namespace == "" && binding.ObjectMeta.Namespace == role.Namespace {
				return true
			}

			return false
		})
		if role.Spec.AllowDashboard {
			matchedBindings = append(matchedBindings, &proxyv1alpha1.RoleBinding{
				Subjects: []proxyv1alpha1.Subject{
					{
						Kind:       "Backend",
						Name:       "dashboard",
						Permission: "all",
					},
				},
			})
		}

		sort.Slice(matchedBindings, func(i, j int) bool {
			return matchedBindings[i].Name < matchedBindings[j].Name
		})
		for _, binding := range matchedBindings {
			for _, subject := range binding.Subjects {
				switch subject.Kind {
				case "Backend":
					namespace := role.Namespace
					if subject.Namespace != "" {
						namespace = subject.Namespace
					}
					backendHost := ""
					if bn, ok := backendMap[namespace+"/"+subject.Name]; ok {
						backendHost = bn.Name + "." + bn.Spec.Layer
						if bn.Spec.Layer == "" {
							backendHost = bn.Name
						}
					} else {
						continue
					}

					bindings = append(bindings, &configv2.Binding{
						Permission: subject.Permission,
						Backend:    backendHost,
					})
				case "RpcPermission":
					bindings = append(bindings, &configv2.Binding{
						RPC: subject.Name,
					})
				}
			}
		}

		roles[i] = &configv2.Role{
			Name:        role.Name,
			Title:       role.Spec.Title,
			Description: role.Spec.Description,
			Bindings:    bindings,
		}
	}
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Name < roles[j].Name
	})
	roleBinary, err := yaml.Marshal(roles)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return roleBinary, nil
}

func (ConfigConverter) RPCPermission(permissions []*proxyv1alpha1.RpcPermission) ([]byte, error) {
	rpcPermissions := make([]*configv2.RPCPermission, len(permissions))
	for i, v := range permissions {
		rpcPermissions[i] = &configv2.RPCPermission{
			Name:  v.Name,
			Allow: v.Spec.Allow,
		}
	}
	sort.Slice(rpcPermissions, func(i, j int) bool {
		return rpcPermissions[i].Name < rpcPermissions[j].Name
	})
	rpcPermissionBinary, err := yaml.Marshal(rpcPermissions)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return rpcPermissionBinary, nil
}
