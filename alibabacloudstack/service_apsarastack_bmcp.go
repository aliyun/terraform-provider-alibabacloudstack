package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

type EvpcService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *EvpcService) ListEvpc(evpcName string) ([]map[string]interface{}, error) {
	request := map[string]interface{}{}
	if evpcName != "" {
		request["EvpcName"] = evpcName
	}

	raw, err := s.client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListEvpc", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	response := raw
	evpcList := response["EvpcList"].([]interface{})

	var result []map[string]interface{}
	for _, evpc := range evpcList {
		result = append(result, evpc.(map[string]interface{}))
	}

	return result, nil
}

func (s *EvpcService) DoEasyAIListEvpcRequest(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{}
	if id != "" {
		request["EvpcId"] = id
	}

	raw, err := s.client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListEvpc", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	response := raw
	evpcList := response["EvpcList"].([]interface{})

	for _, evpc := range evpcList {
		return evpc.(map[string]interface{}), nil
	}

	return nil, nil
}
