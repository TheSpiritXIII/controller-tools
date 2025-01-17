/*
Copyright 2019 The Kubernetes Authors.

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

package crd

import (
	"go/ast"
	"strings"
	"testing"

	"github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
	pkgstest "golang.org/x/tools/go/packages/packagestest"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdmarkers "github.com/TheSpiritXIII/controller-tools/pkg/crd/markers"
	testloader "github.com/TheSpiritXIII/controller-tools/pkg/loader/testutils"
	"github.com/TheSpiritXIII/controller-tools/pkg/markers"
)

func transform(t *testing.T, expr string) *apiext.JSONSchemaProps {
	// this is *very* hacky but I haven’t found a simple way
	// to get an ast.Expr with all the associated metadata required
	// to run typeToSchema upon it:

	moduleName := "github.com/TheSpiritXIII/controller-tools/pkg/crd"
	modules := []pkgstest.Module{
		{
			Name: moduleName,
			Files: map[string]interface{}{
				"test.go": `
				package crd 
				type Test ` + expr,
			},
		},
	}

	pkgs, exported, err := testloader.LoadFakeRoots(pkgstest.Modules, modules, moduleName)
	if exported != nil {
		t.Cleanup(exported.Cleanup)
	}

	if err != nil {
		t.Fatalf("unable to load fake package: %s", err)
	}

	if len(pkgs) != 1 {
		t.Fatal("expected to parse only one package")
	}

	pkg := pkgs[0]
	pkg.NeedTypesInfo()
	failIfErrors(t, pkg.Errors)

	schemaContext := newSchemaContext(pkg, nil, true, false).ForInfo(&markers.TypeInfo{})
	// yick: grab the only type definition
	definedType := pkg.Syntax[0].Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type
	result := typeToSchema(schemaContext, definedType)
	failIfErrors(t, pkg.Errors)
	return result
}

func failIfErrors(t *testing.T, errs []packages.Error) {
	if len(errs) > 0 {
		var msgs []string
		for _, e := range errs {
			msgs = append(msgs, e.Msg)
		}

		t.Fatalf("error loading fake package: %s", strings.Join(msgs, "; "))
	}
}

var arrayOfNumbersSchema *apiext.JSONSchemaProps = &apiext.JSONSchemaProps{
	Type: "array",
	Items: &apiext.JSONSchemaPropsOrArray{
		Schema: &apiext.JSONSchemaProps{
			Type: "number",
		},
	},
}

func Test_Schema_ArrayOfFloat32(t *testing.T) {
	g := gomega.NewWithT(t)

	output := transform(t, "[]float32")
	g.Expect(output).To(gomega.Equal(arrayOfNumbersSchema))
}

func Test_Schema_MapOfStringToArrayOfFloat32(t *testing.T) {
	g := gomega.NewWithT(t)

	output := transform(t, "map[string][]float32")
	g.Expect(output).To(gomega.Equal(&apiext.JSONSchemaProps{
		Type: "object",
		AdditionalProperties: &apiext.JSONSchemaPropsOrBool{
			Allows: true,
			Schema: arrayOfNumbersSchema,
		},
	}))
}

func Test_Schema_ApplyMarkers(t *testing.T) {
	g := gomega.NewWithT(t)

	props := &apiext.JSONSchemaProps{}
	ctx := &schemaContext{}

	var invocations []string

	applyMarkers(ctx, markers.MarkerValues{
		"blah": []interface{}{
			&testPriorityMarker{
				priority: 0, callback: func() {
					invocations = append(invocations, "0")
				},
			},
			&testPriorityMarker{priority: 2, callback: func() {
				invocations = append(invocations, "2")
			}},
			&testPriorityMarker{priority: 11, callback: func() {
				invocations = append(invocations, "11")
			}},
			&defaultPriorityMarker{callback: func() {
				invocations = append(invocations, "default")
			}},
			&testapplyFirstMarker{callback: func() {
				invocations = append(invocations, "applyFirst")
			}},
		}}, props, nil)

	g.Expect(invocations).To(gomega.Equal([]string{"0", "applyFirst", "2", "default", "11"}))
}

type defaultPriorityMarker struct {
	callback func()
}

func (m *defaultPriorityMarker) ApplyToSchema(*apiext.JSONSchemaProps) error {
	m.callback()
	return nil
}

type testPriorityMarker struct {
	priority crdmarkers.ApplyPriority
	callback func()
}

func (m *testPriorityMarker) ApplyPriority() crdmarkers.ApplyPriority {
	return m.priority
}

func (m *testPriorityMarker) ApplyToSchema(*apiext.JSONSchemaProps) error {
	m.callback()
	return nil
}

type testapplyFirstMarker struct {
	callback func()
}

func (m *testapplyFirstMarker) ApplyFirst() {}
func (m *testapplyFirstMarker) ApplyToSchema(*apiext.JSONSchemaProps) error {
	m.callback()
	return nil
}
