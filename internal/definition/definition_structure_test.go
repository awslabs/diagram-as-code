package definition

import (
	"reflect"
	"testing"
)

func TestDefinitionStructure_LoadDefinitions(t *testing.T) {
	url := "https://raw.githubusercontent.com/awslabs/diagram-as-code/main/internal/definition/testdata/aws-icons.zip"
	tc := &DefinitionStructure{
		Definitions: map[string]*Definition{
			"Icons": {
				Type: "Zip",
				ZipFile: DefinitionZipFile{
					SourceType: "url",
					Url:        url,
				},
			},
			"IconsDirectory": {
				Type: "Directory",
				Directory: DefinitionDirectory{
					Source: "Icons",
					Path:   "icons/",
				},
				Parent: &Definition{
					Type: "Zip",
					ZipFile: DefinitionZipFile{
						SourceType: "url",
						Url:        url,
					},
				},
			},

			"IconResource": {
				Type: "Resource",
				Icon: &DefinitionIcon{
					Source: "IconsDirectory",
					Path:   "resource.png",
				},
				Parent: &Definition{
					Type: "Directory",
					Directory: DefinitionDirectory{
						Source: "Icons",
						Path:   "icons/",
					},
					Parent: &Definition{
						Type: "Zip",
						ZipFile: DefinitionZipFile{
							SourceType: "url",
							Url:        url,
						},
					},
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
		ds.Definitions["Icons"].CacheFilePath = ""
		ds.Definitions["IconsDirectory"].CacheFilePath = ""
		ds.Definitions["IconsDirectory"].Parent.CacheFilePath = ""
		ds.Definitions["IconResource"].CacheFilePath = ""
		ds.Definitions["IconResource"].Parent.CacheFilePath = ""
		ds.Definitions["IconResource"].Parent.Parent.CacheFilePath = ""

		if reflect.DeepEqual(tc.Definitions, ds.Definitions) != true {
			t.Errorf("Definition mismatch\n===load from file===\n%v\b======\n===expects===\n%v\n===", tc.Definitions, ds.Definitions)
		}
	})
}
