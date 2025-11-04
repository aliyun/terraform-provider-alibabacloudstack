package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ApiGateWayV2Service struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *ApiGateWayV2Service) DescribeApiGatewayV2Instance(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"gwInstanceId": id,
	}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetInstanceInfo", "/gatewayInstance/getInstanceInfo", nil, request, request)
	if err != nil {
		return nil, err
	}
	data, ok := response["data"]
	if !ok {
		return nil, errmsgs.Error("GetInstanceInfo Failed! %v", response)
	}
	return data.(map[string]interface{}), nil
}

func (s *ApiGateWayV2Service) ApiGateWayV2InstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeApiGatewayV2Instance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["status"]) == failState {
				return object, fmt.Sprint(object["status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["status"])))
			}
		}
		return object, fmt.Sprint(object["status"]), nil
	}
}

func (s *ApiGateWayV2Service) GetCustomDeployConfig(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"gwInstanceId": id,
	}
	response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetCustomDeployConfig", "/gatewayInstance/getCustomDeployConfig", nil, request, request)
	if err != nil {
		return nil, err
	}
	data, ok := response["data"]
	if !ok {
		return nil, errmsgs.Error("CreateInstance Failed! %v", response)
	}
	return data.(map[string]interface{}), nil
}

 func (s *ApiGateWayV2Service) DescribeApigatewayv2K8sCluster(id string) (map[string]interface{}, error) {
    reqBody := map[string]interface{}{
        "current": 1,
        "size":    100,
    }

    response, err := s.client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListClusters", "/k8s/listClusters", nil, nil, reqBody)
    if err != nil {
        return nil, err
    }

    if records, ok := response["data"].(map[string]interface{})["records"].([]interface{}); ok {
        for _, record := range records {
            item := record.(map[string]interface{})
            if k8sClusterCode, ok := item["k8sClusterCode"].(string); ok && k8sClusterCode == id {
                return item, nil
            }
        }
    }

    return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Apigatewayv2 K8s Cluster %s was not found", id))
}
