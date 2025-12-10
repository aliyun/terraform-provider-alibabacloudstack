package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

type Schedulerx2Service struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *Schedulerx2Service) DescribeSchedulerx2AppGroup(id string) (map[string]interface{}, error) {
	reqQuery := map[string]interface{}{
		"Action":      "ListGroups",
		"AccessKeyId": s.client.AccessKey,
		"Namespace":   "system_namespace",
		"PageNum":     1,
		"PageSize":    10,
	}
	response, err := s.client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "ListGroups", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	if records, err := jsonpath.Get("$.Data.Records", response); err == nil {
		for _, record := range records.([]interface{}) {
			r := record.(map[string]interface{})
			if fmt.Sprintf("%v", r["Id"]) == id {
				return r, nil
			}
		}
	} else {
		return nil, errmsgs.WrapError(err)
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Schedulerx2 AppGroup %s was not found", id))
}

func (s *Schedulerx2Service) DescribeSchedulerx2Job(id string) (map[string]interface{}, error) {

	reqQuery := map[string]interface{}{
		"Action":      "ListJobs",
		"AccessKeyId": s.client.AccessKey,
		"Namespace":   "system_namespace",
		"PageNum":     1,
		"PageSize":    10,
	}

	response, err := s.client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "ListJobs", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	if records, err := jsonpath.Get("$.Data.Records", response); err == nil {
		for _, record := range records.([]interface{}) {
			r := record.(map[string]interface{})
			if fmt.Sprintf("%v", r["JobId"]) == id {
				return r, nil
			}
		}
	} else {
		return nil, errmsgs.WrapError(err)
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Schedulerx2 Job %s was not found", id))
}
