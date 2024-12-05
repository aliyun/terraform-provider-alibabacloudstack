package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Method struct {
	Name     string
	Request  string
	Response string
}

type Result struct {
	Name      string
	Exist     bool
	Requests  []string
	Responses []string
}

// Rpc 结构体定义
type Rpc struct {
	Name  string // Rpc名称
	State string // 风险、安全、无
}

// Version 结构体定义
type Version struct {
	Version string // 版本
	Rpcs    []Rpc  // 所有的Rpc
}

func getMethodInfo(field *ast.Field) *Method {
	var method *Method = &Method{}
	method.Name = field.Names[0].Name
	var lock1 bool
	var lock2 bool
	if tp, ok := field.Type.(*ast.FuncType); ok {
		lists := tp.Params.List
		if startExpr, ok := lists[1].Type.(*ast.StarExpr); ok {
			if request, ok := startExpr.X.(*ast.Ident); ok {
				method.Request = request.Name
				lock1 = true
			}
		}
		results := tp.Results.List
		if startExpr, ok := results[0].Type.(*ast.StarExpr); ok {
			if response, ok := startExpr.X.(*ast.Ident); ok {
				method.Response = response.Name
				lock2 = true
			}
		}
		if lock1 && lock2 {
			return method
		}
	}
	return nil
}
func getMethodsInfo(interfacetype *ast.InterfaceType) []*Method {
	list := interfacetype.Methods.List
	methods := make([]*Method, 0)
	for _, field := range list {
		method := getMethodInfo(field)
		if method != nil {
			methods = append(methods, method)
		}
	}
	return methods
}
func GetAllRpcAndRequest() ([]*Method, map[string][]string) {
	filename := "terraform/internal/tfplugin5/tfplugin5.pb.go"
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	//扫描得到所有的RPC调用
	var methods []*Method
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == "ProviderServer" {
				if interfaceType, ok := x.Type.(*ast.InterfaceType); ok {
					methods = getMethodsInfo(interfaceType)
					// for _, m := range methods {
					// 	fmt.Println(m)
					// }
				}
			}
		}
		return true
	})

	//扫描得到所有的Requests 和 Response结构体
	Params := make(map[string][]string)
	ast.Inspect(node, func(n ast.Node) bool {
		structType, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if the name ends with "Response"
		if strings.HasSuffix(structType.Name.Name, "Response") || strings.HasSuffix(structType.Name.Name, "Request") {
			if _, ok := Params[structType.Name.Name]; !ok {
				Params[structType.Name.Name] = make([]string, 0)
			}
			// Check the struct type
			if structDecl, ok := structType.Type.(*ast.StructType); ok {
				for _, field := range structDecl.Fields.List {
					// Check if the field has a tag
					if field.Tag != nil {
						for _, fieldName := range field.Names {
							Params[structType.Name.Name] = append(Params[structType.Name.Name], fieldName.Name)
						}
					}
				}
			}
		}

		return true
	})
	return methods, Params
}

func GetAllGrpcProvider() ([]*Method, map[string][]string) {
	//filename := "terraform-provider-apsarastack/vendor/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema/grpc_provider.go"

	filename := "../../vendor/github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server/server.go" // your directory containing Go source files
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	var methods []*Method

	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check if the function has exactly two parameters and two results.
		if len(funcDecl.Type.Params.List) != 2 || len(funcDecl.Type.Results.List) != 2 {
			return true
		}

		// Extract second parameter and first result.
		secondParam := funcDecl.Type.Params.List[1]
		firstResult := funcDecl.Type.Results.List[0]

		// Ensure second parameter is *ast.StarExpr and X is *ast.SelectorExpr.
		secondParamType, ok := secondParam.Type.(*ast.StarExpr)
		if !ok {
			return true
		}
		secondParamSelector, ok := secondParamType.X.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Ensure first result is *ast.StarExpr and X is *ast.SelectorExpr.
		firstResultType, ok := firstResult.Type.(*ast.StarExpr)
		if !ok {
			return true
		}
		firstResultSelector, ok := firstResultType.X.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// Print relevant information.
		var method = new(Method)
		method.Name = funcDecl.Name.Name
		method.Request = secondParamSelector.Sel.Name
		method.Response = firstResultSelector.Sel.Name
		methods = append(methods, method)

		return true
	})
	// for _, m := range methods {
	// 	fmt.Println(m)
	// }

	//扫描provider中的所有request 和 response
	dirname := "../../vendor/github.com/hashicorp/terraform-plugin-go/tfprotov5/" // your directory containing Go source files
	Params := make(map[string][]string)

	filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "test.go") { // process only Go files
			filename = path
			fs = token.NewFileSet()
			node, err = parser.ParseFile(fs, filename, nil, parser.ParseComments)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			//扫描得到所有的Requests 和 Response结构体
			ast.Inspect(node, func(n ast.Node) bool {
				structType, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}

				// Check if the name ends with "Response"
				if strings.HasSuffix(structType.Name.Name, "Response") || strings.HasSuffix(structType.Name.Name, "Request") {
					if _, ok := Params[structType.Name.Name]; !ok {
						Params[structType.Name.Name] = make([]string, 0)
					}
					// Check the struct type
					if structDecl, ok := structType.Type.(*ast.StructType); ok {
						for _, field := range structDecl.Fields.List {
							// Check if the field has a tag

							for _, fieldName := range field.Names {
								Params[structType.Name.Name] = append(Params[structType.Name.Name], fieldName.Name)
							}

						}
					}
				}

				return true
			})
		}
		return nil
	})

	return methods, Params

}
func Tojson(v1 Version) {
	// 创建或打开一个JSON文件
	file, err := os.OpenFile("versions.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// 创建一个新的Version实例

	// 将Version实例追加到JSON文件中
	enc := json.NewEncoder(file)
	err = enc.Encode(v1)
	if err != nil {
		fmt.Println(err)
		return
	}
}
func GetAllRpcAndRequest2() ([]*Method, map[string][]string) {
	filename := "opentofu/internal/tfplugin5/tfplugin5.pb.go"
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	//扫描得到所有的RPC调用
	var methods []*Method
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == "ProviderServer" {
				if interfaceType, ok := x.Type.(*ast.InterfaceType); ok {
					methods = getMethodsInfo(interfaceType)
					// for _, m := range methods {
					// 	fmt.Println(m)
					// }
				}
			}
		}
		return true
	})

	//扫描得到所有的Requests 和 Response结构体
	Params := make(map[string][]string)
	ast.Inspect(node, func(n ast.Node) bool {
		structType, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if the name ends with "Response"
		if strings.HasSuffix(structType.Name.Name, "Response") || strings.HasSuffix(structType.Name.Name, "Request") {
			if _, ok := Params[structType.Name.Name]; !ok {
				Params[structType.Name.Name] = make([]string, 0)
			}
			// Check the struct type
			if structDecl, ok := structType.Type.(*ast.StructType); ok {
				for _, field := range structDecl.Fields.List {
					// Check if the field has a tag
					if field.Tag != nil {
						for _, fieldName := range field.Names {
							Params[structType.Name.Name] = append(Params[structType.Name.Name], fieldName.Name)
						}
					}
				}
			}
		}

		return true
	})
	return methods, Params
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide at least one argument.")
		return
	}
	var coreMethods []*Method
	var coreParms map[string][]string
	var version Version
	version.Version = os.Args[1]
	if strings.HasPrefix(os.Args[1], "terraform") {
		coreMethods, coreParms = GetAllRpcAndRequest()
	} else if strings.HasPrefix(os.Args[1], "opentofu") {
		coreMethods, coreParms = GetAllRpcAndRequest2()
	} else {
		fmt.Println("input error")
		return
	}
	providerMethods, providerParms := GetAllGrpcProvider()

	methodMap := make(map[string]*Method)
	for _, m := range providerMethods {
		methodMap[m.Name] = m
	}
	okRpcs := make([]*Result, 0)
	notExistRpcs := make([]*Result, 0)
	riskRpcs := make([]*Result, 0)
	for _, m := range coreMethods {
		result := new(Result)
		result.Name = m.Name
		if _, ok := methodMap[m.Name]; !ok {
			result.Exist = false
			notExistRpcs = append(notExistRpcs, result)
		} else {
			result.Exist = true
			cm := methodMap[m.Name]
			requestMap := make(map[string]interface{})
			responseMap := make(map[string]interface{})
			for _, s := range providerParms[cm.Request] {
				requestMap[s] = struct{}{}
			}
			for _, s := range providerParms[cm.Response] {
				responseMap[s] = struct{}{}
			}
			for _, s := range coreParms[m.Request] {
				if _, ok := requestMap[s]; !ok {
					result.Requests = append(result.Requests, s)
				}
			}
			for _, s := range coreParms[m.Response] {
				if _, ok := responseMap[s]; !ok {
					result.Responses = append(result.Responses, s)
				}
			}
			if len(result.Requests) == 0 && len(result.Responses) == 0 {
				okRpcs = append(okRpcs, result)
			} else {
				riskRpcs = append(riskRpcs, result)
			}
		}
	}

	for _, result := range okRpcs {
		fmt.Println("无风险RPC:", result.Name)
		version.Rpcs = append(version.Rpcs, Rpc{Name: result.Name, State: ":white_check_mark:"})
	}
	fmt.Println()

	for _, result := range notExistRpcs {
		fmt.Println("未实现RPC:", result.Name)
		version.Rpcs = append(version.Rpcs, Rpc{Name: result.Name, State: ":no_entry_sign:"})

	}
	fmt.Println()

	for _, result := range riskRpcs {
		fmt.Println("风险RPC:", result.Name)
		version.Rpcs = append(version.Rpcs, Rpc{Name: result.Name, State: ":x:"})

		if len(result.Requests) != 0 {
			fmt.Println("缺失入参:")
			for _, request := range result.Requests {
				fmt.Println(request)
			}
		}
		if len(result.Responses) != 0 {
			fmt.Println("缺失出餐:")
			for _, response := range result.Responses {
				fmt.Println(response)
			}
		}
		fmt.Println()
	}
	Tojson(version)
}
