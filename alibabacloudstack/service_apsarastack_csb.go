package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

type CsbService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *CsbService) DescribeCsbProjectDetail(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDtsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	parts, err := ParseResourceId(id, 2)
	action := "GetProject"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"CsbId":       parts[0],
		"ProjectName": parts[1],
	}
	request["Product"] = "CSB"
	request["OrganizationId"] = s.client.Department
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)

	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-18"), StringPointer("AK"), nil, request, &runtime)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data.ProjectList", response)
	i := v.([]interface{})
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.ProjectList", response)
	}
	if len(i) > 0 {
		object = i[0].(map[string]interface{})
	} else {
		return object, WrapErrorf(Error(GetNotFoundMessage("csb", id)), NotFoundWithResponse, response)
	}
	return object, nil
}
