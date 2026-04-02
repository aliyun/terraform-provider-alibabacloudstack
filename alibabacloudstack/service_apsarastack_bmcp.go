package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

type BcmpService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *BcmpService) ListEvpc(evpcName string) ([]map[string]interface{}, error) {
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

func (s *BcmpService) DoEasyAIListEvpcRequest(id string) (map[string]interface{}, error) {
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

func (s *BcmpService) DoEasyAIListSecurityGroupRequest(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{}
	if id != "" {
		request["SgId"] = id
	}

	raw, err := s.client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListSecurityGroup", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	response := raw
	securityGroupList := response["data"].([]interface{})

	for _, securityGroup := range securityGroupList {
		return securityGroup.(map[string]interface{}), nil
	}

	return nil, nil
}

func (s *BcmpService) DoEasyAIListKeyPairRequest(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   100,
	}

	raw, err := s.client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListKeyPair", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	response := raw
	data := response["data"].([]interface{})
	keyPairId := id

	for _, item := range data {
		itemMap := item.(map[string]interface{})
		if itemMap["name"].(string) == keyPairId {
			return itemMap, nil
		}
	}
	return nil, nil
}