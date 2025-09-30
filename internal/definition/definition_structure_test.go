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
		// Clear cache file paths for comparison (these are set dynamically during loading)
		if icons, exists := ds.Definitions["Icons"]; exists {
			icons.CacheFilePath = ""
		} else {
			t.Error("Expected 'Icons' definition not found")
		}
		if iconsDir, exists := ds.Definitions["IconsDirectory"]; exists {
			iconsDir.CacheFilePath = ""
			if iconsDir.Parent != nil {
				iconsDir.Parent.CacheFilePath = ""
			}
		} else {
			t.Error("Expected 'IconsDirectory' definition not found")
		}
		if iconRes, exists := ds.Definitions["IconResource"]; exists {
			iconRes.CacheFilePath = ""
			if iconRes.Parent != nil {
				iconRes.Parent.CacheFilePath = ""
				if iconRes.Parent.Parent != nil {
					iconRes.Parent.Parent.CacheFilePath = ""
				}
			}
		} else {
			t.Error("Expected 'IconResource' definition not found")
		}

		if reflect.DeepEqual(tc.Definitions, ds.Definitions) != true {
			t.Errorf("Definition mismatch\n===load from file===\n%v\b======\n===expects===\n%v\n===", tc.Definitions, ds.Definitions)
		}
	})
}
