package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
)

func ModifyCase(filename string, api *Api, apsaraApi *Api) bool {
	apiname := api.Name
	requestModify := make(map[string]string)
	for _, beforeRequest := range api.Requests {
		for _, afterRequest := range apsaraApi.Requests {
			if strings.EqualFold(beforeRequest.Name, afterRequest.Name) {
				requestModify[beforeRequest.Name] = afterRequest.Name
			}
		}
	}
	if len(requestModify) == 0 {
		return false
	}

	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return false
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BlockStmt:
			lock := false
			//先检查Apiname
			for _, stmt := range x.List {
				if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
					for _, rhs := range assignStmt.Rhs {
						call, ok := rhs.(*ast.CallExpr)
						var pkName, fnName string
						if ok {
							fnName = getFuncName(call)
							pkName = getPackageName(call)
						}

						//部分不规范的product
						if pkName == "cr_ee" {
							pkName = "cr"
						}
						if pkName == "r_kvstore" {
							pkName = "R-kvstore"
						}
						if pkName == "slsPop" {
							pkName = "Sls"
						}

						if fnName == "NewCommonRequest" {
							for _, stmt := range x.List {
								if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
									if lhs, ok := assignStmt.Lhs[0].(*ast.SelectorExpr); ok {
										if lhsx, ok := lhs.X.(*ast.Ident); ok && lhsx.Name == "request" {
											if rhs, ok := assignStmt.Rhs[0].(*ast.BasicLit); ok {
												value := rhs.Value[1 : len(rhs.Value)-1]
												if lhs.Sel.Name == "ApiName" {
													if value == apiname {
														lock = true
													}
												}

											}
										}
									}
								}
							}
						}
					}
				}
			}
			//检查参数并修复
			for _, stmt := range x.List {
				if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
					if lhs, ok := assignStmt.Lhs[0].(*ast.SelectorExpr); ok {
						if lhs.Sel.Name == "QueryParams" {
							rhs, ok := assignStmt.Rhs[0].(*ast.CompositeLit)
							if !ok {
								continue
							}
							elts := rhs.Elts
							for i := 0; i < len(elts); i++ {
								e := elts[i]
								kv := e.(*ast.KeyValueExpr)
								k := kv.Key.(*ast.BasicLit)
								request := k.Value[1 : len(k.Value)-1]
								if request == "AccessKeyId" || request == "RegionId" || request == "AccessKeySecret" || request == "Product" || request == "Action" || request == "Version" || request == "ProductName" || request == "Department" || request == "ResourceGroup" {
									continue
								}
								if request == "X-acs-body" || request == "SignatureVersion" || request == "Language" || request == "SignatureMethod" || request == "AccountInfo" {
									continue
								}
								for beforeRequest, afterRequest := range requestModify {
									if request == beforeRequest && lock == true {
										modify := "\"" + afterRequest + "\""
										k.Value = modify
										break
									}
								}
							}
						} else if lhsx, ok := lhs.X.(*ast.Ident); ok && lhsx.Name == "request" {
							if rhs, ok := assignStmt.Rhs[0].(*ast.BasicLit); ok {
								value := rhs.Value[1 : len(rhs.Value)-1]
								if lhs.Sel.Name == "ApiName" {
									if value == apiname {
										lock = true
									}
								}

							}
						}
					}
				}
			}
		}

		return true
	})

	file, _ := os.Create(filename)
	var output strings.Builder
	if err = printer.Fprint(&output, fs, node); err != nil {
		return false
	}
	file.WriteString(output.String())
	log.Println("参数不规范修复    ", "文件名:", filename, "    sdk:", api.Name)
	log.Println(requestModify)
	log.Println("-----------------------------------")
	return true
}
