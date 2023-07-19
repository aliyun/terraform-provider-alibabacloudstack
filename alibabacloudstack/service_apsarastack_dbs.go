package alibabacloudstack

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type DbsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DbsService) DescribeDbsBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDbsClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeBackupPlanList"

	request := map[string]interface{}{
		"BackupPlanId": id,
		"RegionId":     s.client.RegionId,
	}
	request["product"] = "dbs"
	request["Product"] = "dbs"
	request["ClientToken"] = buildClientToken("DescribeBackupPlanList")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-06"), StringPointer("AK"), request, nil, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabacloudStackSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Items.BackupPlanDetail", response)
	i := v.([]interface{})
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.BackupPlanDetail", response)
	}
	if len(i) > 0 {
		object = i[0].(map[string]interface{})
	} else {
		return object, WrapErrorf(Error(GetNotFoundMessage("dbs", id)), NotFoundWithResponse, response)
	}

	return object, nil
}
