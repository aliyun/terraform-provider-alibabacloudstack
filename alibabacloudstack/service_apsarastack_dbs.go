package alibabacloudstack

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type DbsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DbsService) DescribeDbsBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}

	request := map[string]interface{}{
		"BackupPlanId": id,
		"PageSize":     PageSizeLarge,
		"PageNumber":   1,
	}
	request["ClientToken"] = buildClientToken("DescribeBackupPlanList")

	response, err = s.client.DoTeaRequest("POST", "dbs", "2019-03-06", "DescribeBackupPlanList", "", nil, nil, request)
	if err != nil {
		return object, err
	}

	v, err := jsonpath.Get("$.Items.BackupPlanDetail", response)
	i := v.([]interface{})
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Items.BackupPlanDetail", response)
	}
	if len(i) > 0 {
		object = i[0].(map[string]interface{})
	} else {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("dbs", id)), errmsgs.NotFoundWithResponse, response)
	}

	return object, nil
}
