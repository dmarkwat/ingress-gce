package operator

import (
	"fmt"

	backendconfigv1beta1 "k8s.io/ingress-gce/pkg/apis/backendconfig/v1beta1"
	frontendconfigv1beta1 "k8s.io/ingress-gce/pkg/apis/frontendconfig/v1beta1"

	api_v1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
)

// Ingresses returns the wrapper
func Ingresses(i []*v1beta1.Ingress) *IngressesOperator {
	return &IngressesOperator{i: i}
}

// IngressesOperator is an operator wrapper for a list of Ingresses.
type IngressesOperator struct {
	i []*v1beta1.Ingress
}

// AsList returns the underlying list of Ingresses.
func (op *IngressesOperator) AsList() []*v1beta1.Ingress {
	if op.i == nil {
		return []*v1beta1.Ingress{}
	}
	return op.i
}

// Filter the list of Ingresses based on a predicate.
func (op *IngressesOperator) Filter(f func(*v1beta1.Ingress) bool) *IngressesOperator {
	var i []*v1beta1.Ingress
	for _, ing := range op.i {
		if f(ing) {
			i = append(i, ing)
		}
	}
	return Ingresses(i)
}

// ReferencesService returns the Ingresses that references the given Service.
func (op *IngressesOperator) ReferencesService(svc *api_v1.Service) *IngressesOperator {
	dupes := map[string]bool{}

	var i []*v1beta1.Ingress
	for _, ing := range op.i {
		key := fmt.Sprintf("%s/%s", ing.Namespace, ing.Name)
		if doesIngressReferenceService(ing, svc) && !dupes[key] {
			i = append(i, ing)
			dupes[key] = true
		}
	}
	return Ingresses(i)
}

// ReferencesBackendConfig returns the Ingresses that references the given BackendConfig.
func (op *IngressesOperator) ReferencesBackendConfig(beConfig *backendconfigv1beta1.BackendConfig, svcsOp *ServicesOperator) *IngressesOperator {
	dupes := map[string]bool{}

	var i []*v1beta1.Ingress
	svcs := svcsOp.ReferencesBackendConfig(beConfig).AsList()
	for _, ing := range op.i {
		for _, svc := range svcs {
			key := fmt.Sprintf("%s/%s", ing.Namespace, ing.Name)
			if doesIngressReferenceService(ing, svc) && !dupes[key] {
				i = append(i, ing)
				dupes[key] = true
			}
		}
	}
	return Ingresses(i)
}

// ReferencesFrontendConfig returns the Ingresses that reference the given FrontendConfig.
func (op *IngressesOperator) ReferencesFrontendConfig(feConfig *frontendconfigv1beta1.FrontendConfig) *IngressesOperator {
	dupes := map[string]bool{}

	var i []*v1beta1.Ingress
	for _, ing := range op.i {
		key := fmt.Sprintf("%s/%s", ing.Namespace, ing.Name)
		if doesIngressReferenceFrontendConfig(ing, feConfig) && !dupes[key] {
			i = append(i, ing)
			dupes[key] = true
		}
	}
	return Ingresses(i)
}
