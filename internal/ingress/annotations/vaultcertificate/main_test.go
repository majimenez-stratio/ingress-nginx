/*
Copyright 2018 The Kubernetes Authors.

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

package backendCertVaultPath

import (
	"testing"

	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/ingress-nginx/internal/ingress/annotations/parser"
	"k8s.io/ingress-nginx/internal/ingress/resolver"

	"k8s.io/apimachinery/pkg/util/intstr"
)

func buildIngress() *networking.Ingress {
	return &networking.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: networking.IngressSpec{
			Backend: &networking.IngressBackend{
				ServiceName: "default-backend",
				ServicePort: intstr.FromInt(80),
			},
		},
	}
}
func TestParseInvalidAnnotations(t *testing.T) {
	ing := buildIngress()

	// Test no annotations set
	i, err := NewParser(&resolver.Mock{}).Parse(ing)
	if err != nil {
		t.Errorf("unexpected error parsing ingress with tls-cert-vault")
	}
	_, ok := i.(string)
	if !ok {
		t.Errorf("expected a string type")
	}

	data := map[string]string{}
	ing.SetAnnotations(data)

	// Test with empty annotations
	i, err = NewParser(&resolver.Mock{}).Parse(ing)
	if err != nil {
		t.Errorf("unexpected error parsing ingress with tls-cert-vault")
	}
	_, ok = i.(string)
	if !ok {
		t.Errorf("expected a string type")
	}

	// Test invalid annotation set
	data[parser.GetAnnotationWithPrefix("tls-cert-vault")] = "INVALID&data"
	ing.SetAnnotations(data)

	_, err = NewParser(&resolver.Mock{}).Parse(ing)
	if err == nil {
		t.Errorf("Expected an error as the value is invalid")
	}

}

func TestParseAnnotations(t *testing.T) {
	ing := buildIngress()

	data := map[string]string{}
	data[parser.GetAnnotationWithPrefix("tls-cert-vault")] = "userland/certificates/mycertificado"
	ing.SetAnnotations(data)

	i, err := NewParser(&resolver.Mock{}).Parse(ing)
	if err != nil {
		t.Errorf("unexpected error parsing ingress with tls-cert-vault")
	}
	_, ok := i.(string)
	if !ok {
		t.Errorf("expected a string type")
	}

}
