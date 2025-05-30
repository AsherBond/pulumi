// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//nolint:lll
package dotnet

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pulumi/pulumi/pkg/v3/codegen"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

// DocLanguageHelper is the DotNet-specific implementation of the DocLanguageHelper.
type DocLanguageHelper struct{}

var _ codegen.DocLanguageHelper = DocLanguageHelper{}

// GetDocLinkForPulumiType returns the .Net API doc link for a Pulumi type.
func (d DocLanguageHelper) GetDocLinkForPulumiType(pkg *schema.Package, typeName string) string {
	var filename string
	switch typeName {
	// We use docfx to generate the .NET language docs. docfx adds a suffix
	// to generic classes. The suffix depends on the number of type args the class accepts,
	// which in the case of the Pulumi.Input class is 1.
	case "Pulumi.Input":
		filename = "Pulumi.Input-1"
	default:
		filename = typeName
	}
	return fmt.Sprintf("/docs/reference/pkg/dotnet/Pulumi/%s.html", filename)
}

// GetDocLinkForResourceType returns the .NET API doc URL for a type belonging to a resource provider.
func (d DocLanguageHelper) GetDocLinkForResourceType(pkg *schema.Package, _, typeName string) string {
	typeName = strings.ReplaceAll(typeName, "?", "")
	var packageNamespace string
	if pkg == nil {
		packageNamespace = ""
	} else if pkg.Name != "" {
		info, _ := pkg.Language["csharp"].(CSharpPackageInfo)
		packageNamespace = "." + namespaceName(info.Namespaces, pkg.Name)
	}
	return fmt.Sprintf("/docs/reference/pkg/dotnet/Pulumi%s/%s.html", packageNamespace, typeName)
}

// GetDocLinkForResourceInputOrOutputType returns the doc link for an input or output type of a Resource.
func (d DocLanguageHelper) GetDocLinkForResourceInputOrOutputType(pkg *schema.Package, moduleName, typeName string, input bool) string {
	return d.GetDocLinkForResourceType(pkg, moduleName, typeName)
}

// GetDocLinkForFunctionInputOrOutputType returns the doc link for an input or output type of a Function.
func (d DocLanguageHelper) GetDocLinkForFunctionInputOrOutputType(pkg *schema.Package, moduleName, typeName string, input bool) string {
	return d.GetDocLinkForResourceType(pkg, moduleName, typeName)
}

// GetLanguageTypeString returns the DotNet-specific type given a Pulumi schema type.
func (d DocLanguageHelper) GetTypeName(pkg schema.PackageReference, t schema.Type, input bool, relativeToModule string) string {
	var info CSharpPackageInfo
	if a, err := pkg.Language("csharp"); err == nil {
		info, _ = a.(CSharpPackageInfo)
	}
	mod := &modContext{
		pkg:           pkg,
		mod:           relativeToModule,
		typeDetails:   map[*schema.ObjectType]*typeDetails{},
		namespaces:    info.Namespaces,
		rootNamespace: info.GetRootNamespace(),
	}
	qualifier := "Inputs"
	if !input {
		qualifier = "Outputs"
	}
	return mod.typeString(t, qualifier, input, false /*state*/, true /*requireInitializers*/)
}

func (d DocLanguageHelper) GetResourceName(r *schema.Resource) string {
	return resourceName(r)
}

func (d DocLanguageHelper) GetFunctionName(f *schema.Function) string {
	return tokenToFunctionName(f.Token)
}

// GetResourceFunctionResultName returns the name of the result type when a function is used to lookup
// an existing resource.
func (d DocLanguageHelper) GetResourceFunctionResultName(modName string, f *schema.Function) string {
	funcName := d.GetFunctionName(f)
	return funcName + "Result"
}

func (d DocLanguageHelper) GetMethodName(m *schema.Method) string {
	return Title(m.Name)
}

func (d DocLanguageHelper) GetModuleName(pkg schema.PackageReference, module string) string {
	a, _ := pkg.Language("csharp")
	info, _ := a.(CSharpPackageInfo)
	return namespaceName(info.Namespaces, module)
}

func (d DocLanguageHelper) GetMethodResultName(pkg schema.PackageReference, modName string, r *schema.Resource,
	m *schema.Method,
) string {
	a, _ := pkg.Language("csharp")
	info, _ := a.(CSharpPackageInfo)
	var returnType *schema.ObjectType
	if m.Function.ReturnType != nil {
		if objectType, ok := m.Function.ReturnType.(*schema.ObjectType); ok {
			returnType = objectType
		} else {
			mod := &modContext{
				pkg:           pkg,
				mod:           modName,
				typeDetails:   map[*schema.ObjectType]*typeDetails{},
				namespaces:    info.Namespaces,
				rootNamespace: info.GetRootNamespace(),
			}
			return mod.typeString(m.Function.ReturnType, "", false, false, false)
		}
	}

	if info.LiftSingleValueMethodReturns && returnType != nil && len(returnType.Properties) == 1 {
		mod := &modContext{
			pkg:           pkg,
			mod:           modName,
			typeDetails:   map[*schema.ObjectType]*typeDetails{},
			namespaces:    info.Namespaces,
			rootNamespace: info.GetRootNamespace(),
		}
		return mod.typeString(returnType.Properties[0].Type, "", false, false, false)
	}
	return fmt.Sprintf("%s.%sResult", resourceName(r), d.GetMethodName(m))
}

// GetPropertyName uses the property's csharp-specific language info, if available, to generate
// the property name. Otherwise, returns the PascalCase as the default.
func (d DocLanguageHelper) GetPropertyName(p *schema.Property) (string, error) {
	propLangName := strings.Title(p.Name)

	if raw, ok := p.Language["csharp"].(json.RawMessage); ok {
		val, err := Importer.ImportPropertySpec(raw)
		if err != nil {
			return "", err
		}
		p.Language["csharp"] = val
	}

	names := map[*schema.Property]string{}
	properties := []*schema.Property{p}
	computePropertyNames(properties, names)
	if name, ok := names[p]; ok {
		return name, nil
	}
	return propLangName, nil
}

// GetEnumName returns the enum name specific to C#.
func (d DocLanguageHelper) GetEnumName(e *schema.Enum, typeName string) (string, error) {
	name := fmt.Sprintf("%v", e.Value)
	if e.Name != "" {
		name = e.Name
	}
	return makeSafeEnumName(name, typeName)
}
