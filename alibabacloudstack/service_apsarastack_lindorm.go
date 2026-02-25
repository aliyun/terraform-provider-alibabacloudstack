package alibabacloudstack

import (
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
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
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
