package definition

import (
	"testing"
)

func TestDefinition(t *testing.T) {
	t.Run("ValidDefinition", func(t *testing.T) {
		def := Definition{
			Type: "Type1",
			Icon: &DefinitionIcon{
				Source: "IconSource",
				Path:   "IconPath",
			},
			Label: &DefinitionLabel{
				Title: "LabelTitle",
				Color: "LabelColor",
				Font:  "LabelFont",
			},
			Fill: &DefinitionFill{
				Color: "FillColor",
			},
			Border: &DefinitionBorder{
				Color: "BorderColor",
			},
			Directory: DefinitionDirectory{
				Source: "DirectorySource",
				Path:   "DirectoryPath",
			},
			ZipFile: DefinitionZipFile{
				SourceType: "ZipSourceType",
				Source:     "ZipSource",
				Path:       "ZipPath",
				Url:        "ZipUrl",
			},
			CacheFilePath: "CacheFilePath",
		}

		// Add assertions to verify the fields of the Definition struct
		// For example:
		if def.Type != "Type1" {
			t.Errorf("Expected Type to be 'Type1', got %s", def.Type)
		}

		// Add more assertions for other fields
	})

	t.Run("NestedDefinition", func(t *testing.T) {
		parent := &Definition{
			Type: "ParentType",
		}

		child := &Definition{
			Type:   "ChildType",
			Parent: parent,
		}

		// Add assertions to verify the Parent field
		if child.Parent != parent {
			t.Errorf("Expected child.Parent to be %v, got %v", parent, child.Parent)
		}
	})

}
