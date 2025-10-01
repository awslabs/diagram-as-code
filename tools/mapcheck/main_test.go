package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func TestMapCheckAnalyzer(t *testing.T) {
	testdata := `package testdata

func testMapAccess() {
	m := make(map[string]int)
	
	// Should be flagged - unsafe
	value1 := m["key1"]
	_ = value1
	
	// Should NOT be flagged - safe with ok
	value2, ok := m["key2"]
	if ok {
		_ = value2
	}
	
	// Should NOT be flagged - safe with _
	value3, _ := m["key3"]
	_ = value3
	
	// Should NOT be flagged - safe with _, ok (existence check only)
	_, ok2 := m["key4"]
	if ok2 {
		// key exists
	}
	
	// Should NOT be flagged - map assignment
	m["key6"] = 100
	
	// Should be flagged - unsafe
	if m["key5"] > 0 {
		// do something
	}
}`

	// Create a temporary file for testing
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", testdata, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	// Create type info
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	conf := &types.Config{}
	pkg, err := conf.Check("testdata", fset, []*ast.File{file}, info)
	if err != nil {
		t.Fatal(err)
	}

	// Create analysis pass
	var issues []analysis.Diagnostic
	pass := &analysis.Pass{
		Analyzer:  MapCheckAnalyzer,
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		ResultOf:  make(map[*analysis.Analyzer]interface{}),
		Report: func(d analysis.Diagnostic) {
			issues = append(issues, d)
			t.Logf("Found issue at %s: %s", fset.Position(d.Pos), d.Message)
		},
	}

	// Run inspect analyzer first to set up ResultOf
	inspectResult, err := inspect.Analyzer.Run(pass)
	if err != nil {
		t.Fatal(err)
	}
	pass.ResultOf[inspect.Analyzer] = inspectResult

	// Run the analyzer
	_, err = MapCheckAnalyzer.Run(pass)
	if err != nil {
		t.Fatal(err)
	}

	// We expect 2 issues: value1 := m["key1"] and if m["key5"] > 0
	expectedIssues := 2
	if len(issues) != expectedIssues {
		t.Errorf("Expected %d issues, but found %d", expectedIssues, len(issues))
	}

	// Check specific issue locations and messages
	if len(issues) >= 1 {
		pos1 := fset.Position(issues[0].Pos)
		if pos1.Line != 7 || pos1.Column != 12 {
			t.Errorf("First issue expected at line 7, column 12, but found at line %d, column %d", pos1.Line, pos1.Column)
		}
		expectedMsg := "map access without comma-ok idiom: use 'value, ok := map[key]', 'value, _ := map[key]', or '_, ok := map[key]' instead"
		if issues[0].Message != expectedMsg {
			t.Errorf("First issue message mismatch.\nExpected: %s\nActual: %s", expectedMsg, issues[0].Message)
		}
	}

	if len(issues) >= 2 {
		pos2 := fset.Position(issues[1].Pos)
		if pos2.Line != 30 || pos2.Column != 5 {
			t.Errorf("Second issue expected at line 30, column 5, but found at line %d, column %d", pos2.Line, pos2.Column)
		}
	}

	t.Logf("MapCheck analyzer found %d issues with correct locations and messages", len(issues))
}
