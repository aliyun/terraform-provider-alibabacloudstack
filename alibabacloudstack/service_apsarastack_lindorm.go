package alibabacloudstack

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type LindormService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *LindormService) DescribeLindormInstance(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{
		"InstanceId": id,
	}

	response, err := s.client.DoTeaRequest("GET", "hitsdb", "2020-06-15", "GetLindormInstance", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *LindormService) LindormInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeLindormInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, "InstanceInvalid") {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			if errmsgs.IsExpectedErrors(err, "Forbidden", "Forbidden.RAM") {
				return nil, "CREATING", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		instanceStatus := object["InstanceStatus"].(string)

		for _, failState := range failStates {
			if instanceStatus == failState {
				return object, instanceStatus, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, instanceStatus))
			}
		}
		return object, instanceStatus, nil
	}
}

func (s *LindormService) DescribeLindormConnection(id string) (map[string]interface{}, error) {
	params := strings.Split(id, ":")
	request := map[string]interface{}{
		"InstanceId": params[0],
	}
	resp, err := s.client.DoTeaRequest("GET", "hitsdb", "2020-06-15", "GetLindormInstanceEngineList", "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	data, ok := resp["EngineList"]
	if !ok {
		return nil, errmsgs.Error("GetLindormInstanceEngineList Failed! %v", resp)
	}
	for _, v := range data.([]interface{}) {
		item := v.(map[string]interface{})
		if item["EngineType"] == params[1] {
			return item, nil
		}
	}
	return nil, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Lindorm Connection", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}
