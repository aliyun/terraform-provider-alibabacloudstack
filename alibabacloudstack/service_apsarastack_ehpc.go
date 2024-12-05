package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type EhpcService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *EhpcService) DescribeEhpcJobTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}
	idExist := false
	for {
		response, err = s.client.DoTeaRequest("GET", "EHPC", "2018-04-12", "ListJobTemplates", "", nil, request)
		if err != nil {
			return object, err
		}
		v, err := jsonpath.Get("$.Templates.JobTemplates", response)
		if err != nil {
			return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Templates.JobTemplates", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Ehpc", id)), errmsgs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Id"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Ehpc", id)), errmsgs.NotFoundWithResponse, response)
	}
	return
}
