// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ctl

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/aws-cloudformation/rain/cft/parse"
	"github.com/awslabs/diagram-as-code/internal/definition"
	"github.com/awslabs/diagram-as-code/internal/types"
)

func TestCloudFormationSharedSubnetResource(t *testing.T) {
	cfn, err := parse.Reader(strings.NewReader(`Resources:
  VPC1:
    Type: AWS::EC2::VPC
  Subnet1:
    Type: AWS::EC2::Subnet
    Properties:
      VpcId: !Ref VPC1
  Subnet2:
    Type: AWS::EC2::Subnet
    Properties:
      VpcId: !Ref VPC1
  NLB:
    Type: AWS::ElasticLoadBalancingV2::LoadBalancer
    Properties:
      Subnets:
        - !Ref Subnet1
        - !Ref Subnet2
`))
	if err != nil {
		t.Fatal(err)
	}
	tmpl := TemplateStruct{Diagram: Diagram{Resources: map[string]Resource{
		"Canvas":   {Type: "AWS::Diagram::Canvas", Children: []string{"AWSCloud"}},
		"AWSCloud": {Type: "AWS::Diagram::Cloud"},
		"NLB":      {Type: "AWS::ElasticLoadBalancingV2::LoadBalancer", Title: "NLB", Preset: "Empty"},
	}}}
	ds := definition.DefinitionStructure{Definitions: map[string]*definition.Definition{
		"AWS::Diagram::Canvas":                      {Type: "Group", CFn: definition.DefinitionCFn{HasChildren: true}},
		"AWS::Diagram::Cloud":                       {Type: "Group", CFn: definition.DefinitionCFn{HasChildren: true}},
		"AWS::EC2::VPC":                             {Type: "Group", CFn: definition.DefinitionCFn{HasChildren: true}},
		"AWS::EC2::Subnet":                          {Type: "Group", CFn: definition.DefinitionCFn{HasChildren: true}},
		"AWS::ElasticLoadBalancingV2::LoadBalancer": {Type: "Resource"},
	}}
	if err := convertTemplate(cfn, &tmpl, ds); err != nil {
		t.Fatal(err)
	}
	ensureSingleParent(&tmpl)
	vpc, _ := tmpl.Resources["VPC1"]
	if !reflect.DeepEqual(vpc.Children, []string{"Subnet1", "Subnet2", "NLB"}) {
		t.Fatalf("VPC children = %v, want both subnets and one load balancer", vpc.Children)
	}
	for _, name := range []string{"Subnet1", "Subnet2"} {
		subnet, _ := tmpl.Resources[name]
		if contains(subnet.Children, "NLB") {
			t.Fatalf("shared load balancer is still a child of %s", name)
		}
	}
	resources := make(map[string]*types.Resource)
	if err := loadResources(&tmpl, ds, resources); err != nil {
		t.Fatal(err)
	}
	associateCFnChildren(&tmpl, ds, resources)
	output := filepath.Join(t.TempDir(), "shared-subnet.png")
	if err := createDiagram(resources, &output, &CreateOptions{OverrideFont: "goregular"}); err != nil {
		t.Fatalf("rendering a load balancer referencing sibling subnets: %v", err)
	}
	for name, resource := range resources {
		if !resource.IsDrawn() {
			t.Errorf("resource %s was not drawn", name)
		}
	}
}

func TestEnsureSingleParent(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input map[string][]string
		want  map[string][]string
	}{
		{
			name:  "empty",
			input: map[string][]string{},
			want:  map[string][]string{},
		},
		{
			name:  "single parent",
			input: map[string][]string{"VPC": {"Subnet"}, "Subnet": {"NLB"}, "NLB": nil},
			want:  map[string][]string{"VPC": {"Subnet"}, "Subnet": {"NLB"}, "NLB": nil},
		},
		{
			name:  "ancestor reference",
			input: map[string][]string{"VPC": {"Subnet", "NLB"}, "Subnet": {"NLB"}, "NLB": nil},
			want:  map[string][]string{"VPC": {"Subnet"}, "Subnet": {"NLB"}, "NLB": nil},
		},
		{
			name: "multiple shared children preserve other children",
			input: map[string][]string{
				"VPC": {"Subnet1", "Subnet2"}, "Subnet1": {"NLB2", "Instance", "NLB1"},
				"Subnet2": {"NLB1", "NLB2"}, "Instance": nil, "NLB1": nil, "NLB2": nil,
			},
			want: map[string][]string{
				"VPC": {"Subnet1", "Subnet2", "NLB1", "NLB2"}, "Subnet1": {"Instance"},
				"Subnet2": {}, "Instance": nil, "NLB1": nil, "NLB2": nil,
			},
		},
		{
			name:  "missing resource",
			input: map[string][]string{"VPC": {"Subnet", "Missing"}, "Subnet": {"NLB"}, "NLB": nil},
			want:  map[string][]string{"VPC": {"Subnet", "Missing"}, "Subnet": {"NLB"}, "NLB": nil},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmpl := TemplateStruct{Diagram: Diagram{Resources: make(map[string]Resource)}}
			for name, children := range tc.input {
				tmpl.Resources[name] = Resource{Children: children}
			}
			ensureSingleParent(&tmpl)
			for name, want := range tc.want {
				resource, ok := tmpl.Resources[name]
				if !ok || !reflect.DeepEqual(resource.Children, want) {
					t.Errorf("%s children = %v, want %v", name, resource.Children, want)
				}
			}
		})
	}
}
