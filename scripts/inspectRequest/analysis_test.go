/*
compare
无论requests还是response，都希望api可以覆盖sdk壳子的参数
如果出现：

	sdk的参数不存在于api中，则标记该api缺少参数

否则：

	视为api暂时没问题
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/panjf2000/ants/v2"
)

type FileApis map[string][]*Api

var File = make(FileApis)

func WriteToJson(data FileApis, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON file: %v", err)
	}

	return nil
}

func WriteApisToJson(data []*Api, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %v", err)
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON file: %v", err)
	}

	return nil
}

// ReadFileFromJson 从 JSON 文件中读取并解析成 File 类型的数据
func ReadFromJson(filename string) (FileApis, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening JSON file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var parsedData FileApis
	err = json.Unmarshal(byteValue, &parsedData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	return parsedData, nil
}
func ReadApisFromJson(filename string) ([]*Api, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening JSON file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %v", err)
	}

	var parsedData []*Api
	err = json.Unmarshal(byteValue, &parsedData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data: %v", err)
	}

	return parsedData, nil
}

func TestHighAndCommon(t *testing.T) {
	logInit()
	intersectionApis := make([]*Api, 0)
	firstScanApis := GetApisFromScanFiles()
	scanApis := make([]*Api, 0)
	for _, api := range firstScanApis {
		if api.Type == "common" {
			intersectionApis = append(intersectionApis, api)
		} else if api.Type == "high" {
			scanApis = append(scanApis, api)
		}
	}

	alibabaCloudApis := manifestAlibabaCloudSdk()
	for _, sapi := range scanApis {
		var intersectionApi *Api

		for _, alibabaCloudApi := range alibabaCloudApis {
			if strings.EqualFold(sapi.Name, alibabaCloudApi.Name) && strings.EqualFold(sapi.Product, alibabaCloudApi.Product) {
				intersectionApi = new(Api)
				intersectionApi.name = alibabaCloudApi.name
				intersectionApi.Name = alibabaCloudApi.name
				intersectionApi.Product = alibabaCloudApi.Product
				intersectionApi.Version = alibabaCloudApi.Version
				intersectionApi.FileName = sapi.FileName
				intersectionApi.Type = sapi.Type
				for _, alrq := range alibabaCloudApi.Requests {
					for _, srq := range sapi.Requests {
						if alrq.Name == srq.Name {
							alrq.Name = alrq.RealName
							intersectionApi.Requests = append(intersectionApi.Requests, alrq)
						}
					}
				}
				for _, alrp := range alibabaCloudApi.Responses {
					for _, srp := range sapi.Responses {
						if alrp.Name == srp.Name {
							alrp.Name = alrp.RealName
							intersectionApi.Responses = append(intersectionApi.Responses, alrp)
						}
					}
				}
			}
		}
		if intersectionApi != nil {
			intersectionApis = append(intersectionApis, intersectionApi)
		}
	}

	all := 0

	apsarastackApis, _ := ReadApisFromJson("output/apsaraApis.json")
	log.Println("Finish all getting tasks")
	for _, intersectionApi := range intersectionApis {
		/*
			通过name 和 product 锁定到检索到的api
			进行参数取交集的操作
		*/
		var apapi *Api
		found := false
		for _, apapi = range apsarastackApis {
			if apapi.Name == intersectionApi.Name && apapi.Product == intersectionApi.Product && apapi.Version == intersectionApi.Version {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		if intersectionApi == nil {
			continue
		}
		if apapi == nil {
			continue
		}
		if len(apapi.Requests) == 0 && len(apapi.Responses) == 0 {
			continue
		}

		apsarastackApi := apapi

		log.Println(all+1, ".")
		apsarastackApi.show()
		intersectionApi.show()

		//compare resquests
		RequestMap := make(map[string]interface{})
		ResponseMap := make(map[string]interface{})
		for _, request := range apsarastackApi.Requests {
			RequestMap[request.Name] = struct{}{}
		}
		for _, response := range apsarastackApi.Responses {
			ResponseMap[response.Name] = struct{}{}
		}

		all++
		var bug bool
		var api *Api
		for _, request := range intersectionApi.Requests {
			_, ok := RequestMap[request.Name]
			if !ok {
				_, ok = File[intersectionApi.FileName]
				if !ok {
					File[intersectionApi.FileName] = make([]*Api, 0)
				}
				if bug == false {
					api = new(Api)
					api.Name = intersectionApi.Name
					api.FileName = intersectionApi.FileName
					api.Type = intersectionApi.Type
					api.Product = intersectionApi.Product
					api.Version = intersectionApi.Version
					bug = true
				}
				api.Requests = append(api.Requests, request)
			}
		}
		if api != nil {
			File[intersectionApi.FileName] = append(File[intersectionApi.FileName], api)
		}
	}

	//标准打输出位置
	f, _ := os.Create("output.txt")
	for k, v := range File {
		k = strings.TrimPrefix(k, "../../alibabacloudstack/")
		fmt.Fprintln(f, "文件:", k)
		for _, api := range v {
			fmt.Fprintln(f, "  方法:", api.Name)
			for _, request := range api.Requests {
				fmt.Fprintln(f, "    参数:", request.Name, "缺失")
			}

		}
		fmt.Fprintln(f)
	}
	ans := true
	for k, v := range File {
		k = strings.TrimPrefix(k, "../../alibabacloudstack/")
		fmt.Println("文件:", k)
		for _, api := range v {
			fmt.Println("  方法:", api.Name)
			for _, request := range api.Requests {
				fmt.Println("    参数:", request.Name, "缺失")
				ans = false
			}

		}
		fmt.Println()
	}
	err := WriteToJson(File, "output/FileApis.json")
	ResourceCount := len(File)
	Sdkcount := 0
	for _, x := range File {
		Sdkcount += len(x)
	}

	log.Println("一共有", ResourceCount, "个资源不兼容")
	log.Println("一共有", Sdkcount, "个sdk不兼容")

	if err != nil {
		panic(err)
	}
	if !ans {
		t.Errorf("文档自动化检查不通过")
	}
}

func TestGetAllApsaraSdks(t *testing.T) {
	logInit()
	intersectionApis := make([]*Api, 0)
	firstScanApis := GetApisFromScanFiles()
	scanApis := make([]*Api, 0)
	for _, api := range firstScanApis {
		if api.Type == "common" {
			intersectionApis = append(intersectionApis, api)
		} else if api.Type == "high" {
			scanApis = append(scanApis, api)
		}
	}

	alibabaCloudApis := manifestAlibabaCloudSdk()
	for _, sapi := range scanApis {
		var intersectionApi *Api

		for _, alibabaCloudApi := range alibabaCloudApis {
			if strings.EqualFold(sapi.Name, alibabaCloudApi.Name) && strings.EqualFold(sapi.Product, alibabaCloudApi.Product) {
				intersectionApi = new(Api)
				intersectionApi.name = alibabaCloudApi.name
				intersectionApi.Name = alibabaCloudApi.name
				intersectionApi.Product = alibabaCloudApi.Product
				intersectionApi.Version = alibabaCloudApi.Version
				intersectionApi.FileName = sapi.FileName
				for _, alrq := range alibabaCloudApi.Requests {
					for _, srq := range sapi.Requests {
						if alrq.Name == srq.Name {
							intersectionApi.Requests = append(intersectionApi.Requests, alrq)
						}
					}
				}
				for _, alrp := range alibabaCloudApi.Responses {
					for _, srp := range sapi.Responses {
						if alrp.Name == srp.Name {
							intersectionApi.Responses = append(intersectionApi.Responses, alrp)
						}
					}
				}
			}
		}
		if intersectionApi != nil {
			intersectionApis = append(intersectionApis, intersectionApi)
		}
	}

	apsarastackApis := make([]*Api, len(intersectionApis))

	log.Println("Geting apsarastack Apis")
	var wg sync.WaitGroup
	pool, err := ants.NewPoolWithFunc(10, func(i interface{}) {
		index := i.(int)
		apsarastackApis[index] = GetApiInfo(intersectionApis[index].Name, intersectionApis[index].Product, intersectionApis[index].Version)
		log.Printf("Got %s", intersectionApis[index].Name)
		wg.Done()
	})
	if err != nil {
		log.Fatalf("Failed to create ants pool: %v", err)
	}
	defer pool.Release()

	for i := 0; i < len(intersectionApis); i++ {
		wg.Add(1)
		if err := pool.Invoke(i); err != nil {
			log.Fatalf("Failed to submit task to ants pool: %v", err)
			wg.Done()
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("Finish all getting tasks")
	WriteApisToJson(apsarastackApis, "output/apsaraApis.json")
}

func TestModifyCase(t *testing.T) {
	logInit()
	dirname := "../../alibabacloudstack" // your directory containing Go source files
	File, err := ReadFromJson("output/FileApis.json")
	apsarastackApis, _ := ReadApisFromJson("output/apsaraApis.json")
	if err != nil {
		panic(err)
	}
	for fileName, apis := range File {
		for _, api := range apis {
			var apsaraApi *Api
			for _, apapi := range apsarastackApis {
				if apapi.Name == api.Name && apapi.Product == api.Product && apapi.Version == api.Version {
					apsaraApi = apapi
					break
				}
			}
			err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if path == fileName {
					if api.Type == "common" {
						ModifyCase(path, api, apsaraApi)
					}
				}
				return nil
			})

			if err != nil {
				panic(err)
			}
		}
	}
}
