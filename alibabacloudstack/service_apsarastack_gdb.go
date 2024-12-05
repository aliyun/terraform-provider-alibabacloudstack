package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type GdbService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *GdbService) DoGdbDescribedbinstanceaccesswhitelistRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeGraphDatabaseDbInstance(id)
}

func (s *GdbService) DescribeGraphDatabaseDbInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "gdb", "2019-09-03", "DescribeDBInstanceAttribute", "", nil, request)
	addDebug("DescribeDBInstanceAttribute", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound", "InvalidDBInstanceId.NotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GraphDatabase:DbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Items.DBInstance", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Items.DBInstance", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GraphDatabase", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DBInstanceId"]) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GraphDatabase", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})

	GetDBInstanceAccessWhiteListObject, err := s.GetDBInstanceAccessWhiteList(id)
	if err != nil {
		return nil, err
	}

	object["DBInstanceIPArray"] = GetDBInstanceAccessWhiteListObject["DBInstanceIPArray"]
	return object, nil
}

func (s *GdbService) GetDBInstanceAccessWhiteList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "gdb", "2019-09-03", "DescribeDBInstanceAccessWhiteList", "", nil, request)
	addDebug("DescribeDBInstanceAccessWhiteList", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound", "InvalidDBInstanceId.NotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GraphDatabase:DbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Items", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Items", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *GdbService) GraphDatabaseDbInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGraphDatabaseDbInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DBInstanceStatus"]) == failState {
				return object, fmt.Sprint(object["DBInstanceStatus"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["DBInstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["DBInstanceStatus"]), nil
	}
}

func (s *GdbService) DescribeDBInstanceAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "gdb", "2019-09-03", "DescribeDBInstanceAttribute", "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidDBInstance.NotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GraphDatabase:DbInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Items.DBInstance", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Items.DBInstance", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GraphDatabase", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DBInstanceId"]) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("GraphDatabase", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}
