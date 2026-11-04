package definition

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
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
		err := ds.LoadDefinitions("testdata/valid.yaml", false)
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

func TestDefinitionStructure_LoadDefinitionsUnresolvableSource(t *testing.T) {
	// "URL" is not a supported ZipFile.SourceType ("url" is), so Icons never
	// gets a CacheFilePath and the icon that depends on it can never resolve.
	data := `Definitions:
  Icons:
    Type: Zip
    ZipFile:
      SourceType: URL
      Url: https://github.com/awslabs/diagram-as-code/archive/icons.zip
  AWS::S3::Bucket:
    Type: Resource
    Icon:
      Source: Icons
      Path: s3.png
`
	filePath := filepath.Join(t.TempDir(), "definitions.yaml")
	if err := os.WriteFile(filePath, []byte(data), 0600); err != nil {
		t.Fatalf("failed to write definition file: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		ds := &DefinitionStructure{}
		done <- ds.LoadDefinitions(filePath, false)
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected an error for a definition whose source never resolves")
		}
		if !strings.Contains(err.Error(), "AWS::S3::Bucket") {
			t.Errorf("expected error to name the unresolved definition, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("LoadDefinitions did not return for a definition whose source never resolves")
	}
}
