package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type DrdsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DrdsService) DescribeDrdsInstance(id string) (*drds.DescribeDrdsInstanceResponse, error) {
	request := drds.CreateDescribeDrdsInstanceRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.DrdsInstanceId = id
	raw, err := s.client.WithDrdsClient(func(drdsClient *drds.Client) (interface{}, error) {
		return drdsClient.DescribeDrdsInstance(request)
	})

	response, ok := raw.(*drds.DescribeDrdsInstanceResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDrdsInstanceId.NotFound"}) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	
	if response.Data.Status == "5" {
		return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return response, nil
}

func (s *DrdsService) DrdsInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDrdsInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object.Data.Status == failState {
				return object, object.Data.Status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object.Data.Status))
			}
		}

		return object, object.Data.Status, nil
	}
}

func (s *DrdsService) WaitDrdsInstanceConfigEffect(id string, item map[string]string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		effected := false
		object, err := s.DescribeDrdsInstance(id)

		if err != nil {
			if errmsgs.NotFoundError(err) {
				return errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
			}
			return errmsgs.WrapError(err)
		}

		if value, ok := item["description"]; ok {
			if object.Data.Description == value {
				effected = true
			}
		}

		if effected {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Data, item, errmsgs.ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}

	return nil
}
