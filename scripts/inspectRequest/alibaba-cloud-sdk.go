package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// parseFile parses a single .go file for request and response structs and returns a slice of Api structs
func parseFile(path string) ([]*Api, error) {
	apis := []*Api{}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var currentApi *Api
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if structType, ok := x.Type.(*ast.StructType); ok {
				name := x.Name.Name
				if strings.HasSuffix(name, "Request") {
					if currentApi != nil {
						apis = append(apis, currentApi)
					}
					currentApi = &Api{
						Name:     strings.TrimSuffix(name, "Request"),
						Requests: extractFields(structType),
					}
				} else if strings.HasSuffix(name, "Response") {
					if currentApi != nil && currentApi.Name == strings.TrimSuffix(name, "Response") {
						currentApi.Responses = extractFields(structType)
					}
				}
			}
		case *ast.FuncDecl:
			if currentApi != nil && x.Name.Name == "Create"+currentApi.Name+"Request" {
				for _, stmt := range x.Body.List {
					switch expr := stmt.(type) {
					case *ast.ExprStmt:
						if callExpr, ok := expr.X.(*ast.CallExpr); ok {
							if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
								if fun.Sel.Name == "InitWithApiInfo" {
									if len(callExpr.Args) >= 3 {
										currentApi.Product = getValue(callExpr.Args[0])
										currentApi.Version = getValue(callExpr.Args[1])
										currentApi.name = getValue(callExpr.Args[2])
									}
								}
							}
						}
					}
				}
				if currentApi != nil {
					apis = append(apis, currentApi)
				}
				currentApi = nil
			}
		}
		return true
	})

	if currentApi != nil {
		apis = append(apis, currentApi)
	}
	return apis, nil
}

// extractFields extracts field information from a struct type
func extractFields(structType *ast.StructType) []Parm {
	fields := []Parm{}
	for _, field := range structType.Fields.List {
		typ := exprToString(field.Type)
		realName := ""
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			realName = tag.Get("name")
		}
		for _, nameIdent := range field.Names {
			fields = append(fields, Parm{
				Name:     nameIdent.Name,
				Type:     typ,
				RealName: realName,
			})
		}
	}
	return fields
}

// exprToString converts an ast.Expr to its string representation
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// getValue extracts the string value from ast.BasicLit
func getValue(expr ast.Expr) string {
	if lit, ok := expr.(*ast.BasicLit); ok {
		return strings.Trim(lit.Value, "\"")
	}
	return ""
}

func manifestAlibabaCloudSdk() []*Api {
	apis := make([]*Api, 0)
	dir := "../../local_vendor/github.com/aliyun/alibaba-cloud-sdk-go/services"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			newApis, err := parseFile(path)
			if err != nil {
				return err
			}
			apis = append(apis, newApis...)
		}
		return nil
	})
	if err != nil {
		log.Println("Error scanning directory:", err)
		return nil
	}

	return apis
}
