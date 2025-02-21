package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type Dms_enterpriseService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *Dms_enterpriseService) DescribeDmsEnterpriseInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"Host":     parts[0],
		"Port":     parts[1],
		"PageSize": PageSizeLarge,
		"PageNumber": 1,
	}
	response, err = s.client.DoTeaRequest("POST", "dms-enterprise", "2018-11-01", "GetInstance", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InstanceNoEnoughNumber"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DmsEnterpriseInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Instance", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Instance", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Dms_enterpriseService) DoDms_EnterpriseGetuserRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeDmsEnterpriseUser(id)
}

func (s *Dms_enterpriseService) DescribeDmsEnterpriseUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"Uid":      id,
		"PageSize": PageSizeLarge,
		"PageNumber": 1,
	}
	response, err = s.client.DoTeaRequest("POST", "dms-enterprise", "2018-11-01", "GetUser", "", nil, nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.User", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.User", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
