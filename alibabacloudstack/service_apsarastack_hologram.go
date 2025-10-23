package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type HologramService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *HologramService) DescribeHologramInstance(id string) (map[string]interface{}, error) {
	pattern := fmt.Sprintf("/api/v1/instances/%s", id)
	response, err := s.client.DoTeaRequest("GET", "Hologram", "2022-06-01", "GetInstance", pattern, nil, nil, nil)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologress_instance", "ListInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	instance, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologress_instance", "ListInstances", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if instance != nil {
		return instance.(map[string]interface{}), nil
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("hologress_instance", id))
}

func (s *HologramService) HologramInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeHologramInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["InstanceStatus"] == failState {
				return object, object["InstanceStatus"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["InstanceStatus"]))
			}
		}
		return object, object["InstanceStatus"].(string), nil
	}
}

func (s *HologramService) CalculateQuota(node int, cluster string) (map[string]interface{}, error) {

	body := map[string]interface{}{
		"node":    node,
		"cluster": cluster,
	}
	request := map[string]interface{}{
		"x-acs-body": body,
	}
	response, err := s.client.DoTeaRequest("POST", "Hologram", "2022-06-01", "CalculateQuota", "/api/v1/instances/calculateQuota", nil, nil, request)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologress_instance", "CalculateQuota", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	quota, err := jsonpath.Get("$.Quota", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologress_instance", "CalculateQuota", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return quota.(map[string]interface{}), nil
}
