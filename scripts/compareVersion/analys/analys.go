package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Rpc struct {
	Name  string `json:"Name"`
	State string `json:"State"`
}

type Version struct {
	Version string `json:"Version"`
	Rpcs    []Rpc  `json:"Rpcs"`
}

func addReadme(versions []Version) {
	// 生成表格
	table := generateTable(versions)

	// 读取README文件
	readmeContent, err := os.ReadFile("../../../README.md")
	if err != nil {
		fmt.Println("Error reading README.md:", err)
		return
	}

	// 找到插入点的位置
	insertPointIndex := bytes.Index(readmeContent, []byte("<!-- INSERT TABLE HERE -->"))
	if insertPointIndex == -1 {
		fmt.Println("Insert point not found in README.md")
		return
	}

	// 删除插入点后面的内容
	newReadmeContent := append(readmeContent[:insertPointIndex+len("<!-- INSERT TABLE HERE -->")+2], []byte("\n"+table)...)

	// 写入新的README文件
	err = os.WriteFile("../../../README.md", newReadmeContent, 0644)
	if err != nil {
		fmt.Println("Error writing to README.md:", err)
		return
	}

	fmt.Println("Table successfully inserted into README.md")
}

func generateTable(versions []Version) string {
	var buffer bytes.Buffer

	// 添加 Rpc 名称作为第一列表头
	buffer.WriteString(":white_check_mark::当前功能被Provider支持")
	buffer.WriteString("\n")

	buffer.WriteString(":x::当前功能在该Provider存在风险")
	buffer.WriteString("\n")

	buffer.WriteString(":no_entry_sign::当前功能在该Provider下不可用")
	buffer.WriteString("\n")

	buffer.WriteString("| Rpc Name ")
	for _, v := range versions {
		buffer.WriteString(fmt.Sprintf(" | %s ", v.Version))
	}
	buffer.WriteString("|\n")

	// 分割线
	buffer.WriteString("| --- ")
	for range versions {
		buffer.WriteString(" | --- ")
	}
	buffer.WriteString("|\n")

	// 数据行
	var AllRpcs []Rpc
	for _, version := range versions {
		for _, Rpc := range version.Rpcs {
			var lock bool
			for _, rpc := range AllRpcs {
				if Rpc == rpc {
					lock = true
				}
			}
			if !lock {
				AllRpcs = append(AllRpcs, Rpc)
			}
		}
	}
	for _, rpc := range AllRpcs {
		buffer.WriteString(fmt.Sprintf("| %s ", rpc.Name))
		for _, v := range versions {
			var lock bool
			for _, r := range v.Rpcs {
				if r.Name == rpc.Name {
					buffer.WriteString(fmt.Sprintf(" | %s ", r.State))
					lock = true
					break
				}
			}
			if !lock {
				buffer.WriteString(fmt.Sprintf(" | %s ", " "))
			}
		}
		buffer.WriteString("|\n")
	}

	return buffer.String()
}
func main() {
	// 打开包含JSON数据的文件
	file, err := os.Open("../versions.json") // 假设文件名为data.json
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var versions []Version

	// 逐行扫描文件内容
	for scanner.Scan() {
		line := scanner.Text()
		var v Version

		// 解码每一行的JSON数据到Version结构体
		err := json.Unmarshal([]byte(line), &v)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}

		// 将解析后的Version添加到切片中
		versions = append(versions, v)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

	addReadme(versions)
}
