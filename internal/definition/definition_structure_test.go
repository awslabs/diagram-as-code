package definition

import (
	"reflect"
	"testing"
)

func TestDefinitionStructure_LoadDefinitions(t *testing.T) {
	url := "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/internal/definition/testdata/aws-icons.zip"
	tc := &DefinitionStructure{
		Definitions: map[string]*Definition{
			"aws_resource": {
				Type: "Resource",
				Icon: &DefinitionIcon{
					Source: "aws_icons",
					Path:   "icons/resource.png",
				},
				Parent: &Definition{
					Type: "Zip",
					ZipFile: DefinitionZipFile{
						SourceType: "url",
						Url:        url,
					},
				},
			},
			"aws_icons": {
				Type: "Zip",
				ZipFile: DefinitionZipFile{
					SourceType: "url",
					Url:        url,
				},
			},
		},
	}

	t.Run("Valid YAML file", func(t *testing.T) {
		ds := &DefinitionStructure{}
		err := ds.LoadDefinitions("testdata/valid.yaml")
		if err != nil {
			t.Errorf("Failed to laod definition file: %v", err)
		}
		ds.Definitions["aws_resource"].CacheFilePath = ""
		ds.Definitions["aws_resource"].Parent.CacheFilePath = ""
		ds.Definitions["aws_icons"].CacheFilePath = ""

		if reflect.DeepEqual(tc.Definitions, ds.Definitions) != true {
			t.Errorf("Definition mismatch\n===load from file===\n%v\b======\n===expects===\n%v\n===", tc.Definitions, ds.Definitions)
		}
	})
}
