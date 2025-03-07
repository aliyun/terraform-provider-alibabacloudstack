package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type CsbService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *CsbService) DescribeCsbProjectDetail(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, err
	}
	request := map[string]interface{}{
		"CsbId":       parts[0],
		"ProjectName": parts[1],
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "CSB", "2017-11-18", "GetProject", "", nil, nil, request)
	if err != nil {
		return object, err
	}

	v, err := jsonpath.Get("$.Data.ProjectList", response)
	i := v.([]interface{})
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Data.ProjectList", response)
	}
	if len(i) > 0 {
		object = i[0].(map[string]interface{})
	} else {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("csb", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}
