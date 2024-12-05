package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type OosService struct {
	client *connectivity.AlibabacloudStackClient
}
func (s *OosService) DoOosGettemplateRequest(id string) (object map[string]interface{}, err error) {
    return s.DescribeOosTemplate(id)
}

func (s *OosService) DescribeOosTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"TemplateName": id,
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "Oos", "2019-06-01", "GetTemplate", "", nil, request)
	addDebug("GetTemplate", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExists.Template"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OosTemplate", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Template", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Template", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *OosService) DoOosListexecutionsRequest(id string) (object map[string]interface{}, err error) {
    return s.DescribeOosExecution(id)
}
func (s *OosService) DescribeOosExecution(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"ExecutionId": id,
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "Oos", "2019-06-01", "ListExecutions", "", nil, request)
	addDebug("ListExecutions", response, request)
	if err != nil {
		return
	}
	v, err := jsonpath.Get("$.Executions", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Executions", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OOS", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["ExecutionId"].(string) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("OOS", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *OosService) OosExecutionStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeOosExecution(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}
