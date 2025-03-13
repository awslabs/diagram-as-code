package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

type JsonStructure struct {
	Data struct {
		Items []struct {
			Fields struct {
				Button1DropdownItems string `json:"button1DropdownItems"`
			} `json:"fields"`
		} `json:"items"`
	} `json:"data"`
}

// Update the Slide struct to use the Shape type
type Slide struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/presentationml/2006/main sld"`
	CSld    struct {
		SpTree struct {
			Pictures []struct {
				NvPicPr struct {
					CNvPr struct {
						ID    string `xml:"id,attr"`
						Name  string `xml:"name,attr"`
						Descr string `xml:"descr,attr"`
					} `xml:"http://schemas.openxmlformats.org/presentationml/2006/main cNvPr"`
				} `xml:"nvPicPr"`
				BlipFill struct {
					Blip struct {
						Embed string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships embed,attr"`
					} `xml:"http://schemas.openxmlformats.org/drawingml/2006/main blip"`
				} `xml:"blipFill"`
			} `xml:"http://schemas.openxmlformats.org/presentationml/2006/main pic"`
			AlternativePics []struct {
				NvPicPr struct {
					CNvPr struct {
						ID    string `xml:"id,attr"`
						Name  string `xml:"name,attr"`
						Descr string `xml:"descr,attr"`
					} `xml:"http://schemas.openxmlformats.org/presentationml/2006/main cNvPr"`
				} `xml:"nvPicPr"`
				BlipFill struct {
					Blip struct {
						Embed string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships embed,attr"`
					} `xml:"http://schemas.openxmlformats.org/drawingml/2006/main blip"`
				} `xml:"blipFill"`
			} `xml:"http://schemas.openxmlformats.org/presentationml/2006/main altPic"`
			Groups []Group `xml:"http://schemas.openxmlformats.org/presentationml/2006/main grpSp"`
		} `xml:"spTree"`
	} `xml:"cSld"`
}

// Update the Group struct to use the Shape type
type Group struct {
	NvGrpSpPr struct {
		CNvPr struct {
			ID    string `xml:"id,attr"`
			Name  string `xml:"name,attr"`
			Descr string `xml:"descr,attr"`
		} `xml:"http://schemas.openxmlformats.org/presentationml/2006/main cNvPr"`
	} `xml:"nvGrpSpPr"`
	Sp []struct {
		TxBody struct {
			P struct {
				R struct {
					T string `xml:",chardata"`
				} `xml:"http://schemas.openxmlformats.org/drawingml/2006/main r"`
			} `xml:"http://schemas.openxmlformats.org/drawingml/2006/main p"`
		} `xml:"txBody"`
	} `xml:"http://schemas.openxmlformats.org/presentationml/2006/main sp"`
	Pic struct {
		BlipFill struct {
			Blip struct {
				Embed string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships embed,attr"`
			} `xml:"http://schemas.openxmlformats.org/drawingml/2006/main blip"`
		} `xml:"blipFill"`
	} `xml:"http://schemas.openxmlformats.org/presentationml/2006/main pic"`
	Groups []Group `xml:"http://schemas.openxmlformats.org/presentationml/2006/main grpSp"`
}

func cleanName(name string) string {
	// Implement the same cleaning rules as in the bash script
	name = regexp.MustCompile(` group\.`).ReplaceAllString(name, "")
	name = regexp.MustCompile(` Service icon\.`).ReplaceAllString(name, "")
	// Add other cleaning rules...
	return strings.TrimSpace(name)
}

// Structure for YAML template
type YAMLData struct {
	URL   string
	Icons []IconData
}

type IconData struct {
	Name     string
	ImageRef string
}

func main() {
	url, err := getPPTXUrl()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("PPTx file url: %s\n\n", url)

	// Download and process the zip file
	fmt.Println("Downloading PPTX file...")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Reading zip content...")
	zipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}
	fmt.Printf("Read %d bytes from zip file\n", len(zipBytes))

	// Create a reader for the outer zip content
	fmt.Println("\nCreating zip reader...")
	outerZipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		fmt.Printf("Error creating zip reader: %v\n", err)
		return
	}

	// Find the Light BG PPTX file
	var pptxFile *zip.File
	fmt.Printf("\nSearching for Light BG PPTX file...\n")
	for _, file := range outerZipReader.File {
		fmt.Printf("Found file: %s\n", file.Name)
		if strings.Contains(file.Name, "Light-BG") && strings.HasSuffix(file.Name, ".pptx") {
			pptxFile = file
			fmt.Printf("Found Light BG PPTX file: %s\n", file.Name)
			break
		}
	}

	if pptxFile == nil {
		fmt.Println("Error: Light BG PPTX file not found")
		return
	}

	// Read the PPTX file
	fmt.Println("\nReading PPTX file...")
	pptxRC, err := pptxFile.Open()
	if err != nil {
		fmt.Printf("Error opening PPTX file: %v\n", err)
		return
	}
	defer pptxRC.Close()

	pptxBytes, err := io.ReadAll(pptxRC)
	if err != nil {
		fmt.Printf("Error reading PPTX file: %v\n", err)
		return
	}
	fmt.Printf("Read %d bytes from PPTX file\n", len(pptxBytes))

	// Create a reader for the PPTX content
	fmt.Println("\nCreating PPTX reader...")
	pptxReader, err := zip.NewReader(bytes.NewReader(pptxBytes), int64(len(pptxBytes)))
	if err != nil {
		fmt.Printf("Error creating PPTX reader: %v\n", err)
		return
	}

	// Create map to store image mappings
	imageMappings := make(map[string]string)

	// Process slides
	fmt.Printf("\nProcessing files in PPTX...\n")
	slideCount := 0
	for _, file := range pptxReader.File {
		if strings.HasPrefix(file.Name, "ppt/slides/slide") && strings.HasSuffix(file.Name, ".xml") {
			fmt.Printf("\nProcessing slide file: %s\n", file.Name)
			slideCount++
			processSlide(file, imageMappings, pptxReader, file.Name) // Pass pptxReader here
		}
	}
	fmt.Printf("\nProcessed %d slide files\n", slideCount)
	fmt.Printf("Found %d image mappings\n", len(imageMappings))

	// Generate YAML output
	fmt.Println("\nGenerating YAML file...")
	if err := generateYAML(url, imageMappings); err != nil {
		fmt.Printf("Error generating YAML: %v\n", err)
		return
	}
}

// Main processing function that coordinates the slide processing
func processSlide(file *zip.File, imageMappings map[string]string, pptxReader *zip.Reader, slideFilePath string) {
	content, err := readSlideContent(file)
	if err != nil {
		fmt.Printf("Error reading slide content: %v\n", err)
		return
	}

	slide, err := parseSlideXML(content)
	if err != nil {
		fmt.Printf("Error parsing slide XML: %v\n", err)
		return
	}

	relationships := getSlideRelationships(file, pptxReader)

	processSlideImages(slide, relationships, imageMappings, slideFilePath)
}

// Reads the content of a slide file
func readSlideContent(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening slide: %v", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("error reading slide content: %v", err)
	}
	fmt.Printf("Read %d bytes from slide\n", len(content))

	return content, nil
}

func parseSlideXML(content []byte) (*Slide, error) {
	// Print the first 1000 characters of the XML content for debugging
	fmt.Printf("XML Content preview:\n%s\n", string(content[:min(1000, len(content))]))

	var slide Slide
	err := xml.Unmarshal(content, &slide)
	if err != nil {
		return nil, fmt.Errorf("error parsing XML: %v", err)
	}

	totalPictures := len(slide.CSld.SpTree.Pictures) + len(slide.CSld.SpTree.AlternativePics)
	fmt.Printf("Found %d regular pictures and %d alternative pictures in slide (total: %d)\n",
		len(slide.CSld.SpTree.Pictures),
		len(slide.CSld.SpTree.AlternativePics),
		totalPictures,
	)

	return &slide, nil
}

// Gets relationships for a slide
func getSlideRelationships(file *zip.File, pptxReader *zip.Reader) map[string]string {
	relationships := make(map[string]string)

	slideNum := strings.TrimPrefix(strings.TrimSuffix(filepath.Base(file.Name), ".xml"), "slide")
	relsFile := fmt.Sprintf("ppt/slides/_rels/slide%s.xml.rels", slideNum)

	var found bool
	for _, f := range pptxReader.File {
		if f.Name == relsFile {
			relationships = parseRelationshipsFile(f)
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("WARNING: Could not find relationships file %s\n", relsFile)
	}

	return relationships
}

// Parses the relationships file
func parseRelationshipsFile(file *zip.File) map[string]string {
	relationships := make(map[string]string)

	rc, err := file.Open()
	if err != nil {
		fmt.Printf("Error opening relationships file: %v\n", err)
		return relationships
	}
	defer rc.Close()

	// Define the relationships structure with proper namespace
	type Relationships struct {
		XMLName xml.Name `xml:"http://schemas.openxmlformats.org/package/2006/relationships Relationships"`
		Rels    []struct {
			Id     string `xml:"Id,attr"`
			Type   string `xml:"Type,attr"`
			Target string `xml:"Target,attr"`
		} `xml:"Relationship"`
	}

	var rels Relationships
	if err := xml.NewDecoder(rc).Decode(&rels); err != nil {
		fmt.Printf("Error parsing relationships: %v\n", err)
		return relationships
	}

	// Only map image relationships
	const imageType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	for _, rel := range rels.Rels {
		if rel.Type == imageType {
			relationships[rel.Id] = rel.Target
		}
	}

	fmt.Printf("Found %d image relationships\n", len(relationships))
	return relationships
}

// Update the processSlideImages function to handle shapes separately
func processSlideImages(slide *Slide, relationships map[string]string, imageMappings map[string]string, slideFilePath string) {
	// Process regular pictures
	for _, pic := range slide.CSld.SpTree.Pictures {
		fmt.Printf("\nProcessing picture - Raw Name: %s, Image Ref: %s\n",
			pic.NvPicPr.CNvPr.Descr, pic.BlipFill.Blip.Embed)

		name := cleanIconName(pic.NvPicPr.CNvPr.Descr)
		if shouldSkipImage(name) {
			fmt.Printf("Skipping picture with name: %s\n", name)
			continue
		}

		handleImageMapping(name, pic.BlipFill.Blip.Embed, relationships, imageMappings)
	}

	// Process alternative pictures
	for _, pic := range slide.CSld.SpTree.AlternativePics {
		fmt.Printf("\nProcessing alternative picture - Raw Name: %s, Image Ref: %s\n",
			pic.NvPicPr.CNvPr.Descr, pic.BlipFill.Blip.Embed)

		name := cleanIconName(pic.NvPicPr.CNvPr.Descr)
		if shouldSkipImage(name) {
			fmt.Printf("Skipping picture with name: %s\n", name)
			continue
		}

		handleImageMapping(name, pic.BlipFill.Blip.Embed, relationships, imageMappings)
	}

	// Only process groups for slide25.xml
	if strings.HasSuffix(slideFilePath, "slide25.xml") {
		fmt.Printf("\nProcessing groups for slide25.xml\n")
		processGroups(slide.CSld.SpTree.Groups, relationships, imageMappings)
	}
}

func processGroups(groups []Group, relationships map[string]string, imageMappings map[string]string) {
	for i, group := range groups {
		fmt.Printf("Processing group %d:\n", i+1)
		fmt.Printf("  ID: %s\n", group.NvGrpSpPr.CNvPr.ID)
		fmt.Printf("  Name: %s\n", group.NvGrpSpPr.CNvPr.Name)
		fmt.Printf("  Description: %s\n", group.NvGrpSpPr.CNvPr.Descr)

		if group.Pic.BlipFill.Blip.Embed != "" {
			fmt.Printf("  Image Reference: %s\n", group.Pic.BlipFill.Blip.Embed)
		}

		for j, sp := range group.Sp {
			if text := sp.TxBody.P.R.T; text != "" {
				fmt.Printf("  Text content %d: %s\n", j+1, text)
			}
		}

		// Process the group content...
		if group.NvGrpSpPr.CNvPr.Descr != "" {
			name := cleanIconName(group.NvGrpSpPr.CNvPr.Descr)
			if !shouldSkipImage(name) && group.Pic.BlipFill.Blip.Embed != "" {
				handleImageMapping(name, group.Pic.BlipFill.Blip.Embed, relationships, imageMappings)
			}
		}

		// Process text content and nested groups...
		processGroups(group.Groups, relationships, imageMappings)
	}
}

// Determines if an image should be skipped
func shouldSkipImage(name string) bool {
	return name == "" ||
		name == "Example of an architecture diagram" ||
		name == "Graphic icon"
}

// Handles the mapping of images, including duplicate handling
func handleImageMapping(name, embedId string, relationships map[string]string, imageMappings map[string]string) {
	imagePath, exists := relationships[embedId]
	if !exists {
		fmt.Printf("WARNING: Could not find image path for reference ID %s (name: %s)\n", embedId, name)
		fmt.Printf("Available relationship IDs: %v\n", getKeys(relationships))
		return
	}

	// Clean up the image path by removing "../media/" prefix and keeping only the filename
	cleanPath := filepath.Base(imagePath)

	if existingImage, exists := imageMappings[name]; exists && existingImage != cleanPath {
		handleDuplicateName(name, cleanPath, imageMappings)
	} else {
		fmt.Printf("Adding mapping: %s -> %s\n", name, cleanPath)
		imageMappings[name] = cleanPath
	}
}

// Helper function to get map keys for debugging
func getKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k) // Don't add "rId" prefix since it's already in the key
	}
	return keys
}

// Handles duplicate image names by adding a number suffix
func handleDuplicateName(name, embedId string, imageMappings map[string]string) {
	fmt.Printf("Found duplicate name: %s\n", name)
	for i := 2; i <= 9; i++ {
		newName := fmt.Sprintf("%s(%d)", name, i)
		if _, exists := imageMappings[newName]; !exists {
			fmt.Printf("Using new name: %s\n", newName)
			imageMappings[newName] = embedId
			break
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Add this struct to handle Resource type definitions
type ResourceDefinition struct {
	Type        string
	Name        string
	Label       string
	ImagePath   string
	HasChildren bool
}

// Add this function to read and parse the mappings file
func readMappingsFile(filename string) ([]ResourceDefinition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening mappings file: %v", err)
	}
	defer file.Close()

	var resources []ResourceDefinition
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Split line by comma
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}

		resourceType := strings.TrimSpace(parts[0])
		name := strings.TrimSpace(parts[1])

		if name == "" {
			continue
		}

		// Clean up label by removing (number) suffix
		label := regexp.MustCompile(`\([0-9]\)`).ReplaceAllString(name, "")

		// Check if resource type should have children
		hasChildren := false
		switch resourceType {
		case "AWS::ECS::Cluster", "AWS::EKS::Cluster", "AWS::CodePipeline::Pipeline":
			hasChildren = true
		}

		resources = append(resources, ResourceDefinition{
			Type:        resourceType,
			Name:        name,
			Label:       label,
			HasChildren: hasChildren,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading mappings file: %v", err)
	}

	return resources, nil
}

// Modify the generateYAML function to handle both Preset and Resource types
func generateYAML(url string, imageMappings map[string]string) error {
	fmt.Println("Generating YAML file...")

	// Create output directory if it doesn't exist
	err := os.MkdirAll("../definitions", 0755)
	if err != nil {
		return fmt.Errorf("failed to create definitions directory: %v", err)
	}

	// Create output file
	f, err := os.Create("../definitions/definition-for-aws-icons-light.yaml")
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer f.Close()

	// Write YAML header
	fmt.Fprintln(f, "# This file automatically generated by tools/make-definition-from-pptx. If you find a bug in this file, please report the issue.\n")
	fmt.Fprintln(f, "Definitions:")

	// Write Main and related sections
	writeInitialDefinitions(f, url)

	// Get all keys and sort them
	mappingsKeys := make([]string, 0, len(imageMappings))
	for k := range imageMappings {
		mappingsKeys = append(mappingsKeys, k)
	}
	sort.Strings(mappingsKeys)

	// Write Preset type definitions (with quotes)
	for _, name := range mappingsKeys {
		imagePath := imageMappings[name]
		// Remove any (number) from the display title but keep the original name as the key
		displayName := regexp.MustCompile(`\([0-9]\)`).ReplaceAllString(name, "")
		displayName = strings.TrimSpace(displayName) // Remove any extra spaces

		fmt.Fprintf(f, "  %q:\n", name)
		fmt.Fprintln(f, "    Type: Preset")
		fmt.Fprintln(f, "    Icon:")
		fmt.Fprintln(f, "      Source: ArchitectureIconsPptxMedia")
		fmt.Fprintf(f, "      Path: %q\n", imagePath)
		fmt.Fprintln(f, "    Label:")
		fmt.Fprintf(f, "      Title: %q\n", displayName) // Use the clean display name here
		fmt.Fprintln(f, "      Color: \"rgba(0, 0, 0, 255)\"\n")
	}

	// Read and process mappings file
	resources, err := readMappingsFile("make-definition-from-pptx-mappings")
	if err != nil {
		fmt.Printf("Warning: Could not process mappings file: %v\n", err)
		return nil
	}

	// Write Resource type definitions (without quotes for the key)
	for _, res := range resources {
		imagePath, exists := imageMappings[res.Name]
		if !exists {
			fmt.Printf("Warning: No image found for resource %s\n", res.Name)
			continue
		}

		fmt.Fprintf(f, "  %s:\n", res.Type) // No quotes for Resource type key
		fmt.Fprintln(f, "    Type: Resource")
		fmt.Fprintln(f, "    Icon:")
		fmt.Fprintln(f, "      Source: ArchitectureIconsPptxMedia")
		fmt.Fprintf(f, "      Path: %q\n", imagePath)
		fmt.Fprintln(f, "    Label:")
		fmt.Fprintf(f, "      Title: %q\n", res.Label)
		fmt.Fprintln(f, "      Color: \"rgba(0, 0, 0, 255)\"")
		fmt.Fprintf(f, "    CFn:\n      HasChildren: %v\n\n", res.HasChildren)
	}

	fmt.Printf("Successfully generated YAML file with %d icons and %d resources\n",
		len(imageMappings), len(resources))
	return nil
}

func writeInitialDefinitions(f *os.File, url string) {
	initialDefs := `  Main:
    Type: Zip
    ZipFile:
      SourceType: url
      Url: "%s"

  ArchitectureIconsPptx:
    Type: Zip
    ZipFile:
      SourceType: file
      Source: Main
      Path: "AWS-Architecture-Icons-Deck_For-Light-BG_02072025.pptx"

  ArchitectureIconsPptxMedia:
    Type: Directory
    Directory:
      Source: ArchitectureIconsPptx
      Path: "ppt/media/"

  "AWS::Diagram::Canvas":
    Type: Group
    Border:
      Color: "rgba(0, 0, 0, 0)"
    Fill:
      Color: "rgba(255, 255, 255, 255)"
    CFn:
      HasChildren: true

  AWS::Diagram::Cloud:
    Type: Group
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image10.png"
    Border:
      Color: "rgba(0, 0, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Label:
      Title: "AWS Cloud"
      Color: "rgba(0, 0, 0, 255)"
    CFn:
      HasChildren: true

  AWSCloudNoLogo:
    Type: Preset
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image43.png"
    Border:
      Color: "rgba(0, 0, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Label:
      Title: "AWS Cloud"
      Color: "rgba(0, 0, 0, 255)"

  AWS::Region:
    Type: Group
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image45.png"
    Border:
      Type: "dashed"
      Color: "rgba(0, 164, 166, 255)"
    Label:
      Title: "Region"
      Color: "rgba(0, 0, 0, 255)"
    CFn:
      HasChildren: true

  AWS::EC2::AvailabilityZone:
    Type: Group
    Border:
      Type: "dashed"
      Color: "rgba(0, 164, 166, 255)"
    HeaderAlign: center
    Label:
      Title: "Availability Zone"
      Color: "rgba(0, 0, 0, 255)"
    CFn:
      HasChildren: true

  AWS::AutoScaling::AutoScalingGroup:
    Type: Group
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image24.png"
    Border:
      Type: "dashed"
      Color: "rgba(237, 113, 0, 255)"
    HeaderAlign: center
    Label:
      Title: "Auto Scaling Group"
      Color: "rgba(0, 0, 0, 255)"
    CFn:
      HasChildren: true

  AWS::EC2::VPC:
    Type: Group
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image18.png"
    Label:
      Title: "VPC"
      Color: "rgba(0, 0, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Border:
      Color: "rgba(105, 59, 197, 255)"
    CFn:
      HasChildren: true

  AWS::EC2::Subnet:
    Type: Group
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image29.png"
    Label:
      Title: "Subnet"
      Color: "rgba(0, 0, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Border:
      Color: "rgba(122, 161, 22, 255)"
    CFn:
      HasChildren: true

  PrivateSubnet:
    Type: Preset
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image47.png"
    Label:
      Title: "Private Subnet"
      Color: "rgba(0, 0, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Border:
      Color: "rgba(0, 164, 166, 255)"

  PublicSubnet:
    Type: Preset
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image29.png"
    Label:
      Title: "Public Subnet"
      Color: "rgba(0, 0, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Border:
      Color: "rgba(122, 161, 22, 255)"

  AWS::Diagram::DataCenter:
    Type: Group
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image51.png"
    Label:
      Title: "Corporate data center"
      Color: "rgba(29, 137, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Border:
      Color: "rgba(125, 137, 152, 255)"
    CFn:
      HasChildren: true

  AWS::EC2::SpotFleet:
    Type: Group
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image55.png"
    Label:
      Title: "Spot fleet"
      Color: "rgba(29, 137, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Border:
      Color: "rgba(237, 113, 0, 255)"
    CFn:
      HasChildren: true

  AWS::Diagram::Account:
    Type: Group
    Icon:
      Source: ArchitectureIconsPptxMedia
      Path: "image57.png"
    Label:
      Title: "AWS account"
      Color: "rgba(29, 137, 0, 255)"
    Fill:
      Color: "rgba(0, 0, 0, 0)"
    Border:
      Color: "rgba(231, 21, 123, 255)"
    CFn:
      HasChildren: true

  "VPC":
    Type: Preset
    Icon: 
      Source: ArchitectureIconsPptxMedia
      Path: "image18.png"
    Label:
      Title: "VPC"
      Color: "rgba(0, 0, 0, 255)"
`
	fmt.Fprintf(f, initialDefs, url)
}

func cleanIconName(name string) string {
	// Define all the patterns to remove, matching the sed commands
	patterns := []string{
		` group\.`,
		` group inside VPC cloud group and in the AZ`,
		` Service icon\.`,
		` service icon.*\.`,
		` group with AWS logo.`,
		` service\.`,
		`\nservice\.`,
		` group icon\.`,
		` instance icon for the Database category\.`,
		` resource icon for.*`,
		` instance icon for.*`,
		` storage class icon for.*`,
		` standard category icon\.`,
		`A representation of a.*`,
		`\.$`,
	}

	// Apply each pattern
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		name = re.ReplaceAllString(name, "")
	}

	if name == "AWS Cloud group with cloud" {
		name = "AWS Cloud with cloud"
	}

	// Trim any remaining whitespace
	name = strings.TrimSpace(name)

	return name
}

func processZipFile(url string) (map[string]string, error) {
	fmt.Println("Downloading and processing ZIP file...")

	// Download the zip file
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download zip: %v", err)
	}
	defer resp.Body.Close()

	// Read the zip content
	zipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read zip content: %v", err)
	}

	// Create a reader for the zip content
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %v", err)
	}

	// Process the zip contents
	iconMappings := make(map[string]string)
	for _, file := range zipReader.File {
		// Look for PNG files in the appropriate directory
		if strings.HasSuffix(file.Name, ".png") && !strings.Contains(file.Name, "__MACOSX") {
			// Extract icon name from path
			name := filepath.Base(file.Name)
			name = strings.TrimSuffix(name, ".png")

			// Store mapping
			iconMappings[name] = file.Name
		}
	}

	if len(iconMappings) == 0 {
		return nil, fmt.Errorf("no icons found in zip file")
	}

	fmt.Printf("Found %d icons in zip file\n", len(iconMappings))
	return iconMappings, nil
}

func getPPTXUrl() (string, error) {
	fmt.Println("Fetching AWS architecture icons page...")
	resp, err := http.Get("https://aws.amazon.com/architecture/icons/")
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Reading response body...")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println("Parsing HTML...")
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Find all script tags with type="application/json"
	var jsonScripts []string
	var findJsonScripts func(*html.Node)
	findJsonScripts = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "type" && attr.Val == "application/json" {
					if n.FirstChild != nil {
						jsonScripts = append(jsonScripts, n.FirstChild.Data)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findJsonScripts(c)
		}
	}
	findJsonScripts(doc)

	fmt.Printf("Found %d script tags with type='application/json'\n", len(jsonScripts))

	for i, jsonContent := range jsonScripts {
		fmt.Printf("\nTrying to parse JSON script %d...\n", i+1)

		var jsonData JsonStructure
		if err := json.Unmarshal([]byte(jsonContent), &jsonData); err != nil {
			fmt.Printf("Script %d: JSON parsing failed: %v\n", i+1, err)
			continue
		}

		// Check if this JSON has the structure we're looking for
		if len(jsonData.Data.Items) == 0 {
			fmt.Printf("Script %d: No items found in data\n", i+1)
			continue
		}

		fmt.Printf("Script %d: Successfully parsed JSON with %d items\n", i+1, len(jsonData.Data.Items))

		// Parse the HTML in button1DropdownItems
		for _, item := range jsonData.Data.Items {
			htmlContent := item.Fields.Button1DropdownItems
			doc, err := html.Parse(strings.NewReader(htmlContent))
			if err != nil {
				fmt.Printf("Failed to parse HTML content: %v\n", err)
				continue
			}

			// Find the link with "Microsoft PPTx toolkits" text
			var findLink func(*html.Node) string
			findLink = func(n *html.Node) string {
				if n.Type == html.ElementNode && n.Data == "a" {
					// Check if this link contains the text we're looking for
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.TextNode && strings.TrimSpace(c.Data) == "Microsoft PPTx toolkits" {
							// Get the href attribute
							for _, attr := range n.Attr {
								if attr.Key == "href" {
									return attr.Val
								}
							}
						}
					}
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if result := findLink(c); result != "" {
						return result
					}
				}
				return ""
			}

			if url := findLink(doc); url != "" {
				fmt.Println("Found Microsoft PPTx toolkits URL!")
				return url, nil
			}
		}

		fmt.Printf("Script %d: No Microsoft PPTx toolkits URL found\n", i+1)
	}

	return "", fmt.Errorf("PPTx toolkit URL not found in any script tag")
}
