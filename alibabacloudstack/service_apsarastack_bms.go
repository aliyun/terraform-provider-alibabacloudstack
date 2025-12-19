package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type BmsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *BmsService) DescribeKeyPair(id string) (map[string]interface{}, error) {
	action := "ListKeyPair"
	request := map[string]interface{}{
		"KeyPairName": id,
		"PageNumber":  1,
		"PageSize":    100,
		"DeployType":  "bms",
	}

	resp, err := s.client.DoTeaRequest("POST", "bms", "2024-03-01", action, "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	keyPairs, err := jsonpath.Get("$.data.data", resp)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "ListKeyPair", "$.data.data", resp)
	}
	for _, v := range keyPairs.([]interface{}) {
		keyPair := v.(map[string]interface{})
		if keyPair["name"] == id {
			return keyPair, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
}
