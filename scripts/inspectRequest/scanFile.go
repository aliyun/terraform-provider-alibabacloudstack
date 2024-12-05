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

func GetApisFromScanFiles() []*Api {

	dirname := "../../alibabacloudstack" // your directory containing Go source files
	apis := make([]*Api, 0)

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "test.go") { // process only Go files
			if strings.HasPrefix(name, "resource") || strings.HasPrefix(name, "data") || strings.HasPrefix(name, "service") {
				scanFile(path, &apis)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return nil
	}

	//Printing the results
	// for _, api := range apis {
	// 	fmt.Printf("%s\n", api.FileName)
	// 	fmt.Printf("%s\n", api.Name)
	// 	fmt.Printf("%s\n", api.Product)
	// 	fmt.Printf("Requests:\n")
	// 	for _, req := range api.Requests {
	// 		fmt.Printf("\t%s\n", req.Name)
	// 	}
	// 	fmt.Printf("Responses:\n")
	// 	for _, resp := range api.Responses {
	// 		fmt.Printf("\t%s\n", resp.Name)
	// 	}
	// }
	return apis
}

func scanFile(filename string, apis *[]*Api) {

	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 扫描高阶sdk
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return true
		}
		var api *Api
		requestParams := make(map[string]bool)

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			switch block := n.(type) {
			case *ast.BlockStmt:
				for _, stmt := range block.List {
					if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
						var fnName, pkName string
						leftExpr := assignStmt.Lhs[0]
						rightExpr := assignStmt.Rhs[0]
						call, ok := rightExpr.(*ast.CallExpr)
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
						if strings.HasPrefix(fnName, "Create") && strings.HasSuffix(fnName, "Request") {
							if api != nil {
								for param := range requestParams {
									api.Requests = append(api.Requests, Parm{Name: param})
								}
								if len(api.Requests) > 0 {
									if api.name != "" && api.Product != "" {
										*apis = append(*apis, api)
									}
								}
								api = nil
								requestParams = make(map[string]bool)
							}
							api = new(Api)
							api.Type = "high"
							api.Name = fnName[6 : len(fnName)-7] // Extract the core part of the function name.
							api.name = api.Name
							api.Product = pkName
							api.FileName = filename

							continue
						} else {
							selExpr, ok := leftExpr.(*ast.SelectorExpr)
							if !ok {
								continue
							}
							xIdent, ok := selExpr.X.(*ast.Ident)
							if ok && xIdent.Name == "request" {
								if selExpr.Sel.Name != "IaasProvider" && selExpr.Sel.Name != "Headers" && selExpr.Sel.Name != "QueryParams" && selExpr.Sel.Name != "RegionId" && selExpr.Sel.Name != "Schema" {
									requestParams[selExpr.Sel.Name] = true
								}
							}
						}

					}
				}
			}

			return true
		})

		if api != nil {
			for param := range requestParams {
				api.Requests = append(api.Requests, Parm{Name: param})
			}
			if len(api.Requests) > 0 {
				if api.name != "" && api.Product != "" {
					*apis = append(*apis, api)
				}
			}
			api = nil
			requestParams = make(map[string]bool)
		}
		return true
	})

	//扫描common sdk
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BlockStmt:
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
							var api = new(Api)
							api.Type = "common"
							api.FileName = filename
							requestParams := make(map[string]bool)
							for _, stmt := range x.List {
								if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
									if lhs, ok := assignStmt.Lhs[0].(*ast.SelectorExpr); ok {
										if lhs.Sel.Name == "QueryParams" {
											rhs, ok := assignStmt.Rhs[0].(*ast.CompositeLit)
											if !ok {
												continue
											}
											elts := rhs.Elts
											for _, e := range elts {
												kv := e.(*ast.KeyValueExpr)
												k := kv.Key.(*ast.BasicLit)
												request := k.Value[1 : len(k.Value)-1]
												if request == "AccessKeyId" || request == "RegionId" || request == "AccessKeySecret" || request == "Product" || request == "Action" || request == "Version" || request == "ProductName" || request == "Department" || request == "ResourceGroup" {
													continue
												}
												if request == "X-acs-body" || request == "SignatureVersion" || request == "Language" || request == "SignatureMethod" || request == "AccountInfo" {
													continue
												}
												requestParams[request] = true
											}
										} else if lhsx, ok := lhs.X.(*ast.Ident); ok && lhsx.Name == "request" {
											if rhs, ok := assignStmt.Rhs[0].(*ast.BasicLit); ok {
												value := rhs.Value[1 : len(rhs.Value)-1]
												if lhs.Sel.Name == "ApiName" {
													api.name = value
													api.Name = value
												}
												if lhs.Sel.Name == "Product" {
													api.Product = value
												}
												if lhs.Sel.Name == "Version" {
													api.Version = value
												}
											}
										}
									}
								}
							}
							for param := range requestParams {
								api.Requests = append(api.Requests, Parm{Name: param})
							}
							if len(api.Requests) > 0 {
								if api.name != "" && api.Product != "" {
									*apis = append(*apis, api)
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	/*
		//扫描tea sdk
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.BlockStmt:
				var api = new(Api)
				requestParams := make(map[string]bool)

				for _, stmt := range x.List {
					if assignStmt, ok := stmt.(*ast.IfStmt); ok {
						x := assignStmt.Body
						for _, stmt := range x.List {
							if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
								if lhs, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
									if lhs.Name == "response" {
										if rhs, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
											if len(rhs.Args) == 8 {
												args := rhs.Args[3]
												if args2, ok := args.(*ast.CallExpr); ok {
													if version, ok := args2.Args[0].(*ast.BasicLit); ok {
														api.Version = version.Value[1 : len(version.Value)-1]
													}
												}
											}
										}
									}
								}
							}
						}
					}
					if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
						for _, rhs := range assignStmt.Rhs {
							lhs := assignStmt.Lhs[0]
							request, ok := lhs.(*ast.Ident)
							if !ok {
								continue
							}
							var fun *ast.Ident
							call, ok1 := rhs.(*ast.CallExpr)
							_, ok2 := rhs.(*ast.CompositeLit)
							if ok1 {
								fun, ok = call.Fun.(*ast.Ident)
								if !ok {
									continue
								}
							}

							if ((ok1 && fun.Name == "make") || (ok2)) && request.Name == "request" {
								api.Type = "tea"
								api.FileName = filename
								for _, stmt := range x.List {
									if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
										if lhs, ok := assignStmt.Lhs[0].(*ast.Ident); ok {

											if lhs.Name == "action" {
												if rhs, ok := assignStmt.Rhs[0].(*ast.BasicLit); ok {
													value := rhs.Value[1 : len(rhs.Value)-1]
													api.name = value
													api.Name = value
													continue
												}
											} else if lhs.Name == "response" {
												if rhs, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
													args := rhs.Args[3]
													if args2, ok := args.(*ast.CallExpr); ok {
														if version, ok := args2.Args[0].(*ast.BasicLit); ok {
															api.Version = version.Value[1 : len(version.Value)-1]
														}
													}

												}
											}
										}
										lhs = assignStmt.Lhs[0]
										lhs, ok := lhs.(*ast.IndexExpr)
										if !ok {
											continue
										}
										x := lhs.X
										x1, ok := x.(*ast.Ident)
										if !ok {
											continue
										}
										index := lhs.Index
										index1, ok := index.(*ast.BasicLit)
										if !ok {
											continue
										}

										if x1.Name == "request" {
											attribute := index1.Value[1 : len(index1.Value)-1]
											if attribute == "Product" {
												if rhs, ok := assignStmt.Rhs[0].(*ast.BasicLit); ok {
													value := rhs.Value[1 : len(rhs.Value)-1]
													api.Product = value
												}
											} else {
												request := attribute
												if request == "AccessKeyId" || request == "RegionId" || request == "AccessKeySecret" || request == "Product" || request == "Action" || request == "Version" || request == "ProductName" || request == "Department" || request == "ResourceGroup" {
													continue
												}
												if request == "X-acs-body" || request == "SignatureVersion" || request == "Language" || request == "SignatureMethod" || request == "AccountInfo" {
													continue
												}
												requestParams[request] = true
											}
										}
									}
								}

							}

						}
					}
				}
				for param := range requestParams {
					api.Requests = append(api.Requests, Parm{Name: param})
				}
				if len(api.Requests) > 0 {
					if api.name != "" && api.Product != "" && api.Version != "" {
						*apis = append(*apis, api)
					}
				}
			}

			return true
		})*/
}

func getFuncName(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.SelectorExpr:
		return fn.Sel.Name
	default:
		return ""
	}
}

func getPackageName(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.SelectorExpr:
		if ident, ok := fn.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

func inspectRequest(stmt ast.Stmt, requestParams map[string]bool) {
	ast.Inspect(stmt, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			for _, lhs := range x.Lhs {
				selExpr, ok := lhs.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				xIdent, ok := selExpr.X.(*ast.Ident)
				if ok && xIdent.Name == "request" {
					if selExpr.Sel.Name != "IaasProvider" && selExpr.Sel.Name != "Headers" {
						requestParams[selExpr.Sel.Name] = true
					}
				}
			}
		case *ast.IfStmt:
			for _, stmt := range x.Body.List {
				inspectRequest(stmt, requestParams)
			}
		}
		return true
	})
}
