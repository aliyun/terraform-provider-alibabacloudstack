package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getProductExact(product string) string {
	url := "http://apsara-basicservice-center.aliyun-inc.com/popcodes"
	queryParms := "/?name=" + product
	fullURL := url + queryParms
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		return ""
	}
	request.Header.Add("Authorization", "Token 46723f12ae3fd2f17588a1b14b256dca7f7590a7")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %v\n", err)
		return ""
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("Request failed with status: %s\n", response.Status)
		return ""
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		return ""
	}

	var data ApiResponse1
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Failed to parse JSON: %v\n", err)
		return ""
	}
	return data.Results[0].Uuid
}

func GetApiInfo(name string, product string, version string) *Api {
	var api Api
	api.Name = name
	api.name = name
	api.Version = version
	api.Product = product
	url := "http://apsara-basicservice-center.aliyun-inc.com/basicservices-popapi"
	queryParams := "/?arch_types=x86%2A&edition=019b2be2-d98e-4178-95fe-c343d9d7ebe9&%40name__iexact=" + api.Name + "&%40product__exact=" + getProductExact(api.Product) + "&%40version__iexact=" + api.Version
	fullURL := url + queryParams

	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		return nil
	}

	request.Header.Add("Authorization", "Token 46723f12ae3fd2f17588a1b14b256dca7f7590a7")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %v\n", err)
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Request failed with status: %s\n", response.Status)
		return nil
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		return nil
	}

	var data1 ApiResponse1
	var data2 ApiResponse2
	var data3 ApiResponse3

	err1 := json.Unmarshal(body, &data1)
	err2 := json.Unmarshal(body, &data2)
	err3 := json.Unmarshal(body, &data3)
	if err3 == nil {
		for _, result := range data3.Results {
			if len(result.Parameters.ParameterGroups) == 0 {
				err3 = fmt.Errorf("wrong")
			}
		}
	}
	switch {

	case err3 == nil:

		mp := make(map[string]interface{})
		for _, result := range data3.Results {
			//log.Printf("UUID: %s, Name: %s, Status: %s\n", result.Uuid, result.Name, result.Status)
			p := result.Parameters
			for _, paramGroup := range p.ParameterGroups {
				for _, param := range paramGroup.Parameters {
					_, ok := mp[param.TagName]
					if ok {
						continue
					}
					newParam := &Parm{
						Name: param.TagName,
						Type: param.Type,
					}
					api.Requests = append(api.Requests, *newParam)
					//log.Printf("Param Name: %s, Type: %s\n", param.TagName, param.Type)
					mp[param.TagName] = struct{}{}
				}
			}
		}
		return &api

	case err1 == nil:
		for _, result := range data1.Results {
			//log.Printf("UUID: %s, Name: %s, Status: %s\n", result.Uuid, result.Name, result.Status)
			p := result.Parameters
			for _, param := range p.Parameter {
				newParam := &Parm{
					Name: param.TagName,
					Type: param.Type,
				}
				api.Requests = append(api.Requests, *newParam)
				//log.Printf("Param Name: %s, Type: %s\n", param.TagName, param.Type)
			}
			m := result.ResultMapping.Member
			api.Responses = append(api.Responses, Parm{
				Name: m.TagName,
				Type: m.Type,
			})
			//log.Printf("Response Name: %s, Type: %s\n", m.TagName, m.Type)
		}
		return &api

	case err2 == nil:
		for _, result := range data2.Results {
			//log.Printf("UUID: %s, Name: %s, Status: %s\n", result.Uuid, result.Name, result.Status)
			p := result.Parameters
			for _, param := range p.Parameter {
				newParam := &Parm{
					Name: param.TagName,
					Type: param.Type,
				}
				api.Requests = append(api.Requests, *newParam)
				//log.Printf("Param Name: %s, Type: %s\n", param.TagName, param.Type)
			}
			m := result.ResultMapping.Member
			for i := 0; i < len(m); i++ {
				api.Responses = append(api.Responses, Parm{
					Name: m[i].TagName,
					Type: m[i].Type,
				})
			}
			//log.Printf("Response Name: %s, Type: %s\n", m.TagName, m.Type)
		}
		return &api

	default:
		log.Printf("解析错误: %v\n", api.name)
		log.Printf("Failed to parse JSON: %v\n", err1)
		log.Printf("Failed to parse JSON: %v\n", err2)
		return nil
	}
}
