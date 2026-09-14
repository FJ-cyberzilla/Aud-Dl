package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type ComplexityMetrics struct {
	File       string
	Function   string
	Cyclomatic int
	Cognitive  int
}

func analyzeFile(path string) []ComplexityMetrics {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	var metrics []ComplexityMetrics

	ast.Inspect(node, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			m := ComplexityMetrics{
				File:       path,
				Function:   fn.Name.Name,
				Cyclomatic: calculateCyclomatic(fn),
				Cognitive:  calculateCognitive(fn),
			}
			metrics = append(metrics, m)
		}
		return true
	})

	return metrics
}

func calculateCyclomatic(fn *ast.FuncDecl) int {
	complexity := 1
	ast.Inspect(fn, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause:
			complexity++
		}
		return true
	})
	return complexity
}

func calculateCognitive(fn *ast.FuncDecl) int {
	return computeComplexity(fn.Body, 0)
}

func computeComplexity(node ast.Node, depth int) int {
	complexity := 0
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		switch t := n.(type) {
		case *ast.IfStmt:
			complexity += 1 + depth
			complexity += computeComplexity(t.Body, depth+1)
			if t.Else != nil {
				complexity += computeComplexity(t.Else, depth+1)
			}
			return false // Handled recursively
		case *ast.ForStmt:
			complexity += 1 + depth
			complexity += computeComplexity(t.Body, depth+1)
			return false // Handled recursively
		case *ast.RangeStmt:
			complexity += 1 + depth
			complexity += computeComplexity(t.Body, depth+1)
			return false // Handled recursively
		}
		return true // Continue inspecting
	})
	return complexity
}

func main() {
	var allMetrics []ComplexityMetrics
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" && !strings.Contains(path, "vendor/") && !strings.Contains(path, "tools/") {
			metrics := analyzeFile(path)
			allMetrics = append(allMetrics, metrics...)
		}
		return nil
	})

	reportFile, _ := os.Create("complexity_report.txt")
	defer reportFile.Close()

	fmt.Fprintln(reportFile, "=== Complexity Report ===")
	for _, m := range allMetrics {
		fmt.Fprintf(reportFile, "File: %s | Func: %s | Cyclomatic: %d | Cognitive: %d\n", m.File, m.Function, m.Cyclomatic, m.Cognitive)
	}

	fmt.Println("Complexity report generated: complexity_report.txt")
}