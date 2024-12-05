package main

import (
	"bufio"
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

func scanDataTest(fileName string) *Node {
	file, _ := os.Open(fileName)
	fset := token.NewFileSet()

	node, _ := parser.ParseFile(fset, "", file, parser.ParseComments)
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

func hasTested(filename string, dataOut []string) int {
	inputFile, _ := os.Open(filename)
	defer inputFile.Close()
	// 使用bufio.Scanner逐行读取文件
	scanner := bufio.NewScanner(inputFile)
	mp := make(map[string]interface{})
	for scanner.Scan() {
		line := scanner.Text()
		for _, s := range dataOut {
			if strings.Contains(line, s) {
				mp[s] = struct{}{}
			}
		}
	}
	res := 0
	for _, s := range dataOut {
		_, ok := mp[s]
		if ok {
			res++
		}
	}

	return res
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("请传入一个参数：terraform项目地址")
		os.Exit(1)
	}
	TerraformProviderPath = os.Args[1]
	dirname := TerraformProviderPath + "/alibabacloudstack" // your directory containing Go source files
	fmt.Println(dirname)
	filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()

		if !info.IsDir() && !strings.HasSuffix(path, "test.go") && strings.HasSuffix(path, "go") { // process only Go files
			if strings.HasPrefix(name, "data") {

				testFilename := strings.TrimSuffix(path, ".go")
				testFilename = testFilename + "_test.go"
				_, err := os.Stat(testFilename)
				if err != nil {
					fmt.Printf("[0.00%%] %s: 缺失测试文件\n", name)
					return nil
				}
				out := strings.TrimSuffix(path, ".go")
				outstrs := strings.Split(out, "_")
				probablyOut := make([]string, 0)
				probablyOut = append(probablyOut, outstrs[len(outstrs)-1])
				pout2 := strings.TrimSuffix(strings.TrimPrefix(name, "data_source_apsarastack_"), ".go")
				probablyOut = append(probablyOut, pout2)
				pout3 := outstrs[len(outstrs)-2] + "_" + outstrs[len(outstrs)-1]
				probablyOut = append(probablyOut, pout3)

				root := scanDataTest(path)
				count := 0
				dataOut := make([]string, 0)
				for _, attribute := range root.child {
					var lock bool
					for _, s := range probablyOut {
						if attribute.name == s {
							lock = true
						}
					}
					if len(attribute.child) != 0 && lock {
						count = len(attribute.child)
						for _, child := range attribute.child {
							dataOut = append(dataOut, child.name)
						}
					}
				}
				if count == 0 {
					count = len(root.child)
					for _, child := range root.child {
						dataOut = append(dataOut, child.name)
					}
				}

				if count != 0 {
					tested := hasTested(testFilename, dataOut)
					percentage := (float64(tested) / float64(count)) * 100

					fmt.Printf("[%.2f%%] %s: %d / %d\n", percentage, name, tested, count)
				}
			}
		}
		return nil
	})
	fmt.Println()
	// ans := true
	filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if !info.IsDir() && !strings.HasSuffix(path, "test.go") && strings.HasSuffix(path, "go") { // process only Go files
			if strings.HasPrefix(name, "resource") {
				testFilename := strings.TrimSuffix(path, ".go")
				testFilename = testFilename + "_test.go"
				_, err := os.Stat(testFilename)
				if err != nil {
					fmt.Printf("[0.00%%] %s: 缺失测试文件\n", name)
					return nil
				}
				root := scanDataTest(path)
				count := 0
				dataOut := make([]string, 0)

				count = len(root.child)
				for _, child := range root.child {
					dataOut = append(dataOut, child.name)
				}

				if count != 0 {
					tested := hasTested(testFilename, dataOut)
					percentage := (float64(tested) / float64(count)) * 100

					fmt.Printf("[%.2f%%] %s: %d / %d\n", percentage, name, tested, count)
					if tested != count {
						// ans = false
					}
				}
			}
		}
		return nil
	})
	fmt.Println()
}
