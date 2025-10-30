package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type ApiGateWayV2Service struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *ApiGateWayV2Service) DescribeApiGatewayV2Instace(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"gwInstanceId": id,
	}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetInstanceInfo", "/gatewayInstance/getInstanceInfo", nil, request, nil)
	if err != nil {
		return nil, err
	}
	data, ok := response["data"]
	if !ok {
		return nil, errmsgs.Error("CreateInstance Failed! %v", response)
	}
	return data.(map[string]interface{}), nil
}
