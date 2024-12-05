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

var TerraformProviderPath string

type Node struct {
	name  string
	child []*Node
}

func GetSchema(fileName string) *Node {
	file, _ := os.Open(fileName)
	fset := token.NewFileSet()

	node, _ := parser.ParseFile(fset, "", file, parser.AllErrors)
	var root *ast.KeyValueExpr

	ast.Inspect(node, func(n ast.Node) bool {
		if returnStmt, ok := n.(*ast.ReturnStmt); ok {
			results := returnStmt.Results
			if len(results) == 0 {
				return true
			}
			unaryExpr := results[0]
			if unaryExpr, ok := unaryExpr.(*ast.UnaryExpr); ok {
				if x, ok := unaryExpr.X.(*ast.CompositeLit); ok {
					elts := x.Elts
					for _, kvexpr := range elts {
						if Exprkv, ok := kvexpr.(*ast.KeyValueExpr); ok {
							if Exprk, ok := Exprkv.Key.(*ast.Ident); ok && Exprk.Name == "Schema" {
								root = Exprkv
							}
						}
					}
				}
			}
		}
		return true
	})
	rootNode := new(Node)
	rootNode.child = make([]*Node, 0)

	dfs(root, rootNode)
	return rootNode
}

func GetDsetAttributes(fileName string) []string {
	attributes := make([]string, 0)
	file, _ := os.Open(fileName)
	fset := token.NewFileSet()

	node, _ := parser.ParseFile(fset, "", file, parser.AllErrors)

	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
				if x, ok := fun.X.(*ast.Ident); ok {
					if x.Name == "d" && fun.Sel.Name == "Set" {
						args := call.Args
						if len(args) == 2 {
							if attr, ok := args[0].(*ast.BasicLit); ok {
								value := attr.Value[1 : len(attr.Value)-1]
								attributes = append(attributes, value)
							}
						}
					}
				}
			}
		}
		return true
	})

	ast.Inspect(node, func(n ast.Node) bool {
		if assignment, ok := n.(*ast.AssignStmt); ok {
			if len(assignment.Lhs) == 1 {
				if left, ok := assignment.Lhs[0].(*ast.Ident); ok {
					if left.Name == "mapping" {
						if len(assignment.Rhs) == 1 {
							if compositelit, ok := assignment.Rhs[0].(*ast.CompositeLit); ok {
								if typ, ok := compositelit.Type.(*ast.MapType); ok {
									if _, ok := typ.Value.(*ast.InterfaceType); ok {
										for _, elt := range compositelit.Elts {
											if kv, ok := elt.(*ast.KeyValueExpr); ok {
												if attr, ok := kv.Key.(*ast.BasicLit); ok {
													value := attr.Value[1 : len(attr.Value)-1]
													attributes = append(attributes, value)
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
		return true
	})

	return attributes
}

func dfs(left *ast.KeyValueExpr, leftNode *Node) {
	LeftValue, _ := left.Value.(*ast.CompositeLit)
	LeftElts := &LeftValue.Elts

	for _, leftAttribute := range *LeftElts {
		leftAttribute, _ := leftAttribute.(*ast.KeyValueExpr)
		//判断leftAttribute是文件节点还是目录节点
		var ifSchema bool
		leftKey, _ := leftAttribute.Key.(*ast.BasicLit)
		leftAttributeValue, ok := leftAttribute.Value.(*ast.CompositeLit)
		if ok {
			elts := leftAttributeValue.Elts
			for i := 0; i < len(elts); i++ {
				if kv, ok := elts[i].(*ast.KeyValueExpr); ok {
					if k, ok := kv.Key.(*ast.Ident); ok && k.Name == "Elem" {
						ifSchema = true
						break
					}
				}
			}
		}

		childNode := new(Node)
		childNode.name = strings.TrimSuffix(strings.TrimPrefix(leftKey.Value, "\""), "\"")
		childNode.child = make([]*Node, 0)
		leftNode.child = append(leftNode.child, childNode)

		//如果是目录类型
		if ifSchema {
			var leftChild *ast.KeyValueExpr
			//提取出左孩子
			leftAttributeValue, _ := leftAttribute.Value.(*ast.CompositeLit)
			elts := leftAttributeValue.Elts
			for _, kvExpr := range elts {
				kvExpr := kvExpr.(*ast.KeyValueExpr)
				if kExpr, ok := kvExpr.Key.(*ast.Ident); ok && kExpr.Name == "Elem" {
					kExprValue := kvExpr.Value
					if unaryExpr, ok := kExprValue.(*ast.UnaryExpr); ok {
						x := unaryExpr.X.(*ast.CompositeLit)
						for _, e := range x.Elts {
							if leftchild, ok := e.(*ast.KeyValueExpr); ok {
								if leftchildKey, ok := leftchild.Key.(*ast.Ident); ok && leftchildKey.Name == "Schema" {
									leftChild = leftchild
								}
							}

						}
					}
				}
			}
			if leftChild != nil {
				dfs(leftChild, childNode)
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请传入一个参数：terraform项目地址")
		os.Exit(1)
	}
	TerraformProviderPath = os.Args[1]
	dirname := TerraformProviderPath + "/alibabacloudstack" // your directory containing Go source files
	f, _ := os.Create("output.txt")
	filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if name == "data_source_apsarastack_common.go" {
			fmt.Println()
		}
		if !info.IsDir() && !strings.HasSuffix(path, "test.go") && strings.HasSuffix(path, "go") { // process only Go files
			if strings.HasPrefix(name, "data") || strings.HasPrefix(name, "resource") {
				testFilename := strings.TrimSuffix(path, ".go")
				testFilename = testFilename + "_test.go"
				_, err := os.Stat(testFilename)
				if err != nil {
					return nil
				}
				root := GetSchema(path)
				attributes := GetDsetAttributes(path)
				var dfs func(*Node, string) bool
				dfs = func(n *Node, s string) bool {
					if s == n.name {
						return true
					}
					for _, child := range n.child {
						if dfs(child, s) == true {
							return true
						}
					}
					return false
				}
				for _, attribute := range attributes {
					if !dfs(root, attribute) {
						fmt.Fprintln(f, path, attribute)
					}
				}
			}
		}
		return nil
	})
	fmt.Println()

}
