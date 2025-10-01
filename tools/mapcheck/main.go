package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/singlechecker"
	"golang.org/x/tools/go/ast/inspector"
)

var MapCheckAnalyzer = &analysis.Analyzer{
	Name:     "mapcheck",
	Doc:      "enforce comma-ok idiom for map access",
	Run:      checkMapAccess,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func checkMapAccess(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Track all comma-ok assignments
	commaOkAssignments := make(map[ast.Node]bool)

	// Track assignment LHS (should not be flagged)
	assignmentLHS := make(map[ast.Node]bool)

	// First pass: identify comma-ok assignments and assignment LHS
	nodeFilter1 := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	inspect.Preorder(nodeFilter1, func(n ast.Node) {
		assignStmt := n.(*ast.AssignStmt)

		// Track all LHS IndexExpr in assignments (these should not be flagged)
		for _, lhs := range assignStmt.Lhs {
			if indexExpr, ok := lhs.(*ast.IndexExpr); ok {
				assignmentLHS[indexExpr] = true
			}
		}

		// Check if this is a comma-ok assignment (2 LHS, 1 RHS)
		if len(assignStmt.Lhs) == 2 && len(assignStmt.Rhs) == 1 {
			if indexExpr, ok := assignStmt.Rhs[0].(*ast.IndexExpr); ok {
				if getMapType(pass.TypesInfo, indexExpr.X) != nil {
					commaOkAssignments[indexExpr] = true
				}
			}
		}
	})

	// Second pass: check all map accesses
	nodeFilter2 := []ast.Node{
		(*ast.IndexExpr)(nil),
	}

	inspect.Preorder(nodeFilter2, func(n ast.Node) {
		indexExpr := n.(*ast.IndexExpr)

		// Check if the indexed expression is a map type
		if mapType := getMapType(pass.TypesInfo, indexExpr.X); mapType != nil {
			// Check if this is NOT a comma-ok assignment AND NOT an assignment LHS
			if !commaOkAssignments[indexExpr] && !assignmentLHS[indexExpr] {
				pass.Reportf(indexExpr.Pos(), "map access without comma-ok idiom: use 'value, ok := map[key]', 'value, _ := map[key]', or '_, ok := map[key]' instead")
			}
		}
	})

	return nil, nil
}

func getMapType(info *types.Info, expr ast.Expr) *types.Map {
	if t := info.TypeOf(expr); t != nil {
		if mapType, ok := t.Underlying().(*types.Map); ok {
			return mapType
		}
	}
	return nil
}

func main() {
	singlechecker.Main(MapCheckAnalyzer)
}
